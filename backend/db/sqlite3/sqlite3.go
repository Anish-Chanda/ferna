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
	_ "github.com/mattn/go-sqlite3"
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

	// Enable foreign keys
	if _, err := s.db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}

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
	driver, err := sqlite.WithInstance(s.db, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	fs, path, err := migrations.GetMigrationsFS("sqlite3")
	if err != nil {
		return fmt.Errorf("failed to get migrations filesystem: %w", err)
	}

	sourceDriver, err := iofs.New(fs, path)
	if err != nil {
		return fmt.Errorf("failed to create source driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "sqlite3", driver)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

// parseTime parses SQLite TEXT timestamp to time.Time
func parseTime(s string) (time.Time, error) {
	// Try RFC3339 format first (for manually inserted timestamps)
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	// Try SQLite's CURRENT_TIMESTAMP format
	if t, err := time.Parse("2006-01-02 15:04:05", s); err == nil {
		return t, nil
	}
	// Try RFC3339 without timezone
	if t, err := time.Parse("2006-01-02T15:04:05", s); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("unable to parse timestamp: %s", s)
}

// formatTime formats time.Time to SQLite TEXT timestamp
func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

// ---- AUTH STUFF ----

func (s *SQLiteDB) CheckIfEmailExists(ctx context.Context, email string) (bool, error) {
	var count int
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("CheckIfEmailExists: %w", err)
	}
	return count > 0, nil
}

func (s *SQLiteDB) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	now := time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = now
	}

	result, err := s.db.ExecContext(ctx, `
		INSERT INTO users (email, password_hash, auth_provider, provider_user_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		user.Email, user.PasswordHash, string(user.AuthProvider), user.ProviderUserID,
		formatTime(user.CreatedAt), formatTime(user.UpdatedAt))
	if err != nil {
		return 0, fmt.Errorf("CreateUser: %w", err)
	}
	return result.LastInsertId()
}

func (s *SQLiteDB) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	var authProvider, createdAt, updatedAt string
	err := s.db.QueryRowContext(ctx, `
		SELECT id, email, password_hash, auth_provider, provider_user_id, created_at, updated_at
		FROM users WHERE email = ?`, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &authProvider, &user.ProviderUserID,
		&createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("GetUserByEmail: %w", err)
	}

	user.AuthProvider = model.AuthProvider(authProvider)
	if user.CreatedAt, err = parseTime(createdAt); err != nil {
		return nil, fmt.Errorf("GetUserByEmail: parse created_at: %w", err)
	}
	if user.UpdatedAt, err = parseTime(updatedAt); err != nil {
		return nil, fmt.Errorf("GetUserByEmail: parse updated_at: %w", err)
	}

	return &user, nil
}

func (s *SQLiteDB) GetUserByID(ctx context.Context, userID int64) (*model.User, error) {
	var user model.User
	var authProvider, createdAt, updatedAt string
	err := s.db.QueryRowContext(ctx, `
		SELECT id, email, password_hash, auth_provider, provider_user_id, created_at, updated_at
		FROM users WHERE id = ?`, userID).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &authProvider, &user.ProviderUserID,
		&createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("GetUserByID: %w", err)
	}

	user.AuthProvider = model.AuthProvider(authProvider)
	if user.CreatedAt, err = parseTime(createdAt); err != nil {
		return nil, fmt.Errorf("GetUserByID: parse created_at: %w", err)
	}
	if user.UpdatedAt, err = parseTime(updatedAt); err != nil {
		return nil, fmt.Errorf("GetUserByID: parse updated_at: %w", err)
	}

	return &user, nil
}

// ---- SPECIES STUFF ----

func (s *SQLiteDB) SearchSpecies(ctx context.Context, query string, limit, offset int) ([]*model.Species, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, common_name, scientific_name, light_pref, default_water_interval_days,
			   default_fertilizer_interval_days, toxicity, care_notes, care_notes_source,
			   created_at, updated_at
		FROM species
		WHERE common_name LIKE ? OR scientific_name LIKE ?
		ORDER BY common_name
		LIMIT ? OFFSET ?`,
		"%"+query+"%", "%"+query+"%", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("SearchSpecies: %w", err)
	}
	defer rows.Close()

	var species []*model.Species
	for rows.Next() {
		var s model.Species
		var lightPref, toxicity, createdAt, updatedAt string
		err := rows.Scan(&s.ID, &s.CommonName, &s.ScientificName, &lightPref,
			&s.DefaultWaterIntervalDays, &s.DefaultFertilizerIntervalDays,
			&toxicity, &s.CareNotes, &s.CareNotesSource, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("SearchSpecies scan: %w", err)
		}

		s.LightPreference = model.LightPreference(lightPref)
		s.Toxicity = model.Toxicity(toxicity)
		if s.CreatedAt, err = parseTime(createdAt); err != nil {
			return nil, fmt.Errorf("SearchSpecies: parse created_at: %w", err)
		}
		if s.UpdatedAt, err = parseTime(updatedAt); err != nil {
			return nil, fmt.Errorf("SearchSpecies: parse updated_at: %w", err)
		}

		species = append(species, &s)
	}
	return species, rows.Err()
}

