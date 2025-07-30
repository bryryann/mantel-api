package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrPostNotFound = errors.New("post not found")
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

func (m PostModel) Get(id int64) (*Post, error) {
	query := `
		SELECT user_id, content, created_at, updated_at, version
		FROM posts
		WHERE id = $1
	`

	post := Post{ID: id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&post.UserID,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (m PostModel) Insert(post *Post) error {
	query := `
		INSERT INTO posts (user_id, content)
		VALUES ($1, $2)
		RETURNING id, created_at`

	args := []any{post.UserID, post.Content}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&post.ID, &post.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}
