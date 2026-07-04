# Docker Deployment Guide for GoTimeDate

A high-performance REST API and WebSocket server for time operations built with Go and Fiber.

**Note:** This project uses pure Go with no CGO dependencies, resulting in small, efficient container images.

## Quick Start

### Using Docker

```bash
# Pull the image
docker pull ghcr.io/shabilullah/gotimedate:latest

# Run the container
docker run -d \
  --name gotimedate \
  -p 8080:8080 \
  -v $(pwd)/logs:/app/logs \
  ghcr.io/shabilullah/gotimedate:latest
```

### Using Docker Compose

```bash
# Start the service
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the service
docker-compose down
```

## Building Locally

```bash
# Build the image
docker build -t gotimedate:local .

# Run the local image
docker run -d -p 8080:8080 gotimedate:local
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `HOST` | `0.0.0.0` | Server host |
| `PREFORK` | `false` | Enable Fiber prefork (spawns child processes per CPU core) |
| `LOG_FILE` | `/app/logs/server.log` | Log file path |
| `LOG_LEVEL` | `info` | Log level (debug, info, warn, error) |
| `ALLOWED_ORIGINS` | `*` | CORS allowed origins |
| `ALLOWED_METHODS` | `GET,POST,PUT,DELETE,OPTIONS` | CORS allowed methods |
| `ALLOWED_HEADERS` | `Origin,Content-Type,Accept,Authorization` | CORS allowed headers |
| `ALLOW_CREDENTIALS` | `true` | CORS allow credentials |
| `MAX_AGE` | `3600` | CORS preflight cache duration |
| `WS_PING_INTERVAL` | `30` | WebSocket ping interval (seconds) |
| `WS_PONG_WAIT` | `60` | WebSocket pong wait timeout (seconds) |
| `WS_WRITE_WAIT` | `10` | WebSocket write timeout (seconds) |

## Production Deployment

### With Custom Configuration

```bash
docker run -d \
  --name gotimedate \
  -p 8080:8080 \
  -e PORT=8080 \
  -e HOST=0.0.0.0 \
  -e PREFORK=true \
  -e LOG_LEVEL=warn \
  -e ALLOWED_ORIGINS="https://example.com,https://app.example.com" \
  -v $(pwd)/logs:/app/logs \
  --restart unless-stopped \
  ghcr.io/shabilullah/gotimedate:latest
```

### Health Check

The container includes a health check that monitors:
- HTTP endpoint: `http://localhost:8080/health`
- Check interval: 30 seconds
- Timeout: 10 seconds
- Retries: 3

## GitHub Container Registry

Images are automatically published to GitHub Container Registry (GHCR) on:
- Push to `main`/`master` branch → `latest` tag
- Push to any branch → `<branch-name>` tag
- Git tags starting with `v` → versioned tags (`v1.2.3`, `1.2`, `1`)
- Pull requests → `pr-<number>` tag
- Commits → `sha-<short-hash>` tag

### Pull Specific Version

```bash
# Latest
docker pull ghcr.io/shabilullah/gotimedate:latest

# Specific version
docker pull ghcr.io/shabilullah/gotimedate:v1.0.0

# Specific commit
docker pull ghcr.io/shabilullah/gotimedate:sha-abc1234
```

## Multi-Platform Support

The workflow builds images for:
- `linux/amd64` (x86_64)
- `linux/arm64` (ARM 64-bit, including Apple Silicon)

## Volumes

- `/app/logs` - Application logs directory

Mount this directory to persist logs between container restarts:

```bash
docker run -v ./logs:/app/logs ghcr.io/shabilullah/gotimedate:latest
```

## Troubleshooting

### View Logs

```bash
# Container logs
docker logs gotimedate

# Follow logs
docker logs -f gotimedate

# Application logs (if volume mounted)
tail -f logs/server.log
```

### Access Container Shell

```bash
docker exec -it gotimedate sh
```

### Check Health Status

```bash
docker inspect --format='{{.State.Health.Status}}' gotimedate
```

## Architecture

### Multi-Stage Build

The Dockerfile uses a two-stage build process:

1. **Builder stage** (`golang:1.24-alpine`)
   - Downloads dependencies
   - Compiles the Go application with `CGO_ENABLED=0`
   - Produces a static binary with no C dependencies

2. **Runtime stage** (`alpine:3.21`)
   - Minimal base image (~5MB)
   - Only includes the compiled binary, ca-certificates, and tzdata
   - Final image size: ~15-20MB

### Why CGO_ENABLED=0?

This project uses pure Go packages with no C dependencies (no SQLite, no native libraries). Building with `CGO_ENABLED=0` produces:
- Fully static binaries
- Smaller container images
- Better portability across platforms
- No need for C toolchain in runtime container
