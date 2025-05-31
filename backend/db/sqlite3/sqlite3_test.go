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
