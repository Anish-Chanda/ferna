package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anish-chanda/ferna/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeDB implements db.Database for testing HandleLogin and HandleSignup.
type fakeDB struct {
	GetByEmailFunc  func(ctx context.Context, email string) (*model.User, error)
	CheckExistsFunc func(ctx context.Context, email string) (bool, error)
	CreateUserFunc  func(ctx context.Context, email, passHash string) (int64, error)
}

func (f *fakeDB) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return f.GetByEmailFunc(ctx, email)
}

func (f *fakeDB) CheckIfEmailExists(ctx context.Context, email string) (bool, error) {
	return f.CheckExistsFunc(ctx, email)
}

func (f *fakeDB) CreateUser(ctx context.Context, email, passHash string) (int64, error) {
	return f.CreateUserFunc(ctx, email, passHash)
}

// Unused methods to satisfy interface:
func (f *fakeDB) Connect(dsn string) error                               { return nil }
func (f *fakeDB) Close() error                                           { return nil }
func (f *fakeDB) Migrate() error                                         { return nil }
func (f *fakeDB) SomeUnusedMethod1()                                     { /* no-op */ }
func (f *fakeDB) SomeUnusedMethod2(arg interface{}) (interface{}, error) { return nil, nil }

//-----------------------
// Tests for HandleLogin
//-----------------------

func TestHandleLogin_UserNotFound(t *testing.T) {
	fdb := &fakeDB{
		GetByEmailFunc: func(ctx context.Context, email string) (*model.User, error) {
			return nil, nil // no such user
		},
	}
	ok, err := HandleLogin(fdb, "nouser@example.com", "any")
	assert.False(t, ok)
	assert.NoError(t, err)
}

func TestHandleLogin_DBError(t *testing.T) {
	sentinel := errors.New("DB failure")
	fdb := &fakeDB{
		GetByEmailFunc: func(ctx context.Context, email string) (*model.User, error) {
			return nil, sentinel
		},
	}
	ok, err := HandleLogin(fdb, "user@example.com", "any")
	assert.False(t, ok)
	require.Error(t, err)
	assert.Equal(t, sentinel, err)
}

func TestHandleLogin_VerifyPasswordError(t *testing.T) {
	// Return a user with an invalid PassHash so VerifyPassword fails format check.
	fdb := &fakeDB{
		GetByEmailFunc: func(ctx context.Context, email string) (*model.User, error) {
			return &model.User{
				ID:       1,
				Email:    email,
				PassHash: "invalid-format",
			}, nil
		},
	}
	ok, err := HandleLogin(fdb, "user@example.com", "password")
	assert.False(t, ok)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to verify password")
}

func TestHandleLogin_InvalidCredentials(t *testing.T) {
	// First generate a valid hash for "correctpass"
	validHash, err := HashPassword("correctpass")
	require.NoError(t, err)

	// Return user with that PassHash
	fdb := &fakeDB{
		GetByEmailFunc: func(ctx context.Context, email string) (*model.User, error) {
			return &model.User{
				ID:       2,
				Email:    email,
				PassHash: validHash,
			}, nil
		},
	}
	// Supply wrong password
	ok, err := HandleLogin(fdb, "user@example.com", "wrongpass")
	assert.False(t, ok)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestHandleLogin_Success(t *testing.T) {
	// Generate a valid hash for "mypassword"
	validHash, err := HashPassword("mypassword")
	require.NoError(t, err)

	fdb := &fakeDB{
		GetByEmailFunc: func(ctx context.Context, email string) (*model.User, error) {
			return &model.User{
				ID:       3,
				Email:    email,
				PassHash: validHash,
			}, nil
		},
	}
	ok, err := HandleLogin(fdb, "USER@Example.com", "mypassword")
	assert.True(t, ok)
	assert.NoError(t, err)
}

//------------------------
// Tests for HandleSignup
//------------------------

func TestHandleSignup_InvalidJSON(t *testing.T) {
	fdb := &fakeDB{}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBufferString("not-json"))

	HandleSignup(fdb, rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid request body")
}

func TestHandleSignup_MissingFields(t *testing.T) {
	fdb := &fakeDB{}
	rec := httptest.NewRecorder()

	// Missing password
	body := `{"user":"", "passwd":""}`
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBufferString(body))
	HandleSignup(fdb, rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Email and password are required")
}

func TestHandleSignup_EmailExists(t *testing.T) {
	fdb := &fakeDB{
		CheckExistsFunc: func(ctx context.Context, email string) (bool, error) {
			return true, nil // email already exists
		},
	}
	rec := httptest.NewRecorder()
	payload := map[string]string{"user": "exists@example.com", "passwd": "any"}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))

	HandleSignup(fdb, rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)
	assert.Contains(t, rec.Body.String(), "Email already exists")
}

func TestHandleSignup_CheckExistsError(t *testing.T) {
	sentinel := errors.New("checkExists failure")
	fdb := &fakeDB{
		CheckExistsFunc: func(ctx context.Context, email string) (bool, error) {
			return false, sentinel
		},
	}
	rec := httptest.NewRecorder()
	payload := map[string]string{"user": "fail@example.com", "passwd": "any"}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))

	HandleSignup(fdb, rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Internal server error")
}

func TestHandleSignup_CreateUserError(t *testing.T) {
	fdb := &fakeDB{
		CheckExistsFunc: func(ctx context.Context, email string) (bool, error) {
			return false, nil
		},
		CreateUserFunc: func(ctx context.Context, email, passHash string) (int64, error) {
			return 0, errors.New("createUser failure")
		},
	}
	rec := httptest.NewRecorder()
	payload := map[string]string{"user": "new@example.com", "passwd": "any"}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))

	HandleSignup(fdb, rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Internal server error")
}

func TestHandleSignup_Success(t *testing.T) {
	fdb := &fakeDB{
		CheckExistsFunc: func(ctx context.Context, email string) (bool, error) {
			return false, nil
		},
		CreateUserFunc: func(ctx context.Context, email, passHash string) (int64, error) {
			return 42, nil
		},
	}
	rec := httptest.NewRecorder()
	payload := map[string]string{"user": "Alice@Example.COM", "passwd": "strongpass"}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))

	HandleSignup(fdb, rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp struct {
		Success bool  `json:"success"`
		UserID  int64 `json:"user_id"`
	}
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, int64(42), resp.UserID)
}
