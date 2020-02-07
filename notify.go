package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func notifyAdmins(message string) error {
	m := tgbotapi.NewMessage(config.Telegram.ChatID, message)
	m.ParseMode = tgbotapi.ModeMarkdown
	if tgbot == nil {
		return fmt.Errorf("tgbot is not initialized %v", tgbot)
	}
	_, err := tgbot.Send(m)
	if err != nil {
		return err
	}
	return nil
}
