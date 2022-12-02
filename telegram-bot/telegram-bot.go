package telegram_bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type TelegramBot interface {
	StartBot()
	SendMessage(chatID int64, text string) error
}

type TgBot struct {
	TelegramBotApiKey string
	Bot               *tgbotapi.BotAPI
}

func (t *TgBot) StartBot() {
	bot, err := tgbotapi.NewBotAPI(t.TelegramBotApiKey)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	t.Bot = bot
}
func (t *TgBot) SendMessage(message string) {
	msg := tgbotapi.NewMessage(449202647, message)
	_, err := t.Bot.Send(msg)
	if err != nil {
		return
	}
}
