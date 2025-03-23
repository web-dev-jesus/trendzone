# Dockerfile

# Build stage
FROM golang:1.20-alpine AS builder

# Set working directory
WORKDIR /app

# Install required build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o nfl-data-sync ./cmd/nfl-data-sync

# Final stage
FROM alpine:3.17

# Install ca-certificates and tzdata
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/nfl-data-sync .

# Create non-root user
RUN adduser -D -g '' appuser
USER appuser

# Command to run the executable
ENTRYPOINT ["./nfl-data-sync"]