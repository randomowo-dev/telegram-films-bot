package config

import (
	"crypto/sha256"
	"hash"
)

var (
	telegramBotToken = []byte(getEnv("TELEGRAM_BOT_TOKEN", ""))
	TelegramBotToken = generateSHA256(telegramBotToken)
)

func generateSHA256(data []byte) hash.Hash {
	token := sha256.New()
	_, err := token.Write(data)
	if err != nil {
		panic(err)
	}
	return token
}
