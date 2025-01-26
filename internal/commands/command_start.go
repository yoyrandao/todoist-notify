package commands

import (
	"context"
	"fmt"
	"os"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"todoist-tg/internal/abstractions"
	"todoist-tg/internal/container"
	"todoist-tg/internal/messages"
)

var _ abstractions.CommandHandler = (*StartCommandHandler)(nil)

type StartCommandHandler struct {
	container *container.Container
}

var (
	SUPPORT_API_URL = os.Getenv("SUPPORT_API_URL")
)

func NewStartCommandHandler(container *container.Container) *StartCommandHandler {
	return &StartCommandHandler{container}
}

func (h *StartCommandHandler) Handle(ctx context.Context, update telegram.Update) error {
	message := telegram.NewMessage(update.Message.Chat.ID, messages.Greeting)
	if _, err := h.container.Api.Send(message); err != nil {
		return err
	}

	message = telegram.NewMessage(update.Message.Chat.ID, messages.GettingStarted)
	message.ReplyMarkup = getKeyboard(update.Message.Chat.ID)
	if _, err := h.container.Api.Send(message); err != nil {
		return err
	}

	return nil
}

func getKeyboard(chatId int64) *telegram.InlineKeyboardMarkup {
	var keyboard = telegram.NewInlineKeyboardMarkup(
		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonURL("Authorize", fmt.Sprintf("%s/todoist/authorize?chat_id=%d", SUPPORT_API_URL, chatId)),
		),
	)

	return &keyboard
}
