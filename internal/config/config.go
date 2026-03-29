package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                  string
	SendGridAPIKey        string
	SenderEmail           string
	RequestTimeoutSeconds int
	WorkerCount           int
	QueueSize             int
	ObservabilityEnabled  bool
	ServiceName           string
	Environment           string
	OTLPTraceEndpoint     string
	OTLPInsecure          bool
}

func getInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}

func getBool(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}

func Load() *Config {
	_ = godotenv.Load("configs/.env")

	cfg := &Config{
		Port:                  os.Getenv("PORT"),
		SendGridAPIKey:        os.Getenv("SENDGRID_API_KEY"),
		SenderEmail:           os.Getenv("SENDER_EMAIL"),
		RequestTimeoutSeconds: getInt("REQUEST_TIMEOUT_SECONDS", 30),
		WorkerCount:           getInt("WORKER_COUNT", 10),
		QueueSize:             getInt("QUEUE_SIZE", 1000),
		ObservabilityEnabled:  getBool("OBSERVABILITY_ENABLED", true),
		ServiceName:           os.Getenv("OTEL_SERVICE_NAME"),
		Environment:           os.Getenv("DEPLOYMENT_ENV"),
		OTLPTraceEndpoint:     os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		OTLPInsecure:          getBool("OTEL_EXPORTER_OTLP_INSECURE", true),
	}

	if cfg.Port == "" {
		cfg.Port = "3000"
	}

	if cfg.ServiceName == "" {
		cfg.ServiceName = "beckn-onix-microservice"
	}

	if cfg.Environment == "" {
		cfg.Environment = "dev"
	}

	if cfg.OTLPTraceEndpoint == "" {
		cfg.OTLPTraceEndpoint = "localhost:4317"
	}

	if cfg.SendGridAPIKey == "" || cfg.SenderEmail == "" {
		log.Fatal("Missing SENDGRID_API_KEY or SENDER_EMAIL in .env")
	}

	return cfg
}
