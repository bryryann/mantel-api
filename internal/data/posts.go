package data

import (
	"database/sql"
	"time"
)

type Post struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
	Version   int       `json:"-"`
}

type PostModel struct {
	DB *sql.DB
}
