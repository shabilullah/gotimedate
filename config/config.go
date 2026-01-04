package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"regexp"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

type Config struct {
	StaticDir string

	Port             string
	Host             string
	Prefork          bool
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	MaxAge           int
	WSPingInterval   int
	WSPongWait       int
	WSWriteWait      int
	LogLevel         string
	LogFormat        string
	LogFile          string
	OriginPatterns   []*regexp.Regexp
}

const defaultConfigContent = `# Web Server Configuration
PORT=8080
HOST=localhost

# Performance
# Enable prefork for better performance (multiple processes)
PREFORK=false

# CORS Configuration
# Examples: 
# Single: https://app.example.com
# Wildcard Subdomain: https://*.example.com
# Wildcard Port: http://localhost:*
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
ALLOWED_HEADERS=Content-Type,Authorization,X-Requested-With
ALLOW_CREDENTIALS=true
MAX_AGE=3600

# WebSocket Configuration
WS_PING_INTERVAL=30
WS_PONG_WAIT=60
WS_WRITE_WAIT=10

# Logging
# Available LOG_LEVEL: debug, info, warn, error
LOG_LEVEL=info
LOG_FORMAT=json
LOG_FILE=server.log`

func LoadConfig(defaultHTML []byte) *Config {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Error getting executable path: %v", err)
	}

	exeDir := filepath.Dir(exePath)
	isDevelopment := strings.Contains(exePath, "go-build") || strings.Contains(exePath, "Temp")
	if isDevelopment {
		_ = godotenv.Load(".env")
	}

	configPath := filepath.Join(exeDir, "config.env")
	if !isDevelopment {
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			log.Infof("Config file %s not found, creating default...", configPath)
			err := os.WriteFile(configPath, []byte(defaultConfigContent), 0644)
			if err != nil {
				log.Fatalf("Error creating default config file: %v", err)
			}
		}
		_ = godotenv.Overload(configPath)
	}

	staticDirName := getEnv("STATIC_DIR", "static")
	var staticDirPath string
	if isDevelopment {
		// When running via 'go run', use the project root's static directory
		cwd, _ := os.Getwd()
		staticDirPath = filepath.Join(cwd, staticDirName)
	} else {
		staticDirPath = filepath.Join(exeDir, staticDirName)
	}

	if _, err := os.Stat(staticDirPath); os.IsNotExist(err) {
		log.Infof("Static directory %s not found, creating...", staticDirPath)
		if err := os.MkdirAll(staticDirPath, 0755); err != nil {
			log.Fatalf("Error creating static directory: %v", err)
		}
	}

	htmlPath := filepath.Join(staticDirPath, "websocket-test.html")
	log.Infof("Ensuring latest WebSocket test file at %s", htmlPath)
	if err := os.WriteFile(htmlPath, defaultHTML, 0644); err != nil {
		log.Fatalf("Error updating websocket-test.html: %v", err)
	}

	cfg := &Config{
		StaticDir:        staticDirPath,
		Port:             getEnv("PORT", "8080"),
		Host:             getEnv("HOST", "localhost"),
		Prefork:          getEnvBool("PREFORK", false),
		AllowedOrigins:   splitEnv("ALLOWED_ORIGINS", ","),
		AllowedMethods:   splitEnv("ALLOWED_METHODS", ","),
		AllowedHeaders:   splitEnv("ALLOWED_HEADERS", ","),
		AllowCredentials: getEnvBool("ALLOW_CREDENTIALS", true),
		MaxAge:           getEnvInt("MAX_AGE", 3600),
		WSPingInterval:   getEnvInt("WS_PING_INTERVAL", 30),
		WSPongWait:       getEnvInt("WS_PONG_WAIT", 60),
		WSWriteWait:      getEnvInt("WS_WRITE_WAIT", 10),
		LogLevel:         strings.ToLower(getEnv("LOG_LEVEL", "info")),
		LogFormat:        getEnv("LOG_FORMAT", "json"),
	}

	if isDevelopment {
		cwd, _ := os.Getwd()
		cfg.LogFile = filepath.Join(cwd, getEnv("LOG_FILE", "server.log"))
	} else {
		cfg.LogFile = filepath.Join(exeDir, getEnv("LOG_FILE", "server.log"))
	}

	cfg.CompileOrigins()
	return cfg
}

func (c *Config) CompileOrigins() {
	c.OriginPatterns = nil
	for _, origin := range c.AllowedOrigins {
		if strings.Contains(origin, "*") {
			pattern := regexp.QuoteMeta(origin)
			pattern = strings.ReplaceAll(pattern, "\\*", ".*")
			pattern = "^" + pattern + "$"
			if re, err := regexp.Compile(pattern); err == nil {
				c.OriginPatterns = append(c.OriginPatterns, re)
			}
		}
	}
}

func (c *Config) IsOriginAllowed(origin string) bool {
	if origin == "" {
		return true
	}
	for _, allowed := range c.AllowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
	}
	for _, re := range c.OriginPatterns {
		if re.MatchString(origin) {
			return true
		}
	}
	return false
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func splitEnv(key, sep string) []string {
	val := getEnv(key, "")
	if val == "" {
		return []string{}
	}
	return strings.Split(val, sep)
}

func getEnvBool(key string, fallback bool) bool {
	val := getEnv(key, "")
	if val == "" {
		return fallback
	}
	return strings.ToLower(val) == "true"
}

func getEnvInt(key string, fallback int) int {
	val := getEnv(key, "")
	if val == "" {
		return fallback
	}
	var res int
	fmt.Sscanf(val, "%d", &res)
	return res
}
