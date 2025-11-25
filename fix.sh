#!/bin/bash

# Обновление Dockerfile с правильным копированием docs папки

echo "🔄 Обновляем Dockerfile для всех сервисов..."

# users-service/Dockerfile
cat > users-service/Dockerfile << 'EOF'
# Стадия сборки
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
    -o /app/users-service ./cmd/main.go

# Финальная стадия
FROM alpine:latest

RUN apk --no-cache add ca-certificates wget

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/users-service /users-service
COPY --from=builder /build/docs /docs

USER appuser

EXPOSE 8081

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s \
    CMD wget --quiet --tries=1 --spider http://localhost:8081/health || exit 1

ENTRYPOINT ["/users-service"]
EOF

# orders-service/Dockerfile
cat > orders-service/Dockerfile << 'EOF'
# Стадия сборки
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
    -o /app/orders-service ./cmd/main.go

# Финальная стадия
FROM alpine:latest

RUN apk --no-cache add ca-certificates wget

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/orders-service /orders-service
COPY --from=builder /build/docs /docs

USER appuser

EXPOSE 8082

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s \
    CMD wget --quiet --tries=1 --spider http://localhost:8082/health || exit 1

ENTRYPOINT ["/orders-service"]
EOF

# payments-service/Dockerfile
cat > payments-service/Dockerfile << 'EOF'
# Стадия сборки
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates
RUN adduser -D -g '' appuser

WORKDIR /build

COPY go.mod ./

RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags='-w -s' \
    -o /app/payments-service ./cmd/main.go

# Финальная стадия
FROM alpine:latest

RUN apk --no-cache add ca-certificates wget

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/payments-service /payments-service
COPY --from=builder /build/docs /docs

USER appuser

EXPOSE 8083

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s \
    CMD wget --quiet --tries=1 --spider http://localhost:8083/health || exit 1

ENTRYPOINT ["/payments-service"]
EOF

# delivery-service/Dockerfile
cat > delivery-service/Dockerfile << 'EOF'
# Стадия сборки
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates
RUN adduser -D -g '' appuser

WORKDIR /build

COPY go.mod ./

RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags='-w -s' \
    -o /app/delivery-service ./cmd/main.go

# Финальная стадия
FROM alpine:latest

RUN apk --no-cache add ca-certificates wget

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/delivery-service /delivery-service
COPY --from=builder /build/docs /docs

USER appuser

EXPOSE 8084

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s \
    CMD wget --quiet --tries=1 --spider http://localhost:8084/health || exit 1

ENTRYPOINT ["/delivery-service"]
EOF

echo "✅ Все Dockerfile обновлены с копированием /docs папки"

echo -e "\n🐳 Теперь выполните:\n"
echo "docker-compose down --rmi all"
echo "docker-compose build --no-cache"
echo "docker-compose up"