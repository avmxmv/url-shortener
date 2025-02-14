FROM golang:1.23.4 AS builder

RUN apt-get update && apt-get install -y protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /app
COPY ../../url-shortener .

RUN go mod download
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/link_service.proto
RUN CGO_ENABLED=0 GOOS=linux go build -o /url-shortener ./cmd/server

FROM gcr.io/distroless/base-debian11

COPY --from=builder /url-shortener /url-shortener
COPY --from=builder /app/migrations /migrations

EXPOSE 8080 50051
CMD ["/url-shortener"]