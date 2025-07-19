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

func GetRandomWordByLvl(words []domain.HSKWord, lvl int) *domain.HSKWord {
	var filtered []domain.HSKWord

	for _, word := range words {
		if word.Level == lvl {
			filtered = append(filtered, word)
		}
	}

	if len(filtered) == 0 {
		return nil
	}

	randomIndex := rand.Intn(len(filtered))
	return &filtered[randomIndex]
}
