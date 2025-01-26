package storage

import "time"

type User struct {
	ChatId                      int64     `db:"chat_id"`
	EncryptedTodoistAccessToken string    `db:"encrypted_todoist_access_token"`
	CreatedAt                   time.Time `db:"created_at"`
	ModifiedAt                  time.Time `db:"modified_at"`
}

type NotificationTask struct {
	Id                 string    `db:"id"`
	TaskId             string    `db:"task_id"`
	ChatId             int64     `db:"chat_id"`
	NotificationPeriod string    `db:"notification_period"`
	SendThreshold      time.Time `db:"send_threshold"`
	CreatedAt          time.Time `db:"created_at"`
	ModifiedAt         time.Time `db:"modified_at"`
}
