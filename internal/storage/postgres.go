package storage

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) CreateLink(originalURL string) (string, error) {
	var shortURL string
	err := s.db.QueryRow(
		`INSERT INTO links (original_url, short_url) 
		VALUES ($1, $2) 
		ON CONFLICT (original_url) DO NOTHING
		RETURNING short_url`,
		originalURL, GenerateShortURL(),
	).Scan(&shortURL)

	return shortURL, err
}

func (s *PostgresStorage) GetLink(shortURL string) (string, error) {
	var originalURL string
	err := s.db.QueryRow(
		"SELECT original_url FROM links WHERE short_url = $1",
		shortURL,
	).Scan(&originalURL)

	if err == sql.ErrNoRows {
		return "", ErrNotFound
	}
	return originalURL, err
}
