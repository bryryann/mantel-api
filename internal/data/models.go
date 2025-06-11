package data

import (
	"database/sql"
	"errors"
)

// ErrRecordNotFound is an error returned when a requested record cannot be found.
var (
	ErrRecordNotFound = errors.New("record not found")
)

// Models is a container struct that holds all the individual
// database models used throughout the application.
type Models struct {
	Users   UserModel
	Follows FollowsModel
}

// NewModels initializes and returns a new Models struct,
// wiring up the database connection to each model.
func NewModels(db *sql.DB) *Models {
	return &Models{
		Users:   UserModel{DB: db},
		Follows: FollowsModel{DB: db},
	}
}
