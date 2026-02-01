# ---------- build stage ----------
FROM golang:1.22-alpine AS builder

WORKDIR /app

# зависимости
COPY go.mod go.sum ./
RUN go mod download

# исходники
COPY . .

# сборка
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./cmd/server

# ---------- runtime stage ----------
FROM alpine:3.19

WORKDIR /app

# сертификаты (для https / postgres ssl)
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/app .
COPY config/config.yaml ./config/config.yaml
COPY migrations ./migrations

EXPOSE 8080

CMD ["./app"]
