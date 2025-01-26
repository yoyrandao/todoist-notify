package storage

import (
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var (
	POSTGRES_CONNECT_RETRY_COUNT = 5

	POSTGRESQL_HOST          = os.Getenv("PGHOST")
	POSTGRESQL_PORT          = os.Getenv("PGPORT")
	POSTGRESQL_USERNAME      = os.Getenv("PGUSER")
	POSTGRESQL_PASSWORD      = os.Getenv("PGPASSWORD")
	POSTGRESQL_DATABASE_NAME = os.Getenv("PGDATABASE")
)

var (
	ErrNoRows = "sql: no rows in result set"
)

// GetPostgresConnectionString generates a PostgreSQL connection string.
//
// Parameters: host (string), port (string), username (string), password (string), databaseName (string).
// Returns a string.
func getPostgresConnectionString(host string, port string, username string, password string, databaseName string) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, databaseName,
	)
}

// OpenPostgres generates a PostgreSQL database connection and retries connecting a set number of times.
//
// Parameters:
//   - host: the host of the PostgreSQL database.
//   - port: the port of the PostgreSQL database.
//   - username: the username for the PostgreSQL database.
//   - password: the password for the PostgreSQL database.
//   - databaseName: the name of the PostgreSQL database.
//
// Returns:
//   - *sqlx.DB: the database connection if successful.
//   - error: an error if the connection fails after all retry attempts.
func OpenPostgres() (*sqlx.DB, error) {
	connectionString := getPostgresConnectionString(POSTGRESQL_HOST, POSTGRESQL_PORT, POSTGRESQL_USERNAME, POSTGRESQL_PASSWORD, POSTGRESQL_DATABASE_NAME)

	for attempt := 0; attempt < POSTGRES_CONNECT_RETRY_COUNT; attempt++ {
		db, err := sqlx.Connect("pgx", connectionString)

		if err == nil {
			return db, nil
		}

		time.Sleep(5 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to postgres after %d attempts", POSTGRES_CONNECT_RETRY_COUNT)
}
