package abstractions

import (
	"context"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandHandler interface {
	Handle(context.Context, telegram.Update) error
}
