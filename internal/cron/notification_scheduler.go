package cron

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"
	"todoist-tg/internal/container"
	"todoist-tg/internal/storage"
	"todoist-tg/internal/todoist"
)

type NotificationScheduler struct {
	container *container.Container
}

func NewNotificationScheduler(container *container.Container) *NotificationScheduler {
	return &NotificationScheduler{container}
}

func (s *NotificationScheduler) Run(ctx context.Context) error {
	slog.Info("starting notifications scheduler...")

	users, err := s.container.UserRepository.GetAll(ctx)
	if err != nil {
		return err
	}

	for _, user := range users {
		accessToken, err := s.container.Encryptor.Decrypt(user.EncryptedTodoistAccessToken)
		if err != nil {
			return err
		}

		tasks, err := todoist.NewTodoistClient(accessToken).GetTasks()
		if err != nil {
			return err
		}

		for _, task := range tasks {
			if task.IsCompleted {
				continue
			}

			slog.Debug("processing task", "task_id", task.Id)

			label, ok := containsNotifyLabel(task.Labels)
			if !ok {
				continue
			}

			period := getPeriodFromNotifyLabel(label)
			soonest, err := s.container.NotificationTaskRepository.GetSoonestByTask(task.Id)

			if err == nil {
				if soonest.NotificationPeriod != period {
					if err := s.container.NotificationTaskRepository.DeleteAllByTask(ctx, task.Id); err != nil {
						return err
					}
				}
			} else if err.Error() != storage.ErrNoRows {
				return err
			}

			dates, err := generateNotificationDates(time.Now(), task.Due.Date.Time, period)
			if err != nil {
				return err
			}

			for _, date := range dates {
				if err := s.container.NotificationTaskRepository.CreateOrUpdate(&storage.NotificationTask{
					Id:                 fmt.Sprintf("%s-%s", task.Id, formatTimeToIndex(date)),
					TaskId:             task.Id,
					ChatId:             user.ChatId,
					NotificationPeriod: period,
					SendThreshold:      date,
				}); err != nil {
					slog.Error("cannot create or update notification task", "error", err.Error())
				}
			}
		}
	}

	return nil
}

func containsNotifyLabel(labels []string) (string, bool) {
	for _, label := range labels {
		if strings.Contains(label, "tg-notify-") {
			return label, true
		}
	}

	return "", false
}

func getPeriodFromNotifyLabel(label string) string {
	label = strings.TrimPrefix(label, "tg-notify-")
	return label
}

func formatTimeToIndex(t time.Time) string {
	return t.Format("02012006-1500")
}

func generateNotificationDates(now, endDate time.Time, period string) ([]time.Time, error) {
	nowUTC := now.UTC()
	endUTC := endDate.UTC()
	var result []time.Time

	if nowUTC.After(endUTC) {
		return []time.Time{}, nil
	}

	// Parse the period
	unit := period[len(period)-1:]
	value, err := strconv.Atoi(period[:len(period)-1])
	if err != nil || value <= 0 {
		return nil, errors.New("invalid period format")
	}

	var nextDate time.Time
	switch unit {
	case "d":
		// Start from the current day at 12:00
		nextDate = nowUTC.Truncate(time.Hour)
		for nextDate.Before(nowUTC) {
			nextDate = nextDate.AddDate(0, 0, value)
		}
		for nextDate.Before(endUTC) || nextDate.Equal(endUTC) {
			result = append(result, nextDate)
			nextDate = nextDate.AddDate(0, 0, value)
		}
	case "h":
		// Start from the closest rounded hour
		nextDate = nowUTC.Truncate(time.Hour)
		if nextDate.Before(nowUTC) {
			nextDate = nextDate.Add(time.Hour)
		}
		for nextDate.Before(endUTC) || nextDate.Equal(endUTC) {
			result = append(result, nextDate)
			nextDate = nextDate.Add(time.Duration(value) * time.Hour)
		}
	case "M":
		// Start from the 1st of the current month at 12:00
		nextDate = time.Date(nowUTC.Year(), nowUTC.Month(), 1, 12, 0, 0, 0, time.UTC)
		for nextDate.Before(nowUTC) {
			nextDate = nextDate.AddDate(0, value, 0)
		}
		for nextDate.Before(endUTC) || nextDate.Equal(endUTC) {
			result = append(result, nextDate)
			nextDate = nextDate.AddDate(0, value, 0)
		}
		// Add the end of the range if not already included
		lastMonth := time.Date(endUTC.Year(), endUTC.Month(), 1, 12, 0, 0, 0, time.UTC)
		if !containsDate(result, lastMonth) && (lastMonth.Before(endUTC) || lastMonth.Equal(endUTC)) {
			result = append(result, lastMonth)
		}
	default:
		return nil, errors.New("unsupported period unit, use 'd', 'h', or 'M'")
	}

	// add the end date
	result = append(result, endUTC)

	return result, nil
}

func containsDate(dates []time.Time, date time.Time) bool {
	for _, d := range dates {
		if d.Equal(date) {
			return true
		}
	}
	return false
}
