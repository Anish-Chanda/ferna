package sqlite3

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/anish-chanda/ferna/db"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupInMemoryDB(t *testing.T) db.Database {
	dsn := "file::memory:?mode=memory&cache=shared&_foreign_keys=ON"
	sq := NewSQLiteDB()
	require.NoError(t, sq.Connect(dsn))
	require.NoError(t, sq.Migrate())
	return sq
}

func TestCheckIfEmailExists_Empty(t *testing.T) {
	sq := setupInMemoryDB(t)
	defer sq.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	exists, err := sq.CheckIfEmailExists(ctx, "noone@nowhere.com")
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestCreateUser_And_CheckIfEmailExists(t *testing.T) {
	sq := setupInMemoryDB(t)
	defer sq.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	id, err := sq.CreateUser(ctx, "alice@example.com", "dummy-hash")
	require.NoError(t, err)
	assert.True(t, id > 0)

	exists, err := sq.CheckIfEmailExists(ctx, "alice@example.com")
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestGetUserByEmail_Success(t *testing.T) {
	sq := setupInMemoryDB(t)
	defer sq.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	hash := "s0m3-h4sh"
	id, err := sq.CreateUser(ctx, "bob@example.com", hash)
	require.NoError(t, err)

	userRec, err := sq.GetUserByEmail(ctx, "bob@example.com")
	require.NoError(t, err)
	require.NotNil(t, userRec)

	assert.Equal(t, id, userRec.ID)
	assert.Equal(t, "bob@example.com", userRec.Email)
	assert.Equal(t, hash, userRec.PassHash)
	assert.WithinDuration(t, time.Now().UTC(), userRec.CreatedAt.UTC(), time.Minute)
	assert.WithinDuration(t, time.Now().UTC(), userRec.UpdatedAt.UTC(), time.Minute)
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	sq := setupInMemoryDB(t)
	defer sq.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	userRec, err := sq.GetUserByEmail(ctx, "charlie@nowhere.com")
	require.NoError(t, err)
	assert.Nil(t, userRec)
}

// TestCreateMultipleUsers ensures that inserting two different emails yields two different IDs,
// and that CheckIfEmailExists returns true for both.
func TestCreateMultipleUsers(t *testing.T) {
	sq := setupInMemoryDB(t)
	defer sq.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Insert the first user
	email1 := "user1@example.com"
	hash1 := "hash1"
	id1, err := sq.CreateUser(ctx, email1, hash1)
	require.NoError(t, err)
	assert.True(t, id1 > 0)

	// Insert a second, distinct user
	email2 := "user2@example.com"
	hash2 := "hash2"
	id2, err := sq.CreateUser(ctx, email2, hash2)
	require.NoError(t, err)
	assert.True(t, id2 > 0)
	assert.NotEqual(t, id1, id2, "IDs for distinct users should differ")

	// Check existence of both emails
	exists1, err := sq.CheckIfEmailExists(ctx, email1)
	require.NoError(t, err)
	assert.True(t, exists1, "email1 should exist")

	exists2, err := sq.CheckIfEmailExists(ctx, email2)
	require.NoError(t, err)
	assert.True(t, exists2, "email2 should exist")
}

// TestCreateDuplicateUser verifies that attempting to insert the same email twice returns an error.
func TestCreateDuplicateUser(t *testing.T) {
	sq := setupInMemoryDB(t)
	defer sq.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	email := "duplicate@example.com"
	hash := "some-hash"

	// First insertion should succeed
	id, err := sq.CreateUser(ctx, email, hash)
	require.NoError(t, err)
	assert.True(t, id > 0)

	// Second insertion with the same email must fail
	_, err = sq.CreateUser(ctx, email, hash)
	require.Error(t, err, "Inserting a duplicate email should error out")

	// Verify the error mentions "UNIQUE constraint failed" (SQLite)
	assert.Contains(t, err.Error(), "UNIQUE constraint failed", "Error should indicate UNIQUE constraint failure")
}

// TestContextCancellation ensures that if the context is already canceled or times out,
// each method returns a context-related error.
func TestContextCancellation(t *testing.T) {
	sq := setupInMemoryDB(t)
	defer sq.Close()

	// Create a context and cancel it immediately
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	cancel()
	time.Sleep(5 * time.Millisecond) // allow cancellation to propagate

	// 1) CheckIfEmailExists with canceled context
	_, err := sq.CheckIfEmailExists(ctx, "whoever@example.com")
	require.Error(t, err, "CheckIfEmailExists should error on canceled context")
	assert.True(t, errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded),
		"Error should be context.Canceled or DeadlineExceeded")

	// 2) CreateUser with canceled context
	_, err = sq.CreateUser(ctx, "someone@example.com", "some-hash")
	require.Error(t, err, "CreateUser should error on canceled context")
	assert.True(t, errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded),
		"Error should be context.Canceled or DeadlineExceeded")

	// 3) GetUserByEmail with canceled context
	_, err = sq.GetUserByEmail(ctx, "someone@example.com")
	require.Error(t, err, "GetUserByEmail should error on canceled context")
	assert.True(t, errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded),
		"Error should be context.Canceled or DeadlineExceeded")
}

// TestCloseThenUse verifies that once Close() is called, subsequent calls to any method return an error.
func TestCloseThenUse(t *testing.T) {
	sq := setupInMemoryDB(t)

	// Close immediately
	require.NoError(t, sq.Close(), "Close should not error")

	// Now attempt CheckIfEmailExists
	ctx := context.Background()
	_, err := sq.CheckIfEmailExists(ctx, "postclose@example.com")
	require.Error(t, err, "Call after Close() should error")
	assert.Contains(t, err.Error(), "database is closed", "Error should mention that the DB is closed")

	// Attempt CreateUser
	_, err = sq.CreateUser(ctx, "postclose@example.com", "hash")
	require.Error(t, err, "CreateUser after Close() should error")
	assert.Contains(t, err.Error(), "database is closed", "Error should mention that the DB is closed")

	// Attempt GetUserByEmail
	_, err = sq.GetUserByEmail(ctx, "postclose@example.com")
	require.Error(t, err, "GetUserByEmail after Close() should error")
	assert.Contains(t, err.Error(), "database is closed", "Error should mention that the DB is closed")
}

// TestExecContext verifies that ExecContext works correctly for INSERT, UPDATE, DELETE operations
func TestExecContext_Insert(t *testing.T) {
	sq := setupInMemoryDB(t)
	defer sq.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Test INSERT using ExecContext
	result, err := sq.ExecContext(ctx,
		`INSERT INTO species(common_name, scientific_name, default_watering_frequency_days) VALUES (?, ?, ?)`,
		"Test Plant", "Testus plantus", 7)
	require.NoError(t, err)

	// Check that the insert was successful
	rowsAffected, err := result.RowsAffected()
	require.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)

	// Verify the record was inserted by searching for it
	species, err := sq.SearchSpecies(ctx, "Test Plant", 10, 0)
	require.NoError(t, err)
	require.Len(t, species, 1)
	assert.Equal(t, "Test Plant", species[0].CommonName)
	assert.Equal(t, "Testus plantus", species[0].ScientificName)
	assert.Equal(t, 7, species[0].DefaultWateringFrequency)
}

func TestExecContext_InsertOrIgnore(t *testing.T) {
	sq := setupInMemoryDB(t)
	defer sq.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// First insert
	result1, err := sq.ExecContext(ctx,
		`INSERT OR IGNORE INTO species(common_name, scientific_name, default_watering_frequency_days) VALUES (?, ?, ?)`,
		"Rose", "Rosa rubiginosa", 3)
	require.NoError(t, err)

	rowsAffected, err := result1.RowsAffected()
	require.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)

	// Second insert with same scientific_name should be ignored due to UNIQUE constraint
	result2, err := sq.ExecContext(ctx,
		`INSERT OR IGNORE INTO species(common_name, scientific_name, default_watering_frequency_days) VALUES (?, ?, ?)`,
		"Different Rose", "Rosa rubiginosa", 5)
	require.NoError(t, err)

	rowsAffected, err = result2.RowsAffected()
	require.NoError(t, err)
	assert.Equal(t, int64(0), rowsAffected) // Should be 0 because it was ignored

	// Verify only one record exists
	species, err := sq.SearchSpecies(ctx, "Rosa rubiginosa", 10, 0)
	require.NoError(t, err)
	require.Len(t, species, 1)
	assert.Equal(t, "Rose", species[0].CommonName) // Should still be the original
}

func TestExecContext_Update(t *testing.T) {
	sq := setupInMemoryDB(t)
	defer sq.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// First insert a record
	_, err := sq.ExecContext(ctx,
		`INSERT INTO species(common_name, scientific_name, default_watering_frequency_days) VALUES (?, ?, ?)`,
		"Cactus", "Cactaceae genericus", 14)
	require.NoError(t, err)

	// Update the watering frequency
	result, err := sq.ExecContext(ctx,
		`UPDATE species SET default_watering_frequency_days = ? WHERE scientific_name = ?`,
		21, "Cactaceae genericus")
	require.NoError(t, err)

	rowsAffected, err := result.RowsAffected()
	require.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)

	// Verify the update
	species, err := sq.SearchSpecies(ctx, "Cactus", 10, 0)
	require.NoError(t, err)
	require.Len(t, species, 1)
	assert.Equal(t, 21, species[0].DefaultWateringFrequency)
}

func TestExecContext_Delete(t *testing.T) {
	sq := setupInMemoryDB(t)
	defer sq.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// First insert a record
	_, err := sq.ExecContext(ctx,
		`INSERT INTO species(common_name, scientific_name, default_watering_frequency_days) VALUES (?, ?, ?)`,
		"Temporary Plant", "Temporarius plantus", 1)
	require.NoError(t, err)

	// Verify it exists
	species, err := sq.SearchSpecies(ctx, "Temporary Plant", 10, 0)
	require.NoError(t, err)
	require.Len(t, species, 1)

	// Delete the record
	result, err := sq.ExecContext(ctx,
		`DELETE FROM species WHERE scientific_name = ?`,
		"Temporarius plantus")
	require.NoError(t, err)

	rowsAffected, err := result.RowsAffected()
	require.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)

	// Verify it's gone
	species, err = sq.SearchSpecies(ctx, "Temporary Plant", 10, 0)
	require.NoError(t, err)
	assert.Len(t, species, 0)
}

func TestExecContext_CanceledContext(t *testing.T) {
	sq := setupInMemoryDB(t)
	defer sq.Close()

	// Create a context and cancel it immediately
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	cancel()
	time.Sleep(5 * time.Millisecond) // allow cancellation to propagate

	// ExecContext with canceled context should error
	_, err := sq.ExecContext(ctx,
		`INSERT INTO species(common_name, scientific_name, default_watering_frequency_days) VALUES (?, ?, ?)`,
		"Test", "Test", 1)
	require.Error(t, err, "ExecContext should error on canceled context")
	assert.True(t, errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded),
		"Error should be context.Canceled or DeadlineExceeded")
}
