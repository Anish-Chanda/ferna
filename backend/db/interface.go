package db

import (
	"context"

	"github.com/anish-chanda/ferna/model"
)

type Database interface {
	// ---- auth stuff ----

	// CheckIfEmailExists returns true if a user record with the given email exists.
	CheckIfEmailExists(ctx context.Context, email string) (bool, error)
	// CreateUser inserts a new user and returns the new user ID.
	CreateUser(ctx context.Context, user *model.User) (int64, error)
	// GetUserByEmail fetches a user by email.
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	// GetUserByID fetches a user by ID.
	GetUserByID(ctx context.Context, userID int64) (*model.User, error)

	// ---- species stuff ----

	// SearchSpecies searches for species by common or scientific name.
	SearchSpecies(ctx context.Context, query string, limit, offset int) ([]*model.Species, error)
	// GetSpeciesByID gets species by ID
	GetSpeciesByID(ctx context.Context, speciesID int64) (*model.Species, error)

	// ---- location stuff ----

	// CreateLocation creates a new location for a user
	CreateLocation(ctx context.Context, location *model.Location) (int64, error)
	// ListLocations returns all locations for a user
	ListLocations(ctx context.Context, userID int64) ([]*model.Location, error)
	// GetLocationByID gets a location by ID (with user check)
	GetLocationByID(ctx context.Context, userID, locationID int64) (*model.Location, error)
	// UpdateLocation updates a location
	UpdateLocation(ctx context.Context, location *model.Location) error
	// DeleteLocation deletes a location
	DeleteLocation(ctx context.Context, userID, locationID int64) error

	// ---- user_plants stuff ----

	// CreateUserPlant inserts a new user plant, returning its new ID.
	CreateUserPlant(ctx context.Context, plant *model.UserPlant) (int64, error)
	// GetUserPlantByID fetches a single user plant by user+plant ID.
	GetUserPlantByID(ctx context.Context, userID, plantID int64) (*model.UserPlant, error)
	// ListUserPlants returns a paginated list of a user's plants.
	ListUserPlants(ctx context.Context, userID int64, limit, offset int) ([]*model.UserPlant, error)
	// UpdateUserPlant updates all mutable fields on the given plant.
	UpdateUserPlant(ctx context.Context, plant *model.UserPlant) error
	// DeleteUserPlant removes a plant by user+plant ID.
	DeleteUserPlant(ctx context.Context, userID, plantID int64) error

	// ---- plant_tasks stuff ----

	// CreatePlantTask creates a new plant task
	CreatePlantTask(ctx context.Context, task *model.PlantTask) (int64, error)
	// GetPlantTasksByPlantID gets all tasks for a plant
	GetPlantTasksByPlantID(ctx context.Context, userID, plantID int64) ([]*model.PlantTask, error)
	// GetPlantTaskByID gets a specific plant task
	GetPlantTaskByID(ctx context.Context, userID, taskID int64) (*model.PlantTask, error)
	// UpdatePlantTask updates a plant task
	UpdatePlantTask(ctx context.Context, task *model.PlantTask) error
	// DeletePlantTask deletes a plant task
	DeletePlantTask(ctx context.Context, userID, taskID int64) error
	// GetOverdueTasks gets all overdue tasks for a user
	GetOverdueTasks(ctx context.Context, userID int64) ([]*model.PlantTask, error)

	// ---- care_events stuff ----

	// CreateCareEvent creates a new care event
	CreateCareEvent(ctx context.Context, event *model.CareEvent) (int64, error)
	// GetCareEventsByPlantID gets all care events for a plant
	GetCareEventsByPlantID(ctx context.Context, userID, plantID int64, limit, offset int) ([]*model.CareEvent, error)
	// GetCareEventByID gets a specific care event
	GetCareEventByID(ctx context.Context, userID, eventID int64) (*model.CareEvent, error)
	// UpdateCareEvent updates a care event
	UpdateCareEvent(ctx context.Context, event *model.CareEvent) error
	// DeleteCareEvent deletes a care event
	DeleteCareEvent(ctx context.Context, userID, eventID int64) error

	// ---- other stuff -----

	// Connect establishes a connection to the database.
	Connect(dsn string) error
	// Close closes the database connection.
	Close() error
	// Migrate runs database migrations using go-migrate and embedded sql files to ensure the schema is up-to-date.
	Migrate() error
}
