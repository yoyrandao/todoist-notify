package container

import (
	"todoist-tg/internal/encryption"
	"todoist-tg/internal/state"
	"todoist-tg/internal/storage"
	"todoist-tg/internal/utils"

	"todoist-tg/internal/storage/repositories"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Container struct {
	UserState *state.UserState

	Api       *telegram.BotAPI
	Encryptor *encryption.Encryptor

	UserRepository             *repositories.UserRepository
	NotificationTaskRepository *repositories.NotificationTaskRepository
}

func NewContainer(securityKey string) *Container {
	return &Container{
		Encryptor: encryption.NewEncryptor(securityKey),
		UserState: state.NewUserState(),
	}
}

func (c *Container) RegisterBotApi(api *telegram.BotAPI) *Container {
	c.Api = api
	return c
}

func (c *Container) RegisterRepositories() *Container {
	db, err := storage.OpenPostgres()
	if err != nil {
		utils.LogFatal(err)
	}

	c.UserRepository = repositories.NewUserRepository(db)
	c.NotificationTaskRepository = repositories.NewNotificationTaskRepository(db)
	return c
}
