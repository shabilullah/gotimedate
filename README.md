# Go TimeDate API

A high-performance RESTful API for time operations with real-time WebSocket streaming. Built with Go and the Fiber framework.

> [!IMPORTANT]
> **Disclaimer:** This software is created and tailored specifically for my own personal usage. No support, maintenance, or guarantees are provided. Use it at your own risk.

## Features
- **REST API**: Comprehensive timezone-aware time operations.
- **WebSocket Streaming**: Real-time clock updates with customizable formats.
- **Interactive UI**: Premium dark-mode test interface at `/ws-test`.
- **Auto-Config**: Self-generates configuration on first run.
- **API Docs**: Integrated Swagger UI documentation.

## Getting Started

### Prerequisites
- Go 1.24+

### Installation & Run
```bash
# Clone and enter the project
git clone https://github.com/shabilullah/gotimedate.git
cd gotimedate

# Install dependencies and run
go mod tidy
go run main.go
```
The server starts at `http://localhost:8080`.

### Build Instructions

#### For Windows
```powershell
go build -o build/gotimedate.exe main.go
```

#### For Linux (cross-compile from Windows)
```powershell
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o build/gotimedate main.go
```

#### For Linux/debian (native)
```bash
go build -o build/gotimedate main.go
```

### Linux Service Installation & Updates (Debian/Ubuntu)

The fastest way to install or **update** the API as a background systemd service is using this one-liner:

```bash
curl -sSL https://raw.githubusercontent.com/shabilullah/gotimedate/master/scripts/install.sh | sudo bash
```

Alternatively, if you have already cloned the repository:

```bash
sudo ./scripts/install.sh
```

This script will:
1. Verify OS compatibility (Ubuntu/Debian).
2. Install dependencies (`git`, `golang-go`).
3. **Download/Update**: Clones the repo to `/opt/gotimedate` or pulls the latest `master`.
4. **Safety**: Stops the service automatically before recompiling.
5. **Build**: Rebuilds the binary natively from source.
6. **Persistence**: Configures the systemd service and restarts it.

#### Service Management
```bash
# Check status
sudo systemctl status gotimedate

# Start/Stop/Restart
sudo systemctl start gotimedate
sudo systemctl stop gotimedate
sudo systemctl restart gotimedate
```

#### Uninstallation
To completely remove the service, user, and all data:
```bash
# If using the one-liner from GitHub
curl -sSL https://raw.githubusercontent.com/shabilullah/gotimedate/master/scripts/uninstall.sh | sudo bash
```
```bash
# Or if you have the repo locally
sudo ./scripts/uninstall.sh
```

## API Reference

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/health` | Service health & version |
| `GET` | `/api/v1/time` | Current time (default UTC) |
| `GET` | `/api/v1/time/{tz}` | Time for specific timezone |
| `GET` | `/api/v1/timezones` | List all available timezones |
| `POST` | `/api/v1/time/convert`| Convert time between zones |
| `WS` | `/ws/time` | WebSocket streaming endpoint |

## Usage Examples

### Current Time (UTC)
`GET /api/v1/time`
```json
{
  "timestamp": "2026-01-04T16:58:00Z",
  "timezone": "UTC",
  "unix": 1736035080,
  "unix_offset": 0,
  "formatted": "4:58:00 PM",
  "date": "Saturday, January 4, 2026"
}
```

### Timezone Conversion
`POST /api/v1/time/convert`
```json
{
  "from_timezone": "UTC",
  "to_timezone": "America/New_York",
  "timestamp": "2026-01-04T16:58:00Z"
}
```
**Response:**
```json
{
  "original": { 
    "timestamp": "2026-01-04T16:58:00Z",
    "timezone": "UTC", 
    "unix": 1736035080,
    "unix_offset": 0,
    "formatted": "4:58:00 PM",
    "date": "Saturday, January 4, 2026"
  },
  "converted": { 
    "timestamp": "2026-01-04T11:58:00-05:00",
    "timezone": "America/New_York", 
    "unix": 1736035080,
    "unix_offset": -18000,
    "formatted": "11:58:00 AM",
    "date": "Saturday, January 4, 2026"
  },
  "offset_hours": -5.0
}
```

### WebSocket Protocol
`WS /ws/time`

#### Subscribe/Update (Input)
Send a JSON message to change the timezone or time format:
```json
{
  "action": "subscribe",
  "timezone": "America/New_York",
  "format": "24hour"
}
```

#### Time Update (Output)
The server streams updates every second with accurate timezone conversion:
```json
{
  "type": "time_update",
  "timestamp": "2026-01-05T18:37:40+08:00",
  "data": {
    "timestamp": "2026-01-05T05:37:40-05:00",
    "timezone": "America/New_York",
    "unix": 1736073460,
    "unix_offset": -18000,
    "formatted": "05:37:40",
    "date": "Monday, January 5, 2026"
  }
}
```

## Interactive Tools
- **Live Clock Interface**: [http://localhost:8080/ws-test](http://localhost:8080/ws-test)
- **Swagger Documentation**: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

## Configuration
At runtime, `config.env` is created automatically. You can customize the behavior of the API by modifying this file.

### `config.env` Example
Configure `ALLOWED_ORIGINS` using commas to separate multiple values. It supports standard origins, wildcard subdomains, and wildcard ports for both **REST API** and **WebSocket** connections:

```env
PORT=8080
HOST=localhost

# Performance
# Enable prefork for better performance (multiple processes)
PREFORK=false

# CORS Configuration
# Single origin: https://app.example.com
# Wildcard subdomain: https://*.example.com
# Wildcard port: http://localhost:*
ALLOWED_ORIGINS=http://localhost:3000,https://*.example.com,http://localhost:*

# Logging
# Available LOG_LEVEL: debug, info, warn, error
LOG_LEVEL=info
```

### Development vs Production Configuration

- **Development (`go run`)**: Uses `.env` file if present, ignores `config.env`
- **Production (binary)**: Uses `config.env` exclusively, generates default if missing

### Error Responses
All errors return a standardized JSON format:

```json
{
  "error": true,
  "message": "Error description",
  "code": 400,
  "path": "/api/v1/time",
  "method": "GET"
}
```

## License
MIT License - see [LICENSE](LICENSE) for details.
