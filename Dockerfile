# Stage 1: Build
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code including embedded static files
COPY . .

# Build the application with CGO disabled (pure Go, no C dependencies)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /gotimedate .

# Stage 2: Runtime
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata wget

WORKDIR /app

# Copy binary from builder
COPY --from=builder /gotimedate /app/gotimedate

# Environment variables with defaults
ENV PORT=8080 \
    HOST=0.0.0.0 \
    LOG_FILE=logs/server.log \
    LOG_LEVEL=info

# Create logs directory
RUN mkdir -p /app/logs

# Expose port
EXPOSE 8080

# Run the application
CMD ["/app/gotimedate"]
