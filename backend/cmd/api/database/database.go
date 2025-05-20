package database

import (
	"context"
	"database/sql"
	"time"
)

// Database holds db related variables used by the application.
type Database struct {
	db *sql.DB
}

// OpenConnection initializes a connection to PostgreSQL db, as well as
// setting up db connection configurations.
func OpenConnection(dsn string) (*Database, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}
