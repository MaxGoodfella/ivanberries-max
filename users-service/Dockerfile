FROM golang:1.24-alpine

RUN apk add --no-cache git

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

COPY db/migrations /app/users-service/db/migrations

RUN go build -o /app/users-service-app ./cmd

RUN chmod +x /app/users-service-app

CMD ["/app/users-service-app"]