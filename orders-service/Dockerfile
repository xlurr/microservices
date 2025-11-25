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

FROM alpine:latest
RUN apk --no-cache add ca-certificates wget

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/orders-service /orders-service
COPY --from=builder /build/docs ./docs

USER appuser
EXPOSE 8082

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s \
    CMD wget --quiet --tries=1 --spider http://localhost:8082/health || exit 1

ENTRYPOINT ["/orders-service"]
