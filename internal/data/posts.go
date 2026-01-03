package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/bryryann/mantel/backend/internal/mapper"
	"github.com/bryryann/mantel/backend/internal/validator"
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

func (p Post) ToPublic() any {
	return PostPublic{
		ID:        p.ID,
		UserID:    p.UserID,
		Content:   p.Content,
		CreatedAt: p.CreatedAt,
	}
}

type PostPublic struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
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

func (m PostModel) Delete(postID int64) error {
	query := `
		DELETE FROM posts
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, postID)
	if err != nil {
		return err.Err()
	}

	return nil
}

func (m PostModel) SelectAllFromUser(
	userID int64,
	pagination Pagination,
) ([]PostPublic, error) {
	var sortColumn string
	switch strings.ToLower(pagination.Sort) {
	case "asc", "oldest", "old":
		sortColumn = "created_at ASC"
	case "desc", "newest", "new":
		sortColumn = "created_at DESC"
	default:
		sortColumn = "created_at ASC"
	}

	query := fmt.Sprintf(`
		SELECT id, user_id, content, created_at
		FROM posts
		WHERE user_id = $1
		ORDER BY %s
		LIMIT $2 OFFSET $3
	`, sortColumn)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{userID, pagination.PageSize, pagination.Offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []PostPublic
	for rows.Next() {
		var p PostPublic
		if err := rows.Scan(&p.ID, &p.UserID, &p.Content, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (m PostModel) FindByIDFromUser(postID, userID int64) (*Post, error) {
	query := `
		SELECT id, user_id, content, created_at, updated_at, version
		FROM posts 
		WHERE id = $1 AND user_id = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{postID, userID}

	var post Post
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&post.ID,
		&post.UserID,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
	)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (m PostModel) PatchPost(post *Post) (*Post, error) {
	query := `
		UPDATE posts
		SET content = $3, updated_at = $4
		WHERE id = $1 AND user_id = $2
		RETURNING updated_at
	`

	args := []any{
		post.ID,
		post.UserID,
		post.Content,
		time.Now(),
	}

	patched := *post
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&patched.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &patched, nil
}

func (m PostModel) CheckPostOwnership(postID, userID int64) (bool, error) {
	var exists bool
	query := `SELECT 1 FROM posts WHERE id = $1 AND user_id = $2 LIMIT 1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{postID, userID}

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (m PostModel) Exists(postID int64) (bool, error) {
	checkQuery := `
		SELECT COUNT(*) FROM posts
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var count int
	err := m.DB.QueryRowContext(ctx, checkQuery, postID).Scan(&count)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

func ValidatePost(v *validator.Validator, post *Post) {
	v.Check(post.Content != "", "content", "must be provided")
	v.Check(len(post.Content) <= 500, "content", "must be no more than 500 bytes long")
}
