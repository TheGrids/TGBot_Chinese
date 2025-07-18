package config

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/joho/godotenv"
)

type Config struct {
	TGToken    string
	UseWebHook bool
	WebHookURL string
	Port       string
	LogLevel   string
}

// LoadConfig Парсит переменные окружения из .env и возвращает структуру *Config
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Переменные окружения не найдены")
	}

	useWebHook := os.Getenv("USE_WEBHOOK") == "true"

	return &Config{
		TGToken:    os.Getenv("TG_TOKEN"),
		UseWebHook: useWebHook,
		WebHookURL: os.Getenv("WEBHOOK_URL"),
		Port:       os.Getenv("PORT"),
		LogLevel:   os.Getenv("LOG_LEVEL"),
	}
}
