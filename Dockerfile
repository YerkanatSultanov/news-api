FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o news-api ./cmd/

FROM debian:bookworm-slim

WORKDIR /app

COPY migrations ./migrations
COPY --from=builder /app/news-api .
COPY .env .

EXPOSE 8080

CMD ["./news-api"]
