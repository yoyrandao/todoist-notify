package repositories

import (
	"context"
	"todoist-tg/internal/storage"

	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*storage.User, error) {
	var users []*storage.User
	if err := r.db.SelectContext(ctx, &users, `SELECT * FROM "users"`); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) GetByChatId(ctx context.Context, chatId int64) (*storage.User, error) {
	var user storage.User
	if err := r.db.GetContext(ctx, &user, `SELECT * FROM "users" WHERE chat_id = $1 LIMIT 1`, chatId); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) CreateOrUpdate(ctx context.Context, user *storage.User) (*storage.User, error) {
	// check if user already exists
	var existing *storage.User
	if e, err := r.GetByChatId(ctx, user.ChatId); err == nil {
		existing = e
	}

	// update existing user if it exists or create new
	var query string
	if existing != nil {
		query = `
			UPDATE "users" SET 
				encrypted_todoist_access_token = $2, 
				modified_at = now()
			WHERE chat_id = $1
		`
	} else {
		query = `
			INSERT INTO "users" (
				chat_id, encrypted_todoist_access_token, created_at, modified_at
			) VALUES ($1, $2, now(), now())
		`
	}

	_, err := r.db.ExecContext(ctx, query, user.ChatId, user.EncryptedTodoistAccessToken)
	if err != nil {
		return nil, err
	}

	return r.GetByChatId(ctx, user.ChatId)
}
