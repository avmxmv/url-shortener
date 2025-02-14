package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"url-shortener/api"
	"url-shortener/internal/config"
	"url-shortener/internal/service"
	"url-shortener/internal/storage"
)

func main() {
	// Загрузка конфигурации
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Парсинг флагов
	var storageType string
	flag.StringVar(&storageType, "storage", cfg.StorageType, "storage type (inmem|postgres)")
	flag.Parse()

	// Инициализация хранилища
	var store storage.Storage
	switch storageType {
	case "inmem":
		store = storage.NewInMemStorage()
		log.Println("Using in-memory storage")

	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

		store, err = storage.NewPostgresStorage(dsn)
		if err != nil {
			log.Fatalf("Failed to connect to PostgreSQL: %v", err)
		}
		log.Println("Connected to PostgreSQL")

		// Применение миграций
		if err := runMigrations(dsn); err != nil {
			log.Fatalf("Migrations failed: %v", err)
		}

	default:
		log.Fatalf("Unknown storage type: %s", storageType)
	}

	// Инициализация сервиса
	svc := service.NewService(store)

	// Graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запуск gRPC сервера
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		s := grpc.NewServer()
		api.RegisterLinkServiceServer(s, svc)
		log.Printf("gRPC server listening on :%s", cfg.GRPCPort)

		// Используем контекст для graceful shutdown
		go func() {
			<-ctx.Done()
			s.GracefulStop()
		}()

		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Запуск HTTP сервера
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/create", createHandler(svc))
		mux.HandleFunc("/get/", getHandler(svc))

		server := &http.Server{
			Addr:    fmt.Sprintf(":%s", cfg.HTTPPort),
			Handler: mux,
		}

		log.Printf("HTTP server listening on :%s", cfg.HTTPPort)

		// Используем контекст для graceful shutdown
		go func() {
			<-ctx.Done()
			server.Shutdown(context.Background())
		}()

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down server...")
	cancel() // Отменяем контекст для graceful shutdown
}

func runMigrations(dsn string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS links (
		id SERIAL PRIMARY KEY,
		original_url TEXT NOT NULL UNIQUE,
		short_url TEXT NOT NULL UNIQUE,
		created_at TIMESTAMP DEFAULT NOW()
	)`)
	return err
}

func createHandler(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		originalURL := r.FormValue("url")
		if originalURL == "" {
			http.Error(w, "Missing URL parameter", http.StatusBadRequest)
			return
		}

		resp, err := svc.CreateLink(r.Context(), &api.CreateLinkRequest{
			OriginalUrl: originalURL,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, resp.ShortUrl)
	}
}

func getHandler(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		shortURL := r.URL.Path[len("/get/"):]
		if len(shortURL) != 10 {
			http.Error(w, "Invalid short URL", http.StatusBadRequest)
			return
		}

		resp, err := svc.GetLink(r.Context(), &api.GetLinkRequest{
			ShortUrl: shortURL,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Redirect(w, r, resp.OriginalUrl, http.StatusFound)
	}
}
