package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/bryryann/mantel/backend/internal/mapper"
	"github.com/bryryann/mantel/backend/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

// ErrDuplicateEmail is returned when an attempt is made to create or update
// a user with an email address that already exists in the database.
//
// AnonymousUser is a sentinel value used to represent an unauthenticated or guest user.
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrDuplicateEmail    = errors.New("duplicate email")
	ErrDuplicateUsername = errors.New("duplicate username")

	AnonymousUser = &User{}
)

// User represents a user in the system.
type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Version   int       `json:"-"`
}

// ToPublic maps a variable of type User to UserPublic.
func (u User) ToPublic() any {
	return UserPublic{
		ID:       u.ID,
		Username: u.Username,
	}
}

// UserPublic contains no sensitive information about user. Safe for public exposure.
type UserPublic struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

// IsAnonymous returns true if the user is the predefined AnonymousUser.
// This is used to check if the current user is unauthenticated.
func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

// password represents a user's password, including both the plaintext version
// (used temporarily and never stored) and the hashed version (stored securely).
type password struct {
	// plaintext holds the plaintext password temporarily.
	// This should be nil except during creation or validation.
	plaintext *string

	// hash contains the securely hashed version of the password.
	hash []byte
}

// Set generates a secure hash from the given plaintext password
// and stores both the plaintext and the hash within the struct.
func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

// Matches checks whether the provided plaintext password matches
// the stored hash in the password struct.
func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

// UserModel serves as a wrapper to a SQL database connection and provides
// methods for performing CRUD operations on user records.
type UserModel struct {
	DB *sql.DB
}

// Insert adds a new user record to the database.
func (m UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	args := []any{user.Username, user.Email, user.Password.hash}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

// Get retrieves a user from the database by their unique ID.
func (m UserModel) Get(userId int64) (*User, error) {
	query := `
		SELECT id, created_at, username, email, password_hash, version
		FROM users
		WHERE id = $1`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, userId).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

// GetByUsername retrieves a user from the database by their username.
func (m UserModel) GetByUsername(username string) (*User, error) {
	query := `
		SELECT id, created_at, username, email, password_hash, version
		FROM users
		WHERE username = $1`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

// Update modifies an existing user record in the database with new data.
func (m UserModel) Update(user *User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, password_hash = $3, version = version + 1
		WHERE id = $4 AND version = $5
		RETURNING version`

	args := []any{
		user.Username,
		user.Email,
		user.Password.hash,
		user.ID,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		// TODO: Add ErrEditConflict custom error.
		default:
			return err
		}
	}

	return nil
}

func (m UserModel) Exists(id int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM users WHERE id = $1
		)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, err
	}

	if !exists {
		return false, ErrUserNotFound
	}

	return true, nil
}

// ValidateEmail checks whether the provided email string is valid,
// using methods defined in the validator package.
func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

// ValidatePasswordPlaintext checks whether the provided plaintext password
// meets defined security requirements (e.g. length, character rules).
func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 255, "password", "must be no more than 255 bytes long")
}

// ValidateUser runs validation on a User struct, ensuring all required fields
// are present and meet defined validation rules.
func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Username != "", "username", "must be provided")
	v.Check(len(user.Username) <= 500, "username", "must be no more than 500 bytes long")

	ValidateEmail(v, user.Email)

	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}
