# Stage 1: Build
FROM golang:1.25-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum first to leverage Docker cache
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy the rest of the backend source code
COPY backend/ .

# Build the application
# Using -ldflags="-s -w" to reduce binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ledgeguard-api ./cmd/api

# Stage 2: Runtime
FROM alpine:latest

# Install runtime dependencies (like ca-certificates for HTTPS)
RUN apk add --no-cache ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/ledgeguard-api .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./ledgeguard-api"]
