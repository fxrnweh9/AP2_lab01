FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
COPY vendor ./vendor
RUN go build -mod=vendor -o notification-service ./cmd/notification-service

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/notification-service .
CMD ["./notification-service"]