FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o billing-service .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/billing-service .
COPY .env .

EXPOSE 8080

CMD ["./billing-service"]