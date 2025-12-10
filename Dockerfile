# Build Stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./
# COPY go.sum ./ # No go.sum yet

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway .

# Final Stage
FROM alpine:latest

WORKDIR /app/

# Copy the Pre-built binary from the previous stage
COPY --from=builder /app/api-gateway .
COPY --from=builder /app/config.json .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./api-gateway"]
