package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrFriendshipRequestToSelf = errors.New("cannot send friend request to yourself")
	ErrNoSuchRequest           = errors.New("the friend request does not exist")
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
	checkQuery := `
		SELECT COUNT(*) FROM friendships
		WHERE user_id = $1 AND friend_id = $2 AND status = 'pending'
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var count int
	err := m.DB.QueryRowContext(ctx, checkQuery, fs.UserID, fs.FriendID).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return ErrNoSuchRequest
	}

	query := `
		UPDATE friendships
		SET status = 'accepted', updated_at = $3
		WHERE user_id = $1 AND friend_id = $2 AND status = 'pending'
	`
	args := []any{fs.UserID, fs.FriendID, time.Now()}

	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m FriendshipModel) GetPendingRequests(id int64) ([]Friendship, error) {
	query := `
		SELECT user_id, friend_id, created_at, status
		FROM friendships
		WHERE friend_id = $1 and status = 'pending'
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []Friendship
	for rows.Next() {
		var f Friendship

		err = rows.Scan(&f.UserID, &f.FriendID, &f.CreatedAt, &f.Status)
		if err != nil {
			return nil, err
		}

		requests = append(requests, f)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}
