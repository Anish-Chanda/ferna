package sqlite3

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/anish-chanda/ferna/migrations"
	"github.com/anish-chanda/ferna/model"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

type SQLiteDB struct {
	db *sql.DB
}

func NewSQLiteDB() *SQLiteDB {

	return &SQLiteDB{}
}

func (s *SQLiteDB) Connect(dsn string) error {
	var err error
	s.db, err = sql.Open("sqlite3", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	s.db.SetMaxOpenConns(1)

	// Test the connection
	if err := s.db.Ping(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	return nil
}

func (s *SQLiteDB) Close() error {
	if s.db != nil {
		if err := s.db.Close(); err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}
	}
	return nil
}

func (s *SQLiteDB) Migrate() error {
	fs, path, err := migrations.GetMigrationsFS("sqlite3")
	if err != nil {
		return fmt.Errorf("failed to get migrations fs: %w", err)
	}

	d, err := iofs.New(fs, path)
	if err != nil {
		return fmt.Errorf("failed to create source driver: %w", err)
	}

	driver, err := sqlite.WithInstance(s.db, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}
	m, err := migrate.NewWithInstance("iofs", d, "sqlite3", driver)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("Failed to migrate db: %w", err)
	}
	return nil
}

// user auth methods
func (s *SQLiteDB) CheckIfEmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)"
	err := s.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if email exists: %w", err)
	}
	return exists, nil
}

func (s *SQLiteDB) CreateUser(ctx context.Context, email, passHash string) (int64, error) {
	query := "INSERT INTO users (email, pass_hash) VALUES (?, ?)"
	result, err := s.db.ExecContext(ctx, query, email, passHash)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return result.LastInsertId()
}

func (s *SQLiteDB) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, email, pass_hash, created_at, updated_at FROM users WHERE email = ?`
	row := s.db.QueryRowContext(ctx, query, email)

	var user model.User
	var createdAt, updatedAt string

	err := row.Scan(&user.ID, &user.Email, &user.PassHash, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return nil, nil // User not found
	} else if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	var parseErr error
	user.CreatedAt, parseErr = time.Parse(time.RFC3339, createdAt)

	if parseErr != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", parseErr)
	}

	user.UpdatedAt, parseErr = time.Parse(time.RFC3339, updatedAt)
	if parseErr != nil {
		return nil, fmt.Errorf("failed to parse updated_at: %w", parseErr)
	}

	return &user, nil
}
