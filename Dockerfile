# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/trendzone ./cmd/server

# Run stage
FROM alpine:3.18

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/bin/trendzone .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["/app/trendzone"]