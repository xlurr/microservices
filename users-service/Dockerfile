# Builder stage
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

RUN adduser -D -g '' appuser

WORKDIR /build

COPY go.mod ./

RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest

RUN swag init -g cmd/main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags='-w -s' \
    -o users-service ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates wget

# Создаём пользователя в final stage
RUN adduser -D -g '' appuser

WORKDIR /app

# Копируем бинарник и docs
COPY --from=builder /build/users-service /app/users-service
COPY --from=builder /build/docs ./docs

# Создаём data директорию и инициализируем JSON
RUN mkdir -p /app/data && \
    echo '[]' > /app/data/users.json

# Устанавливаем правильные права доступа
RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 8001

CMD ["/app/users-service"]
