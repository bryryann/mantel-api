package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Like struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	PostID    int64     `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}

type LikePublic struct {
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type LikeModel struct {
	DB *sql.DB
}

func (m *LikeModel) Like(userID, postID int64) (*Like, error) {
	query := `
		INSERT INTO likes (user_id, post_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, post_id) DO NOTHING
		RETURNING id, created_at
	`

	args := []any{userID, postID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	like := Like{
		UserID: userID,
		PostID: postID,
	}
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&like.ID, &like.CreatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, nil
		}
		return nil, err
	}

	return &like, nil
}

func (m *LikeModel) Dislike(userID, postID int64) error {
	query := `
		DELETE FROM likes
		WHERE user_id = $1 AND post_id = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{userID, postID}

	res, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (m *LikeModel) ListLikesFromPost(
	postID int64,
	pagination Pagination,
) ([]LikePublic, error) {
	var sortColumn string
	switch strings.ToLower(pagination.Sort) {
	case "asc", "oldest", "old":
		sortColumn = "created_at ASC"
	case "desc", "newest", "new":
		sortColumn = "created_at DESC"
	default:
		sortColumn = "created_at DESC"
	}

	query := fmt.Sprintf(`
		SELECT user_id, created_at
		FROM likes
		WHERE post_id = $1
		ORDER BY %s
		LIMIT $2 OFFSET $3
	`, sortColumn)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{postID, pagination.PageSize, pagination.Offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var likes []LikePublic
	for rows.Next() {
		var l LikePublic
		if err := rows.Scan(&l.UserID, &l.CreatedAt); err != nil {
			return nil, err
		}
		likes = append(likes, l)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return likes, err
}
