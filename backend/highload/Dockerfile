FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -ldflags="-s -w" -o app main.go

FROM alpine:3.19

WORKDIR /root/
COPY --from=builder /app/app .
EXPOSE 8081

CMD ["./app"]