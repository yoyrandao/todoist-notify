package main

import (
	"context"
	"embed"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todoist-tg/internal/bootstrap"
	"todoist-tg/internal/container"
	"todoist-tg/internal/cron"
	"todoist-tg/internal/storage"
	"todoist-tg/internal/utils"

	"github.com/pressly/goose/v3"
)

var (
	LOG_LEVEL    = utils.GetenvOrDefault("LOG_LEVEL", "debug")
	SECURITY_KEY = os.Getenv("SECURITY_KEY")
)

//go:generate cp -r ../../migrations ./migrations
//go:embed migrations/*.sql
var embedMigrations embed.FS

func migrate() error {
	db, err := storage.OpenPostgres()
	if err != nil {
		return err
	}

	goose.SetBaseFS(embedMigrations)
	goose.SetTableName("_migrations")
	goose.SetDialect("postgres")

	if err := goose.Up(db.DB, "migrations"); err != nil {
		return err
	}

	return nil
}

func main() {
	utils.ConfigureLogging(LOG_LEVEL)

	// configure signal handler for graceful shutdown
	ctx, cancel := signal.NotifyContext(
		context.Background(), os.Interrupt, syscall.SIGTERM,
	)
	defer cancel()

	container := container.NewContainer(SECURITY_KEY).RegisterRepositories()
	notificationScheduler := cron.NewNotificationScheduler(container)

	if err := migrate(); err != nil {
		utils.LogFatal(err)
	}

	bootstrap.Schedule(ctx, notificationScheduler.Run, "notification-scheduler", 30*time.Second)

	// wait for interrupt signal
	<-ctx.Done()

	slog.Info("shutting down...")
}
