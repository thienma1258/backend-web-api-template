package api

import (
	"os"
	"strconv"
	"time"
)

type AppConfig struct {
	RedisConfig             *RedisConfig
	ServiceName             string
	Port                    string
	QuoteAPIKey             string
	GeoApiKey               string
	RecordVisitAfterMinutes int
	LogMode                 LOG_MODE
}

// ServerOption contains some basic configuration for HTTP server
type ServerOption struct {
	ServeTLS        bool
	CertFile        string
	KeyFile         string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

func ternary(originalValue string, defaultValue string) string {
	if len(originalValue) == 0 {
		return defaultValue
	}
	return originalValue
}

func NewAppConfig() *AppConfig {
	config := &AppConfig{}
	config.loadRedisConfig()
	config.loadGlobalConfig()
	return config
}

func (config *AppConfig) loadRedisConfig() {
	redisAdd := ternary(os.Getenv("REDIS_ADDR"), "redis-13423.c12.us-east-1-4.ec2.cloud.redislabs.com:13423")
	pass := ternary(os.Getenv("REDIS_PASS"), "yMmzP3FGlfay2PzlwnpsUwSqD7pvZ1d4")
	config.RedisConfig = &RedisConfig{
		Addr: redisAdd,
		Pass: pass,
		DB:   0,
	}
}

func (config *AppConfig) loadGlobalConfig() {
	config.Port = ternary(os.Getenv("PORT"), "8080")
	config.ServiceName = ternary(os.Getenv("SERVICE_NAME"), "dong-pham-api-challenge")
	config.GeoApiKey = ternary(os.Getenv("GEO_API_KEY"), "")
	config.RecordVisitAfterMinutes, _ = strconv.Atoi(ternary(os.Getenv("RECORD_VISIT_AFTER_MINUTE"), "10"))

}
