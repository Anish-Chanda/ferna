package db

import (
	"context"
	"database/sql"

	"github.com/anish-chanda/ferna/model"
)

type Database interface {
	// ---- auth stuff ----

	// CheckIfEmailExists returns true if a user record with the given email exists.
	// It returns (false, nil) if not found, or (false, err) on error.
	CheckIfEmailExists(ctx context.Context, email string) (bool, error)
	// inserts a new user with the given email and password hash and returns the new user ID.
	CreateUser(ctx context.Context, email, passHash string) (int64, error)
	// fetches a user by email.
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)

	// ---- plant stuff ----

	// CreatePlant inserts a new plant, returning its new ID.
	CreatePlant(ctx context.Context, p *model.Plant) (int64, error)
	// GetPlantByID fetches a single plant by user+plant ID.
	GetPlantByID(ctx context.Context, userID, plantID int64) (*model.Plant, error)
	// ListPlants returns a paginated list of a user’s plants.
	ListPlants(ctx context.Context, userID int64, limit, offset int) ([]*model.Plant, error)
	// UpdatePlant updates all mutable fields on the given plant.
	UpdatePlant(ctx context.Context, p *model.Plant) error
	// DeletePlant removes a plant by userplant ID.
	DeletePlant(ctx context.Context, userID, plantID int64) error
	// SearchSpecies searches for species by common or scientific name.
	SearchSpecies(ctx context.Context, query string, limit, offset int) ([]*model.Species, error)
	// get species with a ID
	GetSpeciesByID(ctx context.Context, speciesID int64) (*model.Species, error)

	// ---- other stuff -----

	// Database operations
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	// establishes a connection to the database.
	Connect(dsn string) error
	// Close closes the database connection.
	Close() error
	//  Migrate runs database migrations using go-migrate and embedded sql files to ensure the schema is up-to-date.
	Migrate() error
}
