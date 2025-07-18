package storage

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"tgbot_chinese/internal/domain"
)

// LoadHSKWords Парсит CSV словарь, возвращает массив HSKWord и ошибку
func LoadHSKWords(filePath string) ([]domain.HSKWord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия CSV файла: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения CSV файла: %w", err)
	}

	var words []domain.HSKWord

	for i, row := range rows {
		if i == 0 {
			continue
		}

		level, err := strconv.Atoi(row[3])
		if err != nil {
			continue
		}

		word := domain.HSKWord{
			Chinese: row[0],
			Pinyin:  row[1],
			Russian: row[2],
			Level:   level,
		}

		words = append(words, word)
	}

	return words, nil
}
