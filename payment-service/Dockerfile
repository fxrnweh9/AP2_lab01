FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
COPY vendor ./vendor
RUN go build -mod=vendor -o payment-service ./cmd/payment-service

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/payment-service .
CMD ["./payment-service"]