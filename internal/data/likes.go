package data

import (
	"database/sql"
	"time"
)

type Like struct {
	Id        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	PostID    int64     `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}

type LikeModel struct {
	DB *sql.DB
}