func (s *SQLiteDB) GetSpeciesByID(ctx context.Context, speciesID int64) (*model.Species, error) {
	var species model.Species
	var lightPref, toxicity, createdAt, updatedAt string
	err := s.db.QueryRowContext(ctx, `
		SELECT id, common_name, scientific_name, light_pref, default_water_interval_days,
			   default_fertilizer_interval_days, toxicity, care_notes, care_notes_source,
			   created_at, updated_at
		FROM species WHERE id = ?`, speciesID).Scan(
		&species.ID, &species.CommonName, &species.ScientificName, &lightPref,
		&species.DefaultWaterIntervalDays, &species.DefaultFertilizerIntervalDays,
		&toxicity, &species.CareNotes, &species.CareNotesSource, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("GetSpeciesByID: %w", err)
	}

	species.LightPreference = model.LightPreference(lightPref)
	species.Toxicity = model.Toxicity(toxicity)
	if species.CreatedAt, err = parseTime(createdAt); err != nil {
		return nil, fmt.Errorf("GetSpeciesByID: parse created_at: %w", err)
	}
	if species.UpdatedAt, err = parseTime(updatedAt); err != nil {
		return nil, fmt.Errorf("GetSpeciesByID: parse updated_at: %w", err)
	}

	return &species, nil
}

// ---- LOCATION STUFF ----

func (s *SQLiteDB) CreateLocation(ctx context.Context, location *model.Location) (int64, error) {
	result, err := s.db.ExecContext(ctx, `
		INSERT INTO location (user_id, name, created_at, updated_at)
		VALUES (?, ?, ?, ?)`,
		location.UserID, location.Name, formatTime(location.CreatedAt), formatTime(location.UpdatedAt))
	if err != nil {
		return 0, fmt.Errorf("CreateLocation: %w", err)
	}
	return result.LastInsertId()
}

func (s *SQLiteDB) ListLocations(ctx context.Context, userID int64) ([]*model.Location, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, user_id, name, created_at, updated_at
		FROM location WHERE user_id = ?
		ORDER BY name`, userID)
	if err != nil {
		return nil, fmt.Errorf("ListLocations: %w", err)
	}
	defer rows.Close()

	var locations []*model.Location
	for rows.Next() {
		var loc model.Location
		var createdAt, updatedAt string
		err := rows.Scan(&loc.ID, &loc.UserID, &loc.Name, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("ListLocations scan: %w", err)
		}

		if loc.CreatedAt, err = parseTime(createdAt); err != nil {
			return nil, fmt.Errorf("ListLocations: parse created_at: %w", err)
		}
		if loc.UpdatedAt, err = parseTime(updatedAt); err != nil {
			return nil, fmt.Errorf("ListLocations: parse updated_at: %w", err)
		}

		locations = append(locations, &loc)
	}
	return locations, rows.Err()
}

func (s *SQLiteDB) GetLocationByID(ctx context.Context, userID, locationID int64) (*model.Location, error) {
	var location model.Location
	var createdAt, updatedAt string
	err := s.db.QueryRowContext(ctx, `
		SELECT id, user_id, name, created_at, updated_at
		FROM location WHERE id = ? AND user_id = ?`, locationID, userID).Scan(
		&location.ID, &location.UserID, &location.Name, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("GetLocationByID: %w", err)
	}

	if location.CreatedAt, err = parseTime(createdAt); err != nil {
		return nil, fmt.Errorf("GetLocationByID: parse created_at: %w", err)
	}
	if location.UpdatedAt, err = parseTime(updatedAt); err != nil {
		return nil, fmt.Errorf("GetLocationByID: parse updated_at: %w", err)
	}

	return &location, nil
}

func (s *SQLiteDB) UpdateLocation(ctx context.Context, location *model.Location) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE location SET name = ?, updated_at = ?
		WHERE id = ? AND user_id = ?`,
		location.Name, formatTime(location.UpdatedAt), location.ID, location.UserID)
	if err != nil {
		return fmt.Errorf("UpdateLocation: %w", err)
	}
	return nil
}

