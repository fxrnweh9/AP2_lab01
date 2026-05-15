FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
COPY vendor ./vendor
RUN go build -mod=vendor -o order-service ./cmd/order_service

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/order-service .
CMD ["./order-service"]