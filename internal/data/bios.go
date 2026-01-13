package data

import (
	"context"
	"database/sql"
	"errors"
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

func (m BiosModel) GetByUserID(userID int64) (*Bios, error) {
	query := `
		SELECT id, user_id, content, created_at, updated_at, version
		FROM bios
		WHERE user_id = $1
	`

	var bio Bios

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, userID).Scan(
		&bio.ID,
		&bio.UserID,
		&bio.Content,
		&bio.CreatedAt,
		&bio.UpdatedAt,
		&bio.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &bio, nil
}

func (m BiosModel) Update(bio *Bios) error {
	query := `
		UPDATE bios
		SET content = $1,
			version = version + 1
		WHERE id = $2 AND version = $3
		RETURNING version`

	args := []any{bio.Content, bio.ID, bio.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&bio.Version)
	if err != nil {
		return err
	}

	return nil
}