func (s *SQLiteDB) DeleteLocation(ctx context.Context, userID, locationID int64) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM location WHERE id = ? AND user_id = ?", locationID, userID)
	if err != nil {
		return fmt.Errorf("DeleteLocation: %w", err)
	}
	return nil
}

// ---- USER_PLANTS STUFF ----

func (s *SQLiteDB) CreateUserPlant(ctx context.Context, plant *model.UserPlant) (int64, error) {
	result, err := s.db.ExecContext(ctx, `
		INSERT INTO user_plants (user_id, species_id, nickname, image_url, notes,
							     water_interval_days_override, fertilizer_interval_days_override,
							     location_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		plant.UserID, plant.SpeciesID, plant.Nickname, plant.ImageURL, plant.Notes,
		plant.WaterIntervalDaysOverride, plant.FertilizerIntervalDaysOverride,
		plant.LocationID, formatTime(plant.CreatedAt), formatTime(plant.UpdatedAt))
	if err != nil {
		return 0, fmt.Errorf("CreateUserPlant: %w", err)
	}
	return result.LastInsertId()
}

func (s *SQLiteDB) GetUserPlantByID(ctx context.Context, userID, plantID int64) (*model.UserPlant, error) {
	var plant model.UserPlant
	var createdAt, updatedAt string
	err := s.db.QueryRowContext(ctx, `
		SELECT id, user_id, species_id, nickname, image_url, notes,
			   water_interval_days_override, fertilizer_interval_days_override,
			   location_id, created_at, updated_at
		FROM user_plants WHERE id = ? AND user_id = ?`, plantID, userID).Scan(
		&plant.ID, &plant.UserID, &plant.SpeciesID, &plant.Nickname, &plant.ImageURL,
		&plant.Notes, &plant.WaterIntervalDaysOverride, &plant.FertilizerIntervalDaysOverride,
		&plant.LocationID, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("GetUserPlantByID: %w", err)
	}

	if plant.CreatedAt, err = parseTime(createdAt); err != nil {
		return nil, fmt.Errorf("GetUserPlantByID: parse created_at: %w", err)
	}
	if plant.UpdatedAt, err = parseTime(updatedAt); err != nil {
		return nil, fmt.Errorf("GetUserPlantByID: parse updated_at: %w", err)
	}

	return &plant, nil
}

func (s *SQLiteDB) ListUserPlants(ctx context.Context, userID int64, limit, offset int) ([]*model.UserPlant, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, user_id, species_id, nickname, image_url, notes,
			   water_interval_days_override, fertilizer_interval_days_override,
			   location_id, created_at, updated_at
		FROM user_plants WHERE user_id = ?
		ORDER BY id
		LIMIT ? OFFSET ?`, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("ListUserPlants: %w", err)
	}
	defer rows.Close()

	var plants []*model.UserPlant
	for rows.Next() {
		var plant model.UserPlant
		var createdAt, updatedAt string
		err := rows.Scan(&plant.ID, &plant.UserID, &plant.SpeciesID, &plant.Nickname,
			&plant.ImageURL, &plant.Notes, &plant.WaterIntervalDaysOverride,
			&plant.FertilizerIntervalDaysOverride, &plant.LocationID, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("ListUserPlants scan: %w", err)
		}

		if plant.CreatedAt, err = parseTime(createdAt); err != nil {
			return nil, fmt.Errorf("ListUserPlants: parse created_at: %w", err)
		}
		if plant.UpdatedAt, err = parseTime(updatedAt); err != nil {
			return nil, fmt.Errorf("ListUserPlants: parse updated_at: %w", err)
		}

		plants = append(plants, &plant)
	}
	return plants, rows.Err()
}

