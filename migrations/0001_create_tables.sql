-- +goose Up

CREATE TABLE IF NOT EXISTS users (
	chat_id bigint NOT NULL,
	encrypted_todoist_access_token varchar(255) NOT NULL,
	created_at timestamp NOT NULL,
	modified_at timestamp NOT NULL,

	CONSTRAINT users_pkey PRIMARY KEY (chat_id)
);

CREATE TABLE IF NOT EXISTS notification_tasks (
	id varchar(255) NOT NULL,
	task_id varchar(255) NOT NULL,
	chat_id bigint NOT NULL,
	notification_period varchar(255) NOT NULL,
	send_threshold timestamp NOT NULL,
	created_at timestamp NOT NULL,
	modified_at timestamp NOT NULL,

	CONSTRAINT notification_tasks_pkey PRIMARY KEY (id)
);

-- +goose Down

DROP TABLE IF EXISTS notification_tasks;
DROP TABLE IF EXISTS users;