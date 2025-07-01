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

func (s *SQLiteDB) SearchSpecies(ctx context.Context, query string, limit, offset int) ([]*model.Species, error) {
	pattern := "%" + query + "%"
	rows, err := s.db.QueryContext(ctx, `
        SELECT id, common_name, scientific_name,
               default_watering_frequency_days, created_at, updated_at
          FROM species
         WHERE common_name   LIKE ? 
            OR scientific_name LIKE ?
         ORDER BY common_name
         LIMIT ? OFFSET ?`,
		pattern, pattern, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("error searching species query: %w", err)
	}

	defer rows.Close()

	var res []*model.Species
	for rows.Next() {
		var sp model.Species
		var created, updated string
		if err := rows.Scan(&sp.ID,
			&sp.CommonName,
			&sp.ScientificName,
			&sp.DefaultWateringFrequency,
			&created,
			&updated); err != nil {
			return nil, fmt.Errorf("scan species: %w", err)
		}

		//parse time
		sp.CreatedAt, err = time.Parse(time.RFC3339, created)
		if err != nil {
			return nil, fmt.Errorf("parse created_at: %w", err)
		}
		sp.UpdatedAt, err = time.Parse(time.RFC3339, updated)
		if err != nil {
			return nil, fmt.Errorf("parse updated_at: %w", err)
		}
		res = append(res, &sp)

	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate species rows: %w", err)
	}
	return res, nil
}

func (s *SQLiteDB) GetSpeciesByID(ctx context.Context, speciesID int64) (*model.Species, error) {
	row := s.db.QueryRowContext(ctx, `
        SELECT id, common_name, scientific_name, default_watering_frequency_days, created_at, updated_at
        FROM species
        WHERE id = ?`, speciesID)

	var sp model.Species
	var created, updated sql.NullString
	if err := row.Scan(
		&sp.ID, &sp.CommonName, &sp.ScientificName, &sp.DefaultWateringFrequency,
		&created, &updated,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("GetSpeciesByID: %w", err)
	}

	sp.CreatedAt, _ = time.Parse(time.RFC3339, created.String)
	sp.UpdatedAt, _ = time.Parse(time.RFC3339, updated.String)
	return &sp, nil
}

// ExecContext executes a query with context that doesn't return rows, like insert, delete etc
func (s *SQLiteDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	return result, nil
}

// CreatePlant inserts a new plant and returns its ID.
func (s *SQLiteDB) CreatePlant(ctx context.Context, p *model.Plant) (int64, error) {
	res, err := s.db.ExecContext(ctx, `
      INSERT INTO plants (
        user_id, species_id, nickname, image_url,
        watering_frequency_days, last_watered_at, note
      ) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		p.UserID, p.SpeciesID, p.Nickname, p.ImageURL,
		p.WateringFrequencyDays, p.LastWateredAt, p.Note,
	)
	if err != nil {
		return 0, fmt.Errorf("CreatePlant: %w", err)
	}
	return res.LastInsertId()
}

// GetPlantByID fetches one plant, ensuring it belongs to the user.
func (s *SQLiteDB) GetPlantByID(ctx context.Context, userID, plantID int64) (*model.Plant, error) {
	row := s.db.QueryRowContext(ctx, `
      SELECT id, user_id, species_id, nickname, image_url,
             watering_frequency_days, last_watered_at, note,
             created_at, updated_at
        FROM plants
       WHERE id = ? AND user_id = ?`, plantID, userID)

	var p model.Plant
	var lw, created, updated sql.NullString
	if err := row.Scan(
		&p.ID, &p.UserID, &p.SpeciesID, &p.Nickname, &p.ImageURL,
		&p.WateringFrequencyDays, &lw, &p.Note,
		&created, &updated,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("GetPlantByID: %w", err)
	}
	if lw.Valid {
		t, err := time.Parse(time.RFC3339, lw.String)
		if err == nil {
			p.LastWateredAt = &t
		}
	}
	p.CreatedAt, _ = time.Parse(time.RFC3339, created.String)
	p.UpdatedAt, _ = time.Parse(time.RFC3339, updated.String)
	return &p, nil
}

// ListPlants returns a page of plants for a user.
func (s *SQLiteDB) ListPlants(ctx context.Context, userID int64, limit, offset int) ([]*model.Plant, error) {
	rows, err := s.db.QueryContext(ctx, `
      SELECT id, user_id, species_id, nickname, image_url,
             watering_frequency_days, last_watered_at, note,
             created_at, updated_at
        FROM plants
       WHERE user_id = ?
       ORDER BY id
       LIMIT ? OFFSET ?`, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("ListPlants: %w", err)
	}
	defer rows.Close()

	var list []*model.Plant
	for rows.Next() {
		var p model.Plant
		var lw, created, updated sql.NullString
		if err := rows.Scan(
			&p.ID, &p.UserID, &p.SpeciesID, &p.Nickname, &p.ImageURL,
			&p.WateringFrequencyDays, &lw, &p.Note,
			&created, &updated,
		); err != nil {
			return nil, fmt.Errorf("ListPlants scan: %w", err)
		}
		if lw.Valid {
			t, err := time.Parse(time.RFC3339, lw.String)
			if err == nil {
				p.LastWateredAt = &t
			}
		}
		p.CreatedAt, _ = time.Parse(time.RFC3339, created.String)
		p.UpdatedAt, _ = time.Parse(time.RFC3339, updated.String)
		list = append(list, &p)
	}
	return list, rows.Err()
}

// UpdatePlant updates an existing plant’s mutable fields.
func (s *SQLiteDB) UpdatePlant(ctx context.Context, p *model.Plant) error {
	_, err := s.db.ExecContext(ctx, `
      UPDATE plants
         SET species_id = ?, nickname = ?, image_url = ?,
             watering_frequency_days = ?, last_watered_at = ?, note = ?
       WHERE id = ? AND user_id = ?`,
		p.SpeciesID, p.Nickname, p.ImageURL,
		p.WateringFrequencyDays, p.LastWateredAt, p.Note,
		p.ID, p.UserID,
	)
	return fmt.Errorf("UpdatePlant: %w", err)
}

// DeletePlant removes a plant by its ID and owner.
func (s *SQLiteDB) DeletePlant(ctx context.Context, userID, plantID int64) error {
	_, err := s.db.ExecContext(ctx,
		`DELETE FROM plants WHERE id = ? AND user_id = ?`, plantID, userID)
	return fmt.Errorf("DeletePlant: %w", err)
}
