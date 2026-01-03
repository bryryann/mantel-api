package data

import (
	"context"
	"database/sql"
	"time"
)

type FeedModel struct {
	DB *sql.DB
}

func (m FeedModel) Fetch(
	userID int64,
	pagination Pagination,
) ([]PostPublic, error) {
	query := `
		WITH audience AS (
			SELECT $1 AS user_id
			
			UNION

			SELECT followee_id
			FROM follows
			WHERE follower_id = $1

			UNION
			
			SELECT
				CASE
					WHEN sender_id = $1 THEN receiver_id
					ELSE sender_id
				END
			FROM friendships
			WHERE status = 'accepted'
				AND ($1 IN (sender_id, receiver_id))
		)
		SELECT
			p.id,
			p.user_id,
			p.content,
			p.created_at
		FROM posts p
		JOIN audience a ON a.user_id = p.user_id
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3`

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

		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Content, &p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}
