FROM golang:1.24.0 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build -o pizza-backend .

FROM debian:12.8-slim AS pizza
WORKDIR /app
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/pizza-backend ./
ENTRYPOINT ["./pizza-backend"]