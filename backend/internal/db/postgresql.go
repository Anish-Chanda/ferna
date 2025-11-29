package db

import (
	"context"
	"fmt"
	"time"

	"github.com/anish-chanda/ferna/internal/logger"
	"github.com/anish-chanda/ferna/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresDB wraps pgxpool.Pool with additional methods
type PostgresDB struct {
	Pool   *pgxpool.Pool
	logger *logger.ServiceLogger
}

// Config represents database configuration
type Config struct {
	DSN             string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

// NewPostgresDB creates a new PostgreSQL database connection pool and pings the database
func NewPostgresDB(ctx context.Context, cfg Config, logger *logger.ServiceLogger) (*PostgresDB, error) {
	// Parse the DSN and configure the pool
	poolConfig, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database DSN: %w", err)
	}

	// Configure connection pool settings
	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns
	poolConfig.MaxConnLifetime = cfg.MaxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime

	// Create the connection pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("PostgreSQL connection pool initialized successfully")
	logger.Debugf("Pool configuration - MaxConns: %d, MinConns: %d", cfg.MaxConns, cfg.MinConns)

	return &PostgresDB{
		Pool:   pool,
		logger: logger,
	}, nil
}

// Close gracefully closes the database connection pool
func (db *PostgresDB) Close() {
	if db.Pool != nil {
		db.logger.Info("Closing PostgreSQL connection pool")
		db.Pool.Close()
	}
}

// Stats returns connection pool statistics
func (db *PostgresDB) Stats() *pgxpool.Stat {
	return db.Pool.Stat()
}

// CheckIfEmailExists returns true if a user record with the given email exists
func (db *PostgresDB) CheckIfEmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	err := db.Pool.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		db.logger.Debugf("Failed to check email existence for %s: %v", email, err)
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}
	db.logger.Debugf("Email existence check for %s: %v", email, exists)
	return exists, nil
}

// CreateUser inserts a new user and returns the user ID
func (db *PostgresDB) CreateUser(ctx context.Context, user *model.User) (uuid.UUID, error) {
	query := `INSERT INTO users (id, avatar_url, auth_provider, email, full_name, password_hash, timezone, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
			  RETURNING id`
			  
	var userID uuid.UUID
	err := db.Pool.QueryRow(ctx, query,
		user.ID,
		user.AvatarURL,
		user.AuthProvider,
		user.Email,
		user.FullName,
		user.PasswordHash,
		user.Timezone,
	).Scan(&userID)
	
	if err != nil {
		db.logger.Debugf("Failed to create user %s: %v", user.Email, err)
		return uuid.Nil, fmt.Errorf("failed to create user: %w", err)
	}
	
	db.logger.Infof("User created successfully: %s", user.Email)
	return userID, nil
}

// GetUserByEmail fetches a user by email
func (db *PostgresDB) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, avatar_url, auth_provider, created_at, email, full_name, password_hash, timezone, updated_at
			  FROM users WHERE email = $1`
			  
	var user model.User
	err := db.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.AvatarURL,
		&user.AuthProvider,
		&user.CreatedAt,
		&user.Email,
		&user.FullName,
		&user.PasswordHash,
		&user.Timezone,
		&user.UpdatedAt,
	)
	
	if err != nil {
		if err.Error() == "no rows in result set" {
			db.logger.Debugf("User not found: %s", email)
			return nil, nil
		}
		db.logger.Debugf("Failed to get user by email %s: %v", email, err)
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	
	db.logger.Debugf("User retrieved successfully: %s", email)
	return &user, nil
}
