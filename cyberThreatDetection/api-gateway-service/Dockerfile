FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY shared/ ./shared/
COPY api-gateway-service/ ./api-gateway-service/

RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway ./api-gateway-service/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/api-gateway .

EXPOSE 8080
CMD ["./api-gateway"]