package db

import (
	"context"
	"database/sql"

	"github.com/anish-chanda/ferna/model"
)

type Database interface {
	// user auth methods

	// CheckIfEmailExists returns true if a user record with the given email exists.
	// It returns (false, nil) if not found, or (false, err) on error.
	CheckIfEmailExists(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, email, passHash string) (int64, error) //returns userid
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)

	// Plant functions

	// SearchSpecies searches for species by common or scientific name.
	SearchSpecies(ctx context.Context, query string, limit, offset int) ([]*model.Species, error)

	// Database operations
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	// establishes a connection to the database.
	Connect(dsn string) error
	// Close closes the database connection.
	Close() error
	Migrate() error
}
