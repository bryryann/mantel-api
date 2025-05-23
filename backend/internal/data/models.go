package data

import "database/sql"

// Models is a container struct that holds all the individual
// database models used throughout the application.
type Models struct {
	Users UserModel
}

// NewModels initializes and returns a new Models struct,
// wiring up the database connection to each model.
func NewModels(db *sql.DB) *Models {
	return &Models{
		Users: UserModel{DB: db},
	}
}
