package db

import (
	"context"

	"github.com/anish-chanda/ferna/model"
)

type Database interface {
	// user auth methods

	// CheckIfEmailExists returns true if a user record with the given email exists.
	// It returns (false, nil) if not found, or (false, err) on error.
	CheckIfEmailExists(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, email, passHash string) (int64, error) //returns userid
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	// UpdatePassword(userID int64, passHash string) error

	// driver funcs

	// establishes a connection to the database.
	Connect(dsn string) error
	// Close closes the database connection.
	Close() error
	Migrate() error
}
