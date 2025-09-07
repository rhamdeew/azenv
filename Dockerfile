# Build stage
FROM golang:1.23-alpine AS builder

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY main.go .

# Build the binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o azenv main.go

# Final stage
FROM alpine:3.19

# Install ca-certificates for HTTPS/TLS functionality
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1050 -S azenv && \
    adduser -u 1050 -S azenv -G azenv

# Create directories for certificates
RUN mkdir -p /app/cert /app/cert-cache && \
    chown -R azenv:azenv /app

# Copy binary from builder stage
COPY --from=builder /app/azenv /app/azenv

# Change to non-root user
USER azenv

# Set working directory
WORKDIR /app

# Expose ports (HTTP and HTTPS)
EXPOSE 8080 8443

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/azenv || exit 1

# Default command
CMD ["./azenv"]
