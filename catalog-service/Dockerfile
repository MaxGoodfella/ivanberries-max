FROM golang:1.24-alpine

RUN apk add --no-cache git

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

COPY db/migrations /app/catalog-service/db/migrations

RUN go build -o /app/catalog-service-app ./cmd

RUN chmod +x /app/catalog-service-app

CMD ["/app/catalog-service-app"]