func (s *SQLiteDB) UpdateUserPlant(ctx context.Context, plant *model.UserPlant) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE user_plants SET nickname = ?, image_url = ?, notes = ?,
							   water_interval_days_override = ?, fertilizer_interval_days_override = ?,
							   location_id = ?, updated_at = ?
		WHERE id = ? AND user_id = ?`,
		plant.Nickname, plant.ImageURL, plant.Notes, plant.WaterIntervalDaysOverride,
		plant.FertilizerIntervalDaysOverride, plant.LocationID, formatTime(plant.UpdatedAt),
		plant.ID, plant.UserID)
	if err != nil {
		return fmt.Errorf("UpdateUserPlant: %w", err)
	}
	return nil
}

func (s *SQLiteDB) DeleteUserPlant(ctx context.Context, userID, plantID int64) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM user_plants WHERE id = ? AND user_id = ?", plantID, userID)
	if err != nil {
		return fmt.Errorf("DeleteUserPlant: %w", err)
	}
	return nil
}

// ---- PLANT_TASKS STUFF ----

func (s *SQLiteDB) CreatePlantTask(ctx context.Context, task *model.PlantTask) (int64, error) {
	snoozedUntil := (*string)(nil)
	if task.SnoozedUntil != nil {
		str := formatTime(*task.SnoozedUntil)
		snoozedUntil = &str
	}
	nextDueAt := (*string)(nil)
	if task.NextDueAt != nil {
		str := formatTime(*task.NextDueAt)
		nextDueAt = &str
	}

	result, err := s.db.ExecContext(ctx, `
		INSERT INTO plant_tasks (plant_id, task_type, snoozed_until, interval_days,
								 tolerance_days, next_due_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		task.PlantID, string(task.TaskType), snoozedUntil, task.IntervalDays,
		task.ToleranceDays, nextDueAt, formatTime(task.CreatedAt), formatTime(task.UpdatedAt))
	if err != nil {
		return 0, fmt.Errorf("CreatePlantTask: %w", err)
	}
	return result.LastInsertId()
}

