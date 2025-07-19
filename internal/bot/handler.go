package bot

import (
	"fmt"
	"strconv"
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

		log.Info().Str("command", update.Message.Text).Msg("Введена команда")

		if num, err := strconv.Atoi(update.Message.Text[:1]); err == nil && num >= 1 && num <= 6 {
			word := GetRandomWordByLvl(words, num)
			if word != nil {
				response := fmt.Sprintf("🔹 Слово уровня HSK %d:\n%s [%s] - %s", num, word.Chinese, word.Pinyin, word.Russian)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, response))
			} else {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Не найдено слов для этого уровня."))
			}
		} else { // TODO Временные костыли
			switch {
			case update.Message == nil:
				continue

			case update.Message.IsCommand():
				handleCommands(bot, update.Message, words)
			case update.Message.Text == "📚 Рандом":
				sendRandomWord(bot, update.Message, words)
			case update.Message.Text == "✍️ Выбрать уровень":
				changeLvl(bot, update.Message, words)
			case update.Message.Text == "❌ Выход":
				sendWelcome(bot, update.Message.Chat.ID)
			default:
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда"))
			}
		}
	}
}

func handleCommands(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, words []domain.HSKWord) {
	switch msg.Command() {
	case "start":
		sendWelcome(bot, msg.Chat.ID)
	case "word":
		sendRandomWord(bot, msg, words)
	}
}

func sendWelcome(bot *tgbotapi.BotAPI, chatID int64) {
	text := "<b>Добро пожаловать в HSK Dealer!🎉</b>\n\n" +
		"Данный бот поможет тебе подготовиться к HSK📚"
	msg := tgbotapi.NewMessage(chatID, text)

	msg.ParseMode = "HTML" // Включаем HTML-разметку

	// Создаем клавиатуру
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📚 Рандом"),
			tgbotapi.NewKeyboardButton("✍️ Выбрать уровень"),
		),
		// tgbotapi.NewKeyboardButtonRow(
		// 	tgbotapi.NewKeyboardButton("📊 Прогресс"),
		// 	tgbotapi.NewKeyboardButton("🆘 Помощь"),
		// ),
	)

	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func sendRandomWord(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, words []domain.HSKWord) {
	text := GetRandomWord(words)
	message := tgbotapi.NewMessage(msg.Chat.ID, text)
	bot.Send(message)
}

func changeLvl(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, words []domain.HSKWord) {
	message := tgbotapi.NewMessage(msg.Chat.ID, "Выберите уровень HSK:")

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("1️⃣"),
			tgbotapi.NewKeyboardButton("2️⃣"),
			tgbotapi.NewKeyboardButton("3️⃣"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("4️⃣"),
			tgbotapi.NewKeyboardButton("5️⃣"),
			tgbotapi.NewKeyboardButton("6️⃣"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("❌ Выход"),
		),
	)

	message.ReplyMarkup = keyboard
	bot.Send(message)
}
