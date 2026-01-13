package data

import (
	"context"
	"database/sql"
	"time"
)

type Bios struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
}

type BiosModel struct {
	DB *sql.DB
}

func (m BiosModel) Insert(user *User, content string) error {
	query := `
		INSERT INTO bios (user_id, content)
		VALUES ($1, $2)
		ON CONFLICT (user_id)
		DO UPDATE SET
			content = EXCLUDED.content,
			updated_at = NOW(),
			version = bios.version + 1`

	args := []any{user.ID, content}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