func (s *SQLiteDB) GetPlantTasksByPlantID(ctx context.Context, userID, plantID int64) ([]*model.PlantTask, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT pt.id, pt.plant_id, pt.task_type, pt.snoozed_until, pt.interval_days,
			   pt.tolerance_days, pt.next_due_at, pt.created_at, pt.updated_at
		FROM plant_tasks pt
		JOIN user_plants up ON pt.plant_id = up.id
		WHERE pt.plant_id = ? AND up.user_id = ?
		ORDER BY pt.id`, plantID, userID)
	if err != nil {
		return nil, fmt.Errorf("GetPlantTasksByPlantID: %w", err)
	}
	defer rows.Close()

	var tasks []*model.PlantTask
	for rows.Next() {
		var task model.PlantTask
		var taskType, snoozedUntil, nextDueAt, createdAt, updatedAt sql.NullString
		err := rows.Scan(&task.ID, &task.PlantID, &taskType, &snoozedUntil,
			&task.IntervalDays, &task.ToleranceDays, &nextDueAt, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("GetPlantTasksByPlantID scan: %w", err)
		}

		task.TaskType = model.TaskType(taskType.String)
		if snoozedUntil.Valid {
			if t, err := parseTime(snoozedUntil.String); err == nil {
				task.SnoozedUntil = &t
			}
		}
		if nextDueAt.Valid {
			if t, err := parseTime(nextDueAt.String); err == nil {
				task.NextDueAt = &t
			}
		}
		if task.CreatedAt, err = parseTime(createdAt.String); err != nil {
			return nil, fmt.Errorf("GetPlantTasksByPlantID: parse created_at: %w", err)
		}
		if task.UpdatedAt, err = parseTime(updatedAt.String); err != nil {
			return nil, fmt.Errorf("GetPlantTasksByPlantID: parse updated_at: %w", err)
		}

		tasks = append(tasks, &task)
	}
	return tasks, rows.Err()
}

func (s *SQLiteDB) GetPlantTaskByID(ctx context.Context, userID, taskID int64) (*model.PlantTask, error) {
	var task model.PlantTask
	var taskType, snoozedUntil, nextDueAt, createdAt, updatedAt sql.NullString
	err := s.db.QueryRowContext(ctx, `
		SELECT pt.id, pt.plant_id, pt.task_type, pt.snoozed_until, pt.interval_days,
			   pt.tolerance_days, pt.next_due_at, pt.created_at, pt.updated_at
		FROM plant_tasks pt
		JOIN user_plants up ON pt.plant_id = up.id
		WHERE pt.id = ? AND up.user_id = ?`, taskID, userID).Scan(
		&task.ID, &task.PlantID, &taskType, &snoozedUntil, &task.IntervalDays,
		&task.ToleranceDays, &nextDueAt, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("GetPlantTaskByID: %w", err)
	}

	task.TaskType = model.TaskType(taskType.String)
	if snoozedUntil.Valid {
		if t, err := parseTime(snoozedUntil.String); err == nil {
			task.SnoozedUntil = &t
		}
	}
	if nextDueAt.Valid {
		if t, err := parseTime(nextDueAt.String); err == nil {
			task.NextDueAt = &t
		}
	}
	if task.CreatedAt, err = parseTime(createdAt.String); err != nil {
		return nil, fmt.Errorf("GetPlantTaskByID: parse created_at: %w", err)
	}
	if task.UpdatedAt, err = parseTime(updatedAt.String); err != nil {
		return nil, fmt.Errorf("GetPlantTaskByID: parse updated_at: %w", err)
	}

	return &task, nil
}

func (s *SQLiteDB) UpdatePlantTask(ctx context.Context, task *model.PlantTask) error {
	snoozedUntil := (*string)(nil)
	if task.SnoozedUntil != nil {
		str := formatTime(*task.SnoozedUntil)
		snoozedUntil = &str
	}
	nextDueAt := (*string)(nil)
	if task.NextDueAt != nil {
		str := formatTime(*task.NextDueAt)
		nextDueAt = &str
	}

	_, err := s.db.ExecContext(ctx, `
		UPDATE plant_tasks SET task_type = ?, snoozed_until = ?, interval_days = ?,
							   tolerance_days = ?, next_due_at = ?, updated_at = ?
		WHERE id = ?`,
		string(task.TaskType), snoozedUntil, task.IntervalDays, task.ToleranceDays,
		nextDueAt, formatTime(task.UpdatedAt), task.ID)
	if err != nil {
		return fmt.Errorf("UpdatePlantTask: %w", err)
	}
	return nil
}

func (s *SQLiteDB) DeletePlantTask(ctx context.Context, userID, taskID int64) error {
	_, err := s.db.ExecContext(ctx, `
		DELETE FROM plant_tasks 
		WHERE id = ? AND plant_id IN (SELECT id FROM user_plants WHERE user_id = ?)`,
		taskID, userID)
	if err != nil {
		return fmt.Errorf("DeletePlantTask: %w", err)
	}
	return nil
}

func (s *SQLiteDB) GetOverdueTasks(ctx context.Context, userID int64) ([]*model.PlantTask, error) {
	now := formatTime(time.Now())
	rows, err := s.db.QueryContext(ctx, `
		SELECT pt.id, pt.plant_id, pt.task_type, pt.snoozed_until, pt.interval_days,
			   pt.tolerance_days, pt.next_due_at, pt.created_at, pt.updated_at
		FROM plant_tasks pt
		JOIN user_plants up ON pt.plant_id = up.id
		WHERE up.user_id = ? 
		  AND pt.next_due_at IS NOT NULL 
		  AND pt.next_due_at <= ?
		  AND (pt.snoozed_until IS NULL OR pt.snoozed_until <= ?)
		ORDER BY pt.next_due_at`, userID, now, now)
	if err != nil {
		return nil, fmt.Errorf("GetOverdueTasks: %w", err)
	}
	defer rows.Close()

	var tasks []*model.PlantTask
	for rows.Next() {
		var task model.PlantTask
		var taskType, snoozedUntil, nextDueAt, createdAt, updatedAt sql.NullString
		err := rows.Scan(&task.ID, &task.PlantID, &taskType, &snoozedUntil,
			&task.IntervalDays, &task.ToleranceDays, &nextDueAt, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("GetOverdueTasks scan: %w", err)
		}

		task.TaskType = model.TaskType(taskType.String)
		if snoozedUntil.Valid {
			if t, err := parseTime(snoozedUntil.String); err == nil {
				task.SnoozedUntil = &t
			}
		}
		if nextDueAt.Valid {
			if t, err := parseTime(nextDueAt.String); err == nil {
				task.NextDueAt = &t
			}
		}
		if task.CreatedAt, err = parseTime(createdAt.String); err != nil {
			return nil, fmt.Errorf("GetOverdueTasks: parse created_at: %w", err)
		}
		if task.UpdatedAt, err = parseTime(updatedAt.String); err != nil {
			return nil, fmt.Errorf("GetOverdueTasks: parse updated_at: %w", err)
		}

		tasks = append(tasks, &task)
	}
	return tasks, rows.Err()
}

// ---- CARE_EVENTS STUFF ----

func (s *SQLiteDB) CreateCareEvent(ctx context.Context, event *model.CareEvent) (int64, error) {
	result, err := s.db.ExecContext(ctx, `
		INSERT INTO care_events (plant_id, task_id, event_type, happened_at, notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		event.PlantID, event.TaskID, string(event.EventType), formatTime(event.HappenedAt),
		event.Notes, formatTime(event.CreatedAt), formatTime(event.UpdatedAt))
	if err != nil {
		return 0, fmt.Errorf("CreateCareEvent: %w", err)
	}
	return result.LastInsertId()
}

