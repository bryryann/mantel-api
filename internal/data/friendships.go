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
	ErrInvalidFriendshipStatus = errors.New("given status is not valid for friendships")
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

func (m FriendshipModel) GetFriends(
	userID int64,
	pagination Pagination,
) ([]UserPublic, error) {
	query := `
		SELECT u.id, u.username
		FROM friendships f
		JOIN users u ON u.id =
			CASE
				WHEN f.sender_id = $1 THEN receiver_id
				ELSE f.sender_id
			END
		WHERE (f.sender_id = $1 OR f.receiver_id = $1)
			AND f.status = 'accepted'
		LIMIT $2 OFFSET $3
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{userID, pagination.PageSize, pagination.Offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []UserPublic
	for rows.Next() {
		var f UserPublic
		if err := rows.Scan(&f.ID, &f.Username); err != nil {
			return nil, err

		}
		friends = append(friends, f)
	}

	return friends, nil
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

func (m FriendshipModel) GetFriendshipStatus(userID, friendID int64) (string, error) {
	var status string
	query := `
		SELECT COALESCE(
			(
				SELECT status
				FROM friendships
				WHERE LEAST(sender_id, receiver_id)
					  = LEAST($1::int, $2::int)
				  AND GREATEST(sender_id, receiver_id)
					  = GREATEST($1::int, $2::int)
				LIMIT 1
			),
			'none'
		) AS status;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, userID, friendID).Scan(&status)
	if err != nil {
		return "", err
	}

	return status, nil
}

func (m FriendshipModel) PatchFriendship(fs *Friendship) (*Friendship, error) {
	exists, err := friendshipRequestExists(m.DB, fs)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrNoSuchRequest
	}

	query := `
		UPDATE friendships
		SET status = $1, updated_at = $3
		WHERE id = $2 AND receiver_id = $4
		RETURNING sender_id, updated_at, created_at
	`
	args := []any{
		fs.Status,
		fs.ID,
		time.Now(),
		fs.ReceiverID,
	}

	patched := *fs
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = m.DB.QueryRowContext(ctx, query, args...).Scan(&patched.SenderID, &patched.UpdatedAt, &patched.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &patched, nil

}

func (m FriendshipModel) GetSentPendingRequests(
	senderID int64,
	pagination Pagination,
) ([]Friendship, error) {
	query := `
		SELECT sender_id, receiver_id, created_at, status
		FROM friendships
		WHERE sender_id = $1 AND status = 'pending'
		LIMIT $2 OFFSET $3`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{senderID, pagination.PageSize, pagination.Offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
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

func (m FriendshipModel) GetReceivedPendingRequests(
	receiverID int64,
	pagination Pagination,
) ([]Friendship, error) {
	query := `
		SELECT sender_id, receiver_id, created_at, status
		FROM friendships
		WHERE receiver_id = $1 AND status = 'pending'
		LIMIT $2 OFFSET $3`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{receiverID, pagination.PageSize, pagination.Offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
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

func friendshipRequestExists(db *sql.DB, fs *Friendship) (bool, error) {
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
