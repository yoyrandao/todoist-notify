package commands

import (
	"context"
	"fmt"
	"log/slog"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"todoist-tg/internal/abstractions"
	"todoist-tg/internal/container"
)

var _ abstractions.CommandHandler = (*UnknownCommandHandler)(nil)

type UnknownCommandHandler struct {
	container *container.Container
}

func NewUnknownCommandHandler(container *container.Container) *UnknownCommandHandler {
	return &UnknownCommandHandler{container}
}

func (h *UnknownCommandHandler) Handle(_ context.Context, update telegram.Update) error {
	slog.Debug(fmt.Sprintf("unknown command: %s", update.Message.Command()))

	unknownCommandMessage := telegram.NewMessage(
		update.Message.Chat.ID,
		fmt.Sprintf("unknown command: %s", update.Message.Command()),
	)
	_, err := h.container.Api.Send(unknownCommandMessage)
	return err
}
