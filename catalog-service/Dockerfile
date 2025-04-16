FROM golang:1.24-alpine

RUN apk add --no-cache git

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /app

COPY users-service/ ./users-service

COPY catalog-service/go.mod catalog-service/go.sum ./catalog-service/
WORKDIR /app/catalog-service
RUN go mod tidy

COPY catalog-service/ ./

RUN go build -o /app/catalog-service-app ./cmd

RUN chmod +x /app/catalog-service-app

CMD ["/app/catalog-service-app"]
