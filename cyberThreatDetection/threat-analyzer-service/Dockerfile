FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY shared/ ./shared/
COPY threat-analyzer-service/ ./threat-analyzer-service/

RUN CGO_ENABLED=0 GOOS=linux go build -o threat-service ./threat-analyzer-service/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/threat-service .

EXPOSE 8082
CMD ["./threat-service"]