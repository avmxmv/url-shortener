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
│   ├── 20250216005319_add_index_to_links.sql
│   └── migrations.go
├── .env
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

# Клонируйте проект

```
git clone https://github.com/avmxmv/url-shortener.git
```

# Запуск сервера:

## Через docker: 
```
make docker-compose-up
```

## Остановка: 
```
make docker-compose-down
```

# Локально:

## Сборка проекта:
```
make build
```

## inmem: 
```
make run
```

## postgres: 
```
make run-postgres
```

## Генерация gRPC кода: 
```
make generate
```

## Применение миграций:
```
migrate-up
```

## Запуск тестов:
```
make tests
```

# Примеры запросов:

## gRPC:

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

## HTTP API:

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
