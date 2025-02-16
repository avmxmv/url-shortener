Структура проекта:

```
.
├── api
│   ├── link_service.pb.go
│   ├── link_service.proto
│   └── link_service_grpc.pb.go
├── cmd
│   └── server
│       └── main.go
├── internal
│   ├── config
│   │   └── config.go
│   ├── handler
│   │   ├── create.go
│   │   └── get.go
│   ├── storage
│   │   ├── inmem.go
│   │   ├── inmem_test.go
│   │   ├── postgres.go
│   │   ├── postgres_test.go
│   │   ├── shortener.go
│   │   └── storage.go
│   └── service
│       ├── service.go
│       └── service_test.go
├── migrations
│   ├── 20250214102327_create_links_table.sql
│   └── migrations.go
├── Dockerfile
├── docker-compose.yml
├── .env
├── go.mod
├── go.sum
└── README.md
```

Клонируйте проект

```git clone https://github.com/avmxmv/url-shortener.git```

Запуск сервера:

через docker - ```docker-compose up -d --build```

Локально:

inmem - ```go run ./cmd/server --storage=inmem```

postgres - ```go run ./cmd/server --storage=postgres```

Примеры запросов

grpc:

CreateLink:

```
grpcurl -plaintext -d '{"original_url": "https://avmxmv.com"}' localhost:50051 api.LinkService/CreateLink
```

Пример ответа:

```
{
  "shortUrl": "JHt3oBNWn4"
}
```

GetLink :

```
grpcurl -plaintext -d '{"short_url": "JHt3oBNWn4"}' localhost:50051 api.LinkService/GetLink
```

Пример ответа:

```
{
  "originalUrl": "https://avmxmv.com"
}
```

Пример ответа когда не надено:

```
ERROR:
  Code: Unknown
  Message: not found
```

HTTP API:

POST:

```
curl -X POST -d "url=https://avmxmv.com" http://localhost:8080/create
```

Пример ответа:

```fZKh7hZIeV```

GET:

```
curl http://localhost:8080/get/fZKh7hZIeV
```

Пример ответа:

```
<a href="https://avmxmv.com">Found</a>
```

Пример ответа при Invalid URL:

```
Invalid short URL
```

Пример ответа когда не найдено:

```
not found
```

Запуск тестов:

```
go test ./...
```