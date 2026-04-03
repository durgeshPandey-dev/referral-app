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

func Load() *Config {
	if err := godotenv.Load(os.Getenv("CONFIG_PATH")); err != nil {
		log.Panicf("failed to load config: %v", err)
	}

	cfg := &Config{
		Port:                  os.Getenv("PORT"),
		SendGridAPIKey:        os.Getenv("SENDGRID_API_KEY"),
		SenderEmail:           os.Getenv("SENDER_EMAIL"),
		RequestTimeoutSeconds: getInt("REQUEST_TIMEOUT_SECONDS", 30),
		WorkerCount:           getInt("WORKER_COUNT", 10),
		QueueSize:             getInt("QUEUE_SIZE", 1000),
	}

	if cfg.Port == "" {
		cfg.Port = "3000"
	}

	if cfg.SendGridAPIKey == "" || cfg.SenderEmail == "" {
		log.Fatal("Missing SENDGRID_API_KEY or SENDER_EMAIL in .env")
	}

	return cfg
}
