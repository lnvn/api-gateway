# Build Stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway .

# FROM alpine:latest
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app/

COPY --from=builder /app/api-gateway .
COPY --from=builder /app/config.json .

EXPOSE 8080

CMD ["./api-gateway"]
