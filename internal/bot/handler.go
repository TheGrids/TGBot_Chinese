package bot

import (
	"tgbot_chinese/internal/domain"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

// HandleUpdates Обрабатывает команды телеграм бота
func HandleUpdates(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel, words []domain.HSKWord) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			log.Info().Str("command", update.Message.Command()).Msg("Введена команда")
			switch update.Message.Command() {
			case "word":
				text := GetRandomWord(words)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
				bot.Send(msg)
			case "start":
				text := "Добро пожаловать в HSK Dealer!🎉\n\n" +
					"Данный бот поможет тебе подготовиться к HSK📚"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
				bot.Send(msg)
			default:
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда"))
			}
		}
	}
}
