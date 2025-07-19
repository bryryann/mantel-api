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
	ID         int64            `json:"id"`
	SenderID   int64            `json:"sender_id"`
	ReceiverID int64            `json:"receiver_id"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"-"`
	Status     FriendshipStatus `json:"status"`
	Version    int              `json:"-"`
}

type FriendshipModel struct {
	DB *sql.DB
}

func (m FriendshipModel) SendRequest(fs *Friendship) error {
	if fs.SenderID == fs.ReceiverID {
		return ErrFriendshipRequestToSelf
	}

	query := `
		INSERT INTO friendships (sender_id, receiver_id)
		VALUES ($1, $2)
		RETURNING created_at, status
		ON CONFLICT (sender_id, receiver_id) DO NOTHING
	`

	args := []any{fs.SenderID, fs.ReceiverID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&fs.CreatedAt, &fs.Status)
	if err != nil {
		return err
	}

	return nil
}

func (m FriendshipModel) AcceptRequest(fs *Friendship) error {
	exists, err := requestExists(m.DB, fs)
	if err != nil {
		return err
	}

	if !exists {
		return ErrNoSuchRequest
	}

	query := `
		UPDATE friendships
		SET status = 'accepted', updated_at = $3
		WHERE sender_id = $1 AND receiver_id = $2 AND status = 'pending'
	`
	args := []any{fs.SenderID, fs.ReceiverID, time.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m FriendshipModel) RejectRequest(fs *Friendship) error {
	exists, err := requestExists(m.DB, fs)
	if err != nil {
		return err
	}

	if !exists {
		return ErrNoSuchRequest
	}

	query := `
		DELETE FROM friendships
		WHERE sender_id = $1 AND receiver_id = $2 AND status = 'pending'
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := m.DB.ExecContext(ctx, query, fs.SenderID, fs.ReceiverID)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (m FriendshipModel) GetPendingRequests(id int64) ([]Friendship, error) {
	query := `
		SELECT sender_id, receiver_id, created_at, status
		FROM friendships
		WHERE receiver_id = $1 and status = 'pending'
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

		err = rows.Scan(&f.SenderID, &f.ReceiverID, &f.CreatedAt, &f.Status)
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

func requestExists(db *sql.DB, fs *Friendship) (bool, error) {
	checkQuery := `
		SELECT COUNT(*) FROM friendships
		WHERE sender_id = $1 AND receiver_id = $2 AND status = 'pending'
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var count int
	err := db.QueryRowContext(ctx, checkQuery, fs.SenderID, fs.ReceiverID).Scan(&count)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}
	return true, nil
}
