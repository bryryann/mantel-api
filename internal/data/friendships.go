package data

import (
	"context"
	"database/sql"
	"time"
)

type FriendshipStatus string

const (
	StatusPending  FriendshipStatus = "pending"
	StatusAccepted FriendshipStatus = "accepted"
	StatusBlocked  FriendshipStatus = "blocked"
)

func (s FriendshipStatus) IsValidFriendshipStatus() bool {
	switch s {
	case StatusPending, StatusAccepted, StatusBlocked:
		return true
	default:
		return false
	}
}

type Friendship struct {
	UserID    int64            `json:"user_id"`
	FriendID  int64            `json:"friend_id"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"-"`
	Status    FriendshipStatus `json:"status"`
	Version   int              `json:"-"`
}

type FriendshipModel struct {
	DB *sql.DB
}

func (m FriendshipModel) Insert(fs *Friendship) error {
	query := `
		INSERT INTO friendships (user_id, friend_id)
		VALUES ($1, $2)
		RETURNING created_at, status
		ON CONFLICT DO NOTHING`

	args := []any{fs.UserID, fs.FriendID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&fs.CreatedAt, &fs.Status)
	if err != nil {
		return err
	}

	return nil
}
