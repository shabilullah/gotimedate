# GoTimeDate API

A high-performance REST API and WebSocket server for time operations built with Go and Fiber.

## Features

- 🚀 Fast and lightweight using Fiber framework
- 🔌 WebSocket support for real-time time updates
- 📝 Swagger/OpenAPI documentation
- 🐳 Docker support with multi-platform builds
- ⚙️ Configurable via environment variables
- 📊 Structured logging
- 🌐 CORS support

## Quick Start

### Local Development

```bash
# Clone the repository
git clone https://github.com/shabilullah/gotimedate.git
cd gotimedate

# Install dependencies
go mod download

# Run the server
go run main.go
```

The server will start on `http://localhost:8080`

### Using Docker

```bash
# Pull and run from GitHub Container Registry
docker pull ghcr.io/shabilullah/gotimedate:latest
docker run -d -p 8080:8080 ghcr.io/shabilullah/gotimedate:latest
```

For detailed Docker deployment instructions, see [README.Docker.md](README.Docker.md)

## API Documentation

Once the server is running, visit:
- Swagger UI: `http://localhost:8080/swagger/`
- Health check: `http://localhost:8080/health`
- API base: `http://localhost:8080/api/v1/`

### Available Endpoints

- `GET /health` - Health check endpoint
- `GET /api/v1/time` - Get current time
- `GET /api/v1/timezones` - List available timezones
- `GET /api/v1/time/:timezone` - Get time in specific timezone
- `POST /api/v1/time/convert` - Convert time between timezones
- `GET /ws/time` - WebSocket endpoint for real-time time updates

## Configuration

Create a `.env` file in the project root:

```env
# Server
PORT=8080
HOST=0.0.0.0
PREFORK=false

# CORS
ALLOWED_ORIGINS=*
ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
ALLOWED_HEADERS=Origin,Content-Type,Accept,Authorization
ALLOW_CREDENTIALS=true
MAX_AGE=3600

# WebSocket
WS_PING_INTERVAL=30
WS_PONG_WAIT=60
WS_WRITE_WAIT=10

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
LOG_FILE=server.log
```

## Project Structure

```
.
├── config/          # Configuration loading
├── handlers/        # HTTP request handlers
├── router/          # Route definitions
├── static/          # Static files (embedded)
├── main.go          # Application entry point
├── Dockerfile       # Multi-stage Docker build
└── docker-compose.yml
```

## Building

### Local Binary

```bash
go build -ldflags="-s -w" -o gotimedate .
./gotimedate
```

### Docker Image

```bash
docker build -t gotimedate:local .
```

## Automated CI/CD

The project includes GitHub Actions workflow that automatically:
- Builds Docker images on push to main/master
- Publishes to GitHub Container Registry (GHCR)
- Creates versioned tags from git tags (e.g., `v1.0.0`)
- Supports multi-platform builds (amd64, arm64)

## WebSocket Usage

Connect to the WebSocket endpoint for real-time time updates:

```javascript
const ws = new WebSocket('ws://localhost:8080/ws/time');

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Current time:', data.time);
};
```

## Docker Deployment

See [README.Docker.md](README.Docker.md) for comprehensive Docker deployment guide including:
- Quick start with Docker and Docker Compose
- Environment variable configuration
- Production deployment examples
- Multi-platform support
- Troubleshooting

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License
