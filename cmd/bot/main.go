package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"tgbot_chinese/internal/bot"
	"tgbot_chinese/internal/config"
	"tgbot_chinese/internal/storage"
)

func main() {
	// Инициализация переменных окружения
	cfg := config.LoadConfig()
	if err := initLogger(cfg); err != nil {
		panic("Ошибка инициализации логгера")
	}

	log.Debug().Str("config", fmt.Sprintf("%+v", cfg)).Msg("Конфигурация загружена")
	log.Info().Msg("Запуск бота...")

	// Регистрация обработчика сигналов для graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Запуск бота в отдельной горутине
	go func() {
		if err := runBot(cfg); err != nil {
			log.Error().Err(err).Msg("Ошибка в работе бота")
		}
	}()

	<-stop
	log.Info().Msg("Получен сигнал завершения...")
	log.Info().Msg("Завершение работы...")
	time.Sleep(100 * time.Millisecond)
}

// runBot Основная логика запуска бота
func runBot(cfg *config.Config) error {
	botAPI, err := tgbotapi.NewBotAPI(cfg.TGToken) // Создание бота с помощью токена
	if err != nil {
		return fmt.Errorf("ошибка создания бота: %w", err)
	}

	words, err := storage.LoadHSKWords("./assets/hsk.csv") // Инициализация массива слов
	if err != nil {
		return fmt.Errorf("ошибка загрузки слов: %w", err)
	}

	if cfg.UseWebHook { // Запуск через вебхуки
		log.Info().Str("url", cfg.WebHookURL).Msg("Запуск через Webhook")

		webhook, err := tgbotapi.NewWebhook(cfg.WebHookURL)
		if err != nil {
			return fmt.Errorf("ошибка создания WebHook конфигурации: %w", err)
		}

		_, err = botAPI.Request(webhook)
		if err != nil {
			return fmt.Errorf("ошибка установки WebHook: %w", err)
		}

		updates := botAPI.ListenForWebhook("/bot")
		log.Info().Str("port", cfg.Port).Msg("Сервер слушает порт")
		go func() {
			log.Info().Msg("Запуск HTTP сервера")
			if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
				log.Fatal().Err(err).Msg("Ошибка HTTP сервера")
			}
		}()
		bot.HandleUpdates(botAPI, updates, words)
	} else { // Запуск через Long Polling
		log.Info().Msg("Запуск через Long Polling")
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		updates := botAPI.GetUpdatesChan(u)
		bot.HandleUpdates(botAPI, updates, words)
	}

	return nil
}

// initLogger Инициализация логгера zerolog
func initLogger(cfg *config.Config) error {
	if err := os.MkdirAll("logs", 0755); err != nil {
		return fmt.Errorf("не удалось создать директорию для логов: %w", err)
	}

	logFileName := fmt.Sprintf("logs/bot_%s.log", time.Now().Format("2006-01-02"))

	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл логов: %w", err)
	}

	multiWriter := io.MultiWriter(
		zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		},
		zerolog.ConsoleWriter{
			Out:        logFile,
			NoColor:    true,
			TimeFormat: time.RFC3339,
		},
	)

	log.Logger = zerolog.New(multiWriter).
		With().
		Timestamp().
		Logger()

	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level) // Установка уровня логирования

	return nil
}
