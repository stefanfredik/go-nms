# --- Builder Stage ---
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go module files and download dependencies first (layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the api-gateway binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/api-gateway ./cmd/api-gateway

# --- Runtime Stage ---
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/api-gateway .

EXPOSE 9090

ENTRYPOINT ["/app/api-gateway"]
