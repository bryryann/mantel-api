package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Follows struct {
	ID         int64     `json:"id"`
	FollowerID int64     `json:"follower_id"`
	FolloweeID int64     `json:"followee_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type FollowData struct {
	FollowersCount int64 `json:"followers_count"`
	FollowingCount int64 `json:"following_count"`
}

type FollowsModel struct {
	DB *sql.DB
}

// Insert adds a new follow record to the database.
func (m FollowsModel) Insert(followerID, followeeID int64) error {
	query := `
		INSERT INTO follows (follower_id, followee_id)
		VALUES ($1, $2)
		ON CONFLICT (follower_id, followee_id) DO NOTHING`

	args := []any{followerID, followeeID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...)
	if err != nil {
		// TODO: Add clearer error messages
		return err.Err()
	}

	return nil
}

// Delete removes a follow record from the follow table.
func (m FollowsModel) Delete(followerID, followeeID int64) error {
	query := `
		DELETE FROM follows
		WHERE follower_id = $1 AND followee_id = $2`

	args := []any{followerID, followeeID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...)
	if err != nil {
		return err.Err()
	}

	return nil
}

// GetFollowers returns a slice with every follower that user with related id has.
func (m FollowsModel) GetFollowers(
	userID int64,
	pagination Pagination,
) ([]UserPublic, error) {
	var sortColumn string
	switch pagination.Sort {
	case "username_asc":
		sortColumn = "u.username ASC"
	case "username_desc":
		sortColumn = "u.username DESC"
	default:
		sortColumn = "u.username ASC"
	}

	query := fmt.Sprintf(`
		SELECT u.id, u.username
		FROM follows f
		JOIN users u ON f.follower_id = u.id
		WHERE f.followee_id = $1
		ORDER BY %s
		LIMIT $2 OFFSET $3`, sortColumn)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{userID, pagination.PageSize, pagination.Offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []UserPublic
	for rows.Next() {
		var u UserPublic
		if err := rows.Scan(&u.ID, &u.Username); err != nil {
			return nil, err
		}
		followers = append(followers, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return followers, nil
}

func (m FollowsModel) GetFollowData(userID int64) (FollowData, error) {
	query := `
		SELECT
			(SELECT COUNT(*) FROM follows WHERE followee_id = $1) AS followers_count,
			(SELECT COUNT(*) FROM follows WHERE follower_id = $1) AS following_count;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var data FollowData
	err := m.DB.QueryRowContext(ctx, query, userID).Scan(&data.FollowersCount, &data.FollowingCount)
	if err != nil {
		return FollowData{}, err
	}

	return data, nil
}

func (m FollowsModel) Exists(followerID, followeeID int64) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1 FROM follows
			WHERE follower_id = $1 AND followee_id = $2
		)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{followerID, followeeID}

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// GetFollowers returns a slice with every follow by the user with given id.
func (m FollowsModel) GetFollowees(
	userID int64,
	pagination Pagination,
) ([]UserPublic, error) {
	var sortColumn string
	switch pagination.Sort {
	case "username_asc":
		sortColumn = "u.username ASC"
	case "username_desc":
		sortColumn = "u.username DESC"
	default:
		sortColumn = "u.username ASC"
	}

	query := fmt.Sprintf(`
		SELECT u.id, u.username
		FROM follows f
		JOIN users u ON f.followee_id = u.id
		WHERE f.follower_id = $1
		ORDER BY %s
		LIMIT $2 OFFSET $3`, sortColumn)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{userID, pagination.PageSize, pagination.Offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followees []UserPublic
	for rows.Next() {
		var u UserPublic
		if err := rows.Scan(&u.ID, &u.Username); err != nil {
			return nil, err
		}
		followees = append(followees, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return followees, nil
}
