package commands

import (
	"context"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"todoist-tg/internal/abstractions"
	"todoist-tg/internal/container"
	"todoist-tg/internal/messages"
	"todoist-tg/internal/storage"
)

var _ abstractions.CommandHandler = (*StartCommandHandler)(nil)

type StartCommandHandler struct {
	container *container.Container
}

func NewStartCommandHandler(container *container.Container) *StartCommandHandler {
	return &StartCommandHandler{container}
}

func (h *StartCommandHandler) Handle(ctx context.Context, update telegram.Update) error {
	if !update.Message.IsCommand() {
		return h.processInput(ctx, update)
	}

	return h.processCommand(update)
}

func (h *StartCommandHandler) processCommand(update telegram.Update) error {
	message := telegram.NewMessage(update.Message.Chat.ID, messages.Greeting)
	if _, err := h.container.Api.Send(message); err != nil {
		return err
	}

	message = telegram.NewMessage(update.Message.Chat.ID, messages.GettingStarted)
	if _, err := h.container.Api.Send(message); err != nil {
		return err
	}

	h.container.UserState.SetActiveHandler(update.Message.Chat.ID, h)
	return nil
}

func (h *StartCommandHandler) processInput(ctx context.Context, update telegram.Update) error {
	encryptedAccessToken, err := h.container.Encryptor.Encrypt(update.Message.Text)
	if err != nil {
		return err
	}

	if _, err = h.container.UserRepository.CreateOrUpdate(ctx, &storage.User{
		ChatId:                      update.Message.Chat.ID,
		EncryptedTodoistAccessToken: encryptedAccessToken,
	}); err != nil {
		return err
	}

	message := telegram.NewMessage(update.Message.Chat.ID, messages.AddingTokenSuccessed)
	if _, err := h.container.Api.Send(message); err != nil {
		return err
	}

	h.container.UserState.DeleteActiveHandler(update.Message.Chat.ID)
	return nil
}
