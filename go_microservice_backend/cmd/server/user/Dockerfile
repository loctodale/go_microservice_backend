# Build Stage
FROM golang:alpine AS builder

# Install required tools
RUN apk add --no-cache openssl

# Set working directory
WORKDIR /build

# Copy all files and download dependencies
COPY ../../.. .
RUN go mod download

# Build the Go application
RUN go build -o crm.shopdev.com ./cmd/server/user

# Final Stage
FROM alpine:latest

# Install CA certificates (required for trusting certs)
RUN apk add --no-cache ca-certificates

# Copy application binary
COPY --from=builder /build/crm.shopdev.com /

# Copy configuration files
COPY ../../../config /config

# Copy SSL certificates
COPY ../../../certs/cert.crt /usr/local/share/ca-certificates/custom-cert.crt

# Add the certificate to the system's trusted store
RUN update-ca-certificates

# Copy private key for the application to use (if needed)
COPY ../../../certs/key.pem /etc/ssl/private/key.pem

# Set the entrypoint
ENTRYPOINT ["/crm.shopdev.com", "config/local.yaml"]
