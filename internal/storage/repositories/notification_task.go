package repositories

import (
	"context"
	"fmt"
	"time"
	"todoist-tg/internal/storage"

	"github.com/jmoiron/sqlx"
)

type NotificationTaskRepository struct {
	db *sqlx.DB
}

func (r *NotificationTaskRepository) GetDatabase() *sqlx.DB {
	return r.db
}

func (r *NotificationTaskRepository) GetSoonestByTask(id string) (*storage.NotificationTask, error) {
	query := `
		SELECT * FROM "notification_tasks" WHERE id LIKE $1 ORDER BY send_threshold LIMIT 1
	`

	var task storage.NotificationTask
	if err := r.db.Get(&task, query, fmt.Sprintf("%s%%", id)); err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *NotificationTaskRepository) GetById(id string) (*storage.NotificationTask, error) {
	var task storage.NotificationTask
	if err := r.db.Get(&task, `SELECT * FROM "notification_tasks" WHERE id = $1 LIMIT 1`, id); err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *NotificationTaskRepository) GetAllExpired(ctx context.Context, t time.Time) ([]*storage.NotificationTask, error) {
	query := `
		SELECT * FROM "notification_tasks" 
		WHERE send_threshold < $1
	`
	var tasks []*storage.NotificationTask
	if err := r.db.SelectContext(ctx, &tasks, query, t); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *NotificationTaskRepository) CreateOrUpdate(task *storage.NotificationTask) error {
	var existing *storage.NotificationTask
	if e, err := r.GetById(task.Id); err == nil {
		existing = e
	}

	var err error
	if existing != nil {
		query := `UPDATE "notification_tasks" SET
			task_id = $2,
			chat_id = $3,
			notification_period = $4,
			modified_at = now()
		WHERE id = $1`

		_, err = r.db.Exec(query,
			task.Id,
			task.TaskId,
			task.ChatId,
			task.NotificationPeriod)
	} else {
		query := `
			INSERT INTO "notification_tasks" 
			(id, task_id, chat_id, notification_period, send_threshold, created_at, modified_at)
			VALUES ($1, $2, $3, $4, $5, now(), now())
		`

		_, err = r.db.Exec(query,
			task.Id,
			task.TaskId,
			task.ChatId,
			task.NotificationPeriod,
			task.SendThreshold)
	}

	return err
}

func (r *NotificationTaskRepository) Delete(ctx context.Context, task *storage.NotificationTask) error {
	query := `
		DELETE FROM "notification_tasks" 
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, task.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *NotificationTaskRepository) DeleteAllByTask(ctx context.Context, id string) error {
	query := `
		DELETE FROM "notification_tasks" WHERE id LIKE $1
	`
	_, err := r.db.ExecContext(ctx, query, fmt.Sprintf("%s%%", id))
	if err != nil {
		return err
	}

	return nil
}

func NewNotificationTaskRepository(db *sqlx.DB) *NotificationTaskRepository {
	return &NotificationTaskRepository{db}
}
