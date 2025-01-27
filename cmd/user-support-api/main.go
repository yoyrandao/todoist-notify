package main

import (
	"embed"
	"os"
	"todoist-tg/internal/bootstrap"
	"todoist-tg/internal/container"
	"todoist-tg/internal/controllers"
	"todoist-tg/internal/storage"
	"todoist-tg/internal/utils"

	"github.com/gin-gonic/gin"
)

var (
	LOG_LEVEL = utils.GetenvOrDefault("LOG_LEVEL", "debug")
	PORT      = utils.GetenvOrDefault("PORT", "8080")

	TELEGRAM_BOT_TOKEN      = os.Getenv("TELEGRAM_BOT_TOKEN")
	TELEGRAM_BOT_DEBUG_MODE = utils.GetenvOrDefault("TELEGRAM_BOT_DEBUG_MODE", "false")

	SECURITY_KEY = os.Getenv("SECURITY_KEY")
)

//go:generate cp -r ../../migrations ./migrations
//go:embed migrations/*.sql
var embedFs embed.FS

func main() {
	utils.ConfigureLogging(LOG_LEVEL)

	bot, err := bootstrap.NewTelegramBot(TELEGRAM_BOT_TOKEN, TELEGRAM_BOT_DEBUG_MODE)
	if err != nil {
		utils.LogFatal(err)
	}

	container := container.NewContainer(SECURITY_KEY).RegisterBotApi(bot.Api).RegisterRepositories()
	authzController := controllers.NewAuthorizationController(container)

	server := bootstrap.NewHTTPHost().WithRouting(func(e *gin.Engine) {
		e.GET("/todoist/authorize", authzController.Authorize)
		e.GET("/todoist/authorize-callback", authzController.AuthorizeCallback)
	}, false)

	if err := storage.Migrate(embedFs, "migrations"); err != nil {
		utils.LogFatal(err)
	}

	if err := server.Run(PORT); err != nil {
		utils.LogFatal(err)
	}
}
