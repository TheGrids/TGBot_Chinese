package bot

import (
	"fmt"
	"math/rand"
	"tgbot_chinese/internal/domain"

	"github.com/rs/zerolog/log"
)

// GetRandomWord Возвращает строку со случайным иероглифом из словаря
func GetRandomWord(words []domain.HSKWord) string {
	if len(words) == 0 {
		log.Fatal().Msg("Словарь пусть")
	}
	word := words[rand.Intn(len(words))]
	return fmt.Sprintf("📘 %s\n📣 %s\n🇷🇺 %s\n📚 HSK: %d", word.Chinese, word.Pinyin, word.Russian, word.Level)
}
