package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port            string
	Environment     string
	LogLevel        string
	ShutdownTimeout int
	
	ServiceName     string
	PodName         string
	
	DatabaseURL     string
	MaxConnections  int
	
	ApplicationInsightsKey string
	KeyVaultName           string
	ServiceBusConnection   string
	RedisConnection        string
}

func Load() (*Config, error) {
	maxConn, _ := strconv.Atoi(getEnv("MAX_CONNECTIONS", "100"))
	shutdownTimeout, _ := strconv.Atoi(getEnv("SHUTDOWN_TIMEOUT", "10"))
	
	podName := getEnv("POD_NAME", "local")
	
	appInsightsKey := getEnv("APPLICATION_INSIGHTS_KEY", "")
	if appInsightsKey == "" && fileExists("/mnt/secrets-store/appinsights-key") {
		content, err := os.ReadFile("/mnt/secrets-store/appinsights-key")
		if err == nil {
			appInsightsKey = string(content)
		}
	}
	
	return &Config{
		Port:                   getEnv("PORT", "8080"),
		Environment:            getEnv("ENVIRONMENT", "development"),
		LogLevel:               getEnv("LOG_LEVEL", "info"),
		ShutdownTimeout:        shutdownTimeout,
		ServiceName:            getEnv("SERVICE_NAME", "azure-go-app"),
		PodName:                podName,
		DatabaseURL:            getEnv("DATABASE_URL", ""),
		MaxConnections:         maxConn,
		ApplicationInsightsKey: appInsightsKey,
		KeyVaultName:           getEnv("KEY_VAULT_NAME", ""),
		ServiceBusConnection:   getEnv("SERVICE_BUS_CONNECTION", ""),
		RedisConnection:        getEnv("REDIS_CONNECTION", ""),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}