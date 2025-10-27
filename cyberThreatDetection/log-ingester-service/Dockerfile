FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY shared/ ./shared/
COPY log-ingester-service/ ./log-ingester-service/

RUN CGO_ENABLED=0 GOOS=linux go build -o log-service ./log-ingester-service/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/log-service .

EXPOSE 8081
CMD ["./log-service"]