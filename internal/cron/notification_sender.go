package cron

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"todoist-tg/internal/container"
	"todoist-tg/internal/todoist"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NotificationSender struct {
	container *container.Container
}

func NewNotificationSender(container *container.Container) *NotificationSender {
	return &NotificationSender{container}
}

func (s *NotificationSender) Run(ctx context.Context) error {
	expired, err := s.container.NotificationTaskRepository.GetAllExpired(ctx, time.Now().UTC())
	if err != nil {
		return err
	}

	for _, task := range expired {
		user, err := s.container.UserRepository.GetByChatId(ctx, task.ChatId)
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		accessToken, err := s.container.Encryptor.Decrypt(user.EncryptedTodoistAccessToken)
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		client := todoist.NewTodoistClient(accessToken)

		todoistTask, err := client.GetTask(task.TaskId)
		if err != nil {
			slog.Error(fmt.Errorf("failed to get todoist task %s, deleting notification: %w", task.Id, err).Error())
			s.container.NotificationTaskRepository.Delete(ctx, task)
			continue
		}

		message := fmt.Sprintf("task: %s\n\n%s", todoistTask.Content, todoistTask.Description)

		if _, err := s.container.Api.Send(telegram.NewMessage(task.ChatId, message)); err != nil {
			slog.Error(err.Error())
			continue
		}

		s.container.NotificationTaskRepository.Delete(ctx, task)
		slog.Info("notification sent to user", "chat_id", task.ChatId)
	}

	return nil
}
