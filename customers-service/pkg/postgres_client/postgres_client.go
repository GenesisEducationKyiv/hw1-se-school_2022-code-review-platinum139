package postgres_client

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func NewPostgresClient(connStr string) (*sql.DB, error) {
	dbClient, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = dbClient.Ping()
	if err != nil {
		return nil, err
	}

	return dbClient, nil
}

func ClosePostgresClient(db *sql.DB) error {
	return db.Close()
}
