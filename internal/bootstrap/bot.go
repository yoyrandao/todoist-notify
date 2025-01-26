package bootstrap

import (
	"fmt"
	"log/slog"
	"strconv"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	Api *telegram.BotAPI

	token       string
	debugEnable bool
}

func NewTelegramBot(token string, debugEnable string) (*TelegramBot, error) {
	bot, err := telegram.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	if enabled, err := strconv.ParseBool(debugEnable); err != nil && enabled {
		bot.Debug = true
	}

	slog.Info(fmt.Sprintf("running bot as account %s", bot.Self.UserName))

	return &TelegramBot{Api: bot}, nil
}

func (b *TelegramBot) HandleUpdates(handlerFunc func(telegram.Update)) {
	u := telegram.NewUpdate(0)
	u.Timeout = 60

	for update := range b.Api.GetUpdatesChan(u) {
		// ignore empty-message updates hacks
		if update.Message == nil {
			continue
		}

		// ignore empty messages
		if update.Message.Text == "" {
			continue
		}

		if update.Message.IsCommand() {
			slog.Debug(fmt.Sprintf("received command: %s from %d", update.Message.Text, update.Message.Chat.ID))
		}

		handlerFunc(update)
	}
}
