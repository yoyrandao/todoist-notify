package main

import (
	"context"
	"embed"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"todoist-tg/internal/abstractions"
	"todoist-tg/internal/bootstrap"
	"todoist-tg/internal/commands"
	"todoist-tg/internal/container"
	"todoist-tg/internal/cron"
	"todoist-tg/internal/storage"
	"todoist-tg/internal/utils"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	LOG_LEVEL = utils.GetenvOrDefault("LOG_LEVEL", "debug")

	TELEGRAM_BOT_TOKEN      = os.Getenv("TELEGRAM_BOT_TOKEN")
	TELEGRAM_BOT_DEBUG_MODE = utils.GetenvOrDefault("TELEGRAM_BOT_DEBUG_MODE", "false")

	SECURITY_KEY = os.Getenv("SECURITY_KEY")
)

//go:generate cp -r ../../migrations ./migrations
//go:embed migrations/*.sql
var embedFs embed.FS

func main() {
	utils.ConfigureLogging(LOG_LEVEL)

	// configure context handler for graceful shutdown
	ctx, cancel := signal.NotifyContext(
		context.Background(), os.Interrupt, syscall.SIGTERM,
	)
	defer cancel()

	bot, err := bootstrap.NewTelegramBot(TELEGRAM_BOT_TOKEN, TELEGRAM_BOT_DEBUG_MODE)
	if err != nil {
		utils.LogFatal(err)
	}

	container := container.NewContainer(SECURITY_KEY).RegisterBotApi(bot.Api).RegisterRepositories()
	notificationSender := cron.NewNotificationSender(container)

	commandHandlerMap := map[string]abstractions.CommandHandler{
		"start":   commands.NewStartCommandHandler(container),
		"unknown": commands.NewUnknownCommandHandler(container),
	}

	if err := storage.Migrate(embedFs, "migrations"); err != nil {
		utils.LogFatal(err)
	}

	bootstrap.Schedule(ctx, notificationSender.Run, "notification-send", 5*time.Minute)

	// start handling updates from telegram
	go bot.HandleUpdates(func(u telegram.Update) {
		// if user has active handler (system waits for some input) - route to active handler
		// otherwise process standard command

		var handler abstractions.CommandHandler
		if h, exists := container.UserState.GetActiveHandler(u.Message.Chat.ID); exists {
			handler = h
		} else if h, exists := commandHandlerMap[u.Message.Command()]; exists {
			handler = h
		} else {
			handler = commandHandlerMap["unknown"]
		}

		if err := handler.Handle(ctx, u); err != nil {
			bot.Api.Send(telegram.NewMessage(u.Message.Chat.ID, "Something went wrong, try again later."))
			slog.Error("", "error", err.Error())
		}
	})

	<-ctx.Done()

	slog.Info("shutting down...")
}
