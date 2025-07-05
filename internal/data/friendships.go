package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrFriendshipRequestToSelf = errors.New("cannot send friend request to yourself")
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

func (m FriendshipModel) SendRequest(fs *Friendship) error {
	if fs.UserID == fs.FriendID {
		return ErrFriendshipRequestToSelf
	}

	query := `
		INSERT INTO friendships (user_id, friend_id)
		VALUES ($1, $2)
		RETURNING created_at, status
		ON CONFLICT (user_id, friend_id) DO NOTHING
	`

	args := []any{fs.UserID, fs.FriendID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&fs.CreatedAt, &fs.Status)
	if err != nil {
		return err
	}

	return nil
}

func (m FriendshipModel) AcceptRequest(fs *Friendship) error {
	query := `
		UPDATE friendships
		SET status = 'accepted', updated_at = $3
		WHERE user_id = $1 AND friend_id = $2 AND status = 'pending'
	`

	args := []any{fs.UserID, fs.FriendID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...)
	if err != nil {
		return err.Err()
	}

	return nil
}