func (s *SQLiteDB) GetCareEventsByPlantID(ctx context.Context, userID, plantID int64, limit, offset int) ([]*model.CareEvent, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT ce.id, ce.plant_id, ce.task_id, ce.event_type, ce.happened_at, ce.notes,
			   ce.created_at, ce.updated_at
		FROM care_events ce
		JOIN user_plants up ON ce.plant_id = up.id
		WHERE ce.plant_id = ? AND up.user_id = ?
		ORDER BY ce.happened_at DESC
		LIMIT ? OFFSET ?`, plantID, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("GetCareEventsByPlantID: %w", err)
	}
	defer rows.Close()

	var events []*model.CareEvent
	for rows.Next() {
		var event model.CareEvent
		var eventType, happenedAt, createdAt, updatedAt string
		err := rows.Scan(&event.ID, &event.PlantID, &event.TaskID, &eventType,
			&happenedAt, &event.Notes, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("GetCareEventsByPlantID scan: %w", err)
		}

		event.EventType = model.EventType(eventType)
		if event.HappenedAt, err = parseTime(happenedAt); err != nil {
			return nil, fmt.Errorf("GetCareEventsByPlantID: parse happened_at: %w", err)
		}
		if event.CreatedAt, err = parseTime(createdAt); err != nil {
			return nil, fmt.Errorf("GetCareEventsByPlantID: parse created_at: %w", err)
		}
		if event.UpdatedAt, err = parseTime(updatedAt); err != nil {
			return nil, fmt.Errorf("GetCareEventsByPlantID: parse updated_at: %w", err)
		}

		events = append(events, &event)
	}
	return events, rows.Err()
}

func (s *SQLiteDB) GetCareEventByID(ctx context.Context, userID, eventID int64) (*model.CareEvent, error) {
	var event model.CareEvent
	var eventType, happenedAt, createdAt, updatedAt string
	err := s.db.QueryRowContext(ctx, `
		SELECT ce.id, ce.plant_id, ce.task_id, ce.event_type, ce.happened_at, ce.notes,
			   ce.created_at, ce.updated_at
		FROM care_events ce
		JOIN user_plants up ON ce.plant_id = up.id
		WHERE ce.id = ? AND up.user_id = ?`, eventID, userID).Scan(
		&event.ID, &event.PlantID, &event.TaskID, &eventType, &happenedAt,
		&event.Notes, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("GetCareEventByID: %w", err)
	}

	event.EventType = model.EventType(eventType)
	if event.HappenedAt, err = parseTime(happenedAt); err != nil {
		return nil, fmt.Errorf("GetCareEventByID: parse happened_at: %w", err)
	}
	if event.CreatedAt, err = parseTime(createdAt); err != nil {
		return nil, fmt.Errorf("GetCareEventByID: parse created_at: %w", err)
	}
	if event.UpdatedAt, err = parseTime(updatedAt); err != nil {
		return nil, fmt.Errorf("GetCareEventByID: parse updated_at: %w", err)
	}

	return &event, nil
}

func (s *SQLiteDB) UpdateCareEvent(ctx context.Context, event *model.CareEvent) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE care_events SET task_id = ?, event_type = ?, happened_at = ?, notes = ?, updated_at = ?
		WHERE id = ?`,
		event.TaskID, string(event.EventType), formatTime(event.HappenedAt),
		event.Notes, formatTime(event.UpdatedAt), event.ID)
	if err != nil {
		return fmt.Errorf("UpdateCareEvent: %w", err)
	}
	return nil
}

func (s *SQLiteDB) DeleteCareEvent(ctx context.Context, userID, eventID int64) error {
	_, err := s.db.ExecContext(ctx, `
		DELETE FROM care_events 
		WHERE id = ? AND plant_id IN (SELECT id FROM user_plants WHERE user_id = ?)`,
		eventID, userID)
	if err != nil {
		return fmt.Errorf("DeleteCareEvent: %w", err)
	}
	return nil
}
