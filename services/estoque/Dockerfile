# ─── Estágio 1: Build ────────────────────────────────────────────────────────
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/estoque ./cmd/main.go

# ─── Estágio 2: Runtime ──────────────────────────────────────────────────────
FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/estoque .

EXPOSE 8080

CMD ["./estoque"]
