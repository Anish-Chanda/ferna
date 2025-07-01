package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anish-chanda/ferna/model"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}

type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) CreatePlant(ctx context.Context, p *model.Plant) (int64, error) {
	args := m.Called(ctx, p)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDatabase) GetPlantByID(ctx context.Context, userID, plantID int64) (*model.Plant, error) {
	args := m.Called(ctx, userID, plantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Plant), args.Error(1)
}

func (m *MockDatabase) ListPlants(ctx context.Context, userID int64, limit, offset int) ([]*model.Plant, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]*model.Plant), args.Error(1)
}

func (m *MockDatabase) UpdatePlant(ctx context.Context, p *model.Plant) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockDatabase) DeletePlant(ctx context.Context, userID, plantID int64) error {
	args := m.Called(ctx, userID, plantID)
	return args.Error(0)
}

// Implement other required interface methods as stubs
func (m *MockDatabase) CheckIfEmailExists(ctx context.Context, email string) (bool, error) {
	return false, nil
}
func (m *MockDatabase) CreateUser(ctx context.Context, email, passHash string) (int64, error) {
	return 0, nil
}
func (m *MockDatabase) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return nil, nil
}
func (m *MockDatabase) SearchSpecies(ctx context.Context, query string, limit, offset int) ([]*model.Species, error) {
	return nil, nil
}
func (m *MockDatabase) GetSpeciesByID(ctx context.Context, speciesID int64) (*model.Species, error) {
	args := m.Called(ctx, speciesID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Species), args.Error(1)
}
func (m *MockDatabase) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}
func (m *MockDatabase) Connect(dsn string) error { return nil }
func (m *MockDatabase) Close() error             { return nil }
func (m *MockDatabase) Migrate() error           { return nil }

func createTestRequest(method, url string, body interface{}) *http.Request {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, url, &buf)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestCreatePlant(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		mockDB := new(MockDatabase)

		// Mock the species lookup
		species := &model.Species{
			ID:                       123,
			CommonName:               "Test Plant",
			DefaultWateringFrequency: 7,
		}
		mockDB.On("GetSpeciesByID", mock.Anything, int64(123)).Return(species, nil)
		mockDB.On("CreatePlant", mock.Anything, mock.AnythingOfType("*model.Plant")).Return(int64(123), nil)

		body := map[string]interface{}{
			"species_id":              123,
			"nickname":                "My Plant",
			"watering_frequency_days": 7,
		}
		req := createTestRequest("POST", "/plants", body)
		req = req.WithContext(context.WithValue(req.Context(), "user", map[string]interface{}{"uid": "1"}))

		w := httptest.NewRecorder()
		handler := CreatePlant(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response model.Plant
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, int64(123), response.ID)
		assert.Equal(t, int64(1), response.UserID)

		mockDB.AssertExpectations(t)
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockDB := new(MockDatabase)

		req := httptest.NewRequest("POST", "/plants", bytes.NewReader([]byte("invalid json")))
		w := httptest.NewRecorder()

		handler := CreatePlant(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request body")
	})

	t.Run("unauthorized - no user context", func(t *testing.T) {
		mockDB := new(MockDatabase)

		body := map[string]interface{}{"species_id": 123}
		req := createTestRequest("POST", "/plants", body)

		w := httptest.NewRecorder()
		handler := CreatePlant(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("database error", func(t *testing.T) {
		mockDB := new(MockDatabase)

		// Mock the species lookup to succeed
		species := &model.Species{
			ID:                       123,
			CommonName:               "Test Plant",
			DefaultWateringFrequency: 7,
		}
		mockDB.On("GetSpeciesByID", mock.Anything, int64(123)).Return(species, nil)
		mockDB.On("CreatePlant", mock.Anything, mock.AnythingOfType("*model.Plant")).Return(int64(0), errors.New("db error"))

		body := map[string]interface{}{"species_id": 123}
		req := createTestRequest("POST", "/plants", body)
		req = req.WithContext(context.WithValue(req.Context(), "user", map[string]interface{}{"uid": "1"}))

		w := httptest.NewRecorder()
		handler := CreatePlant(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to create plant")
		mockDB.AssertExpectations(t)
	})

	t.Run("empty request body", func(t *testing.T) {
		mockDB := new(MockDatabase)

		req := createTestRequest("POST", "/plants", map[string]interface{}{})
		req = req.WithContext(context.WithValue(req.Context(), "user", map[string]interface{}{"uid": "1"}))

		w := httptest.NewRecorder()
		handler := CreatePlant(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "species_id is required")
	})

	t.Run("missing species_id", func(t *testing.T) {
		mockDB := new(MockDatabase)

		body := map[string]interface{}{
			"nickname": "My Plant",
		}
		req := createTestRequest("POST", "/plants", body)
		req = req.WithContext(context.WithValue(req.Context(), "user", map[string]interface{}{"uid": "1"}))

		w := httptest.NewRecorder()
		handler := CreatePlant(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "species_id is required")
	})

	t.Run("species not found", func(t *testing.T) {
		mockDB := new(MockDatabase)
		mockDB.On("GetSpeciesByID", mock.Anything, int64(999)).Return(nil, nil)

		body := map[string]interface{}{
			"species_id": 999,
		}
		req := createTestRequest("POST", "/plants", body)
		req = req.WithContext(context.WithValue(req.Context(), "user", map[string]interface{}{"uid": "1"}))

		w := httptest.NewRecorder()
		handler := CreatePlant(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Species not found")
		mockDB.AssertExpectations(t)
	})

	t.Run("species lookup error", func(t *testing.T) {
		mockDB := new(MockDatabase)
		mockDB.On("GetSpeciesByID", mock.Anything, int64(123)).Return(nil, errors.New("database error"))

		body := map[string]interface{}{
			"species_id": 123,
		}
		req := createTestRequest("POST", "/plants", body)
		req = req.WithContext(context.WithValue(req.Context(), "user", map[string]interface{}{"uid": "1"}))

		w := httptest.NewRecorder()
		handler := CreatePlant(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to get species")
		mockDB.AssertExpectations(t)
	})

	t.Run("defaults applied when not provided", func(t *testing.T) {
		mockDB := new(MockDatabase)

		// Mock the species lookup
		species := &model.Species{
			ID:                       123,
			CommonName:               "Rose",
			DefaultWateringFrequency: 3,
		}
		mockDB.On("GetSpeciesByID", mock.Anything, int64(123)).Return(species, nil)

		// Capture the plant that gets created to verify defaults
		var capturedPlant *model.Plant
		mockDB.On("CreatePlant", mock.Anything, mock.AnythingOfType("*model.Plant")).Run(func(args mock.Arguments) {
			capturedPlant = args.Get(1).(*model.Plant)
		}).Return(int64(456), nil)

		body := map[string]interface{}{
			"species_id": 123,
			// No nickname or watering_frequency_days provided
		}
		req := createTestRequest("POST", "/plants", body)
		req = req.WithContext(context.WithValue(req.Context(), "user", map[string]interface{}{"uid": "1"}))

		w := httptest.NewRecorder()
		handler := CreatePlant(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		// Verify defaults were applied
		assert.Equal(t, "Rose", *capturedPlant.Nickname)
		assert.Equal(t, 3, capturedPlant.WateringFrequencyDays)

		mockDB.AssertExpectations(t)
	})
}

func TestGetPlant(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		mockDB := new(MockDatabase)
		plant := &model.Plant{ID: 123, UserID: 1, SpeciesID: 456}

		mockDB.On("GetPlantByID", mock.Anything, int64(1), int64(123)).Return(plant, nil)

		req := httptest.NewRequest("GET", "/plants/123", nil)
		req = req.WithContext(context.WithValue(req.Context(), "user", map[string]interface{}{"uid": "1"}))
		req = mux.SetURLVars(req, map[string]string{"plantID": "123"})

		w := httptest.NewRecorder()
		handler := GetPlant(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response model.Plant
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, int64(123), response.ID)
		assert.Equal(t, int64(1), response.UserID)

		mockDB.AssertExpectations(t)
	})

	t.Run("invalid plant ID", func(t *testing.T) {
		mockDB := new(MockDatabase)

		req := httptest.NewRequest("GET", "/plants/invalid", nil)
		req = mux.SetURLVars(req, map[string]string{"plantID": "invalid"})

		w := httptest.NewRecorder()
		handler := GetPlant(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid plant ID")
	})

	t.Run("unauthorized - no user context", func(t *testing.T) {
		mockDB := new(MockDatabase)

		req := httptest.NewRequest("GET", "/plants/123", nil)
		req = mux.SetURLVars(req, map[string]string{"plantID": "123"})

		w := httptest.NewRecorder()
		handler := GetPlant(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("database error", func(t *testing.T) {
		mockDB := new(MockDatabase)
		mockDB.On("GetPlantByID", mock.Anything, int64(1), int64(123)).Return(&model.Plant{}, errors.New("db error"))

		req := httptest.NewRequest("GET", "/plants/123", nil)
		req = req.WithContext(context.WithValue(req.Context(), "user", map[string]interface{}{"uid": "1"}))
		req = mux.SetURLVars(req, map[string]string{"plantID": "123"})

		w := httptest.NewRecorder()
		handler := GetPlant(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to get plant")
		mockDB.AssertExpectations(t)
	})

	t.Run("negative plant ID", func(t *testing.T) {
		mockDB := new(MockDatabase)

		req := httptest.NewRequest("GET", "/plants/-1", nil)
		req = mux.SetURLVars(req, map[string]string{"plantID": "-1"})

		w := httptest.NewRecorder()
		handler := GetPlant(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("plant not found", func(t *testing.T) {
		mockDB := new(MockDatabase)
		mockDB.On("GetPlantByID", mock.Anything, int64(1), int64(999)).Return(nil, nil)

		body := map[string]interface{}{
			"nickname": "Updated Plant",
		}
		req := createTestRequest("PATCH", "/plants/999", body)
		req = req.WithContext(context.WithValue(req.Context(), "user", map[string]interface{}{"uid": "1"}))
		req = mux.SetURLVars(req, map[string]string{"plantID": "999"})

		w := httptest.NewRecorder()
		handler := UpdatePlant(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Plant not found")
		mockDB.AssertExpectations(t)
	})

	t.Run("get plant error", func(t *testing.T) {
		mockDB := new(MockDatabase)
		mockDB.On("GetPlantByID", mock.Anything, int64(1), int64(123)).Return(nil, errors.New("db error"))

		body := map[string]interface{}{
			"nickname": "Updated Plant",
		}
		req := createTestRequest("PATCH", "/plants/123", body)
		req = req.WithContext(context.WithValue(req.Context(), "user", map[string]interface{}{"uid": "1"}))
		req = mux.SetURLVars(req, map[string]string{"plantID": "123"})

		w := httptest.NewRecorder()
		handler := UpdatePlant(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to get plant")
		mockDB.AssertExpectations(t)
	})
}

func TestUpdatePlant(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		mockDB := new(MockDatabase)
		
		// Mock GetPlantByID to return existing plant
		existingPlant := &model.Plant{
			ID:                    123,
			UserID:                1,
			SpeciesID:             1,
			Nickname:              stringPtr("Original Plant"),
			WateringFrequencyDays: 7,
		}
		mockDB.On("GetPlantByID", mock.Anything, int64(1), int64(123)).Return(existingPlant, nil)
		mockDB.On("UpdatePlant", mock.Anything, mock.AnythingOfType("*model.Plant")).Return(nil)

		body := map[string]interface{}{
			"nickname":                "Updated Plant",
			"watering_frequency_days": 14,
		}
		req := createTestRequest("PATCH", "/plants/123", body)
		req = req.WithContext(context.WithValue(req.Context(), "user", map[string]interface{}{"uid": "1"}))
		req = mux.SetURLVars(req, map[string]string{"plantID": "123"})

		w := httptest.NewRecorder()
		handler := UpdatePlant(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		
		// Verify response contains updated plant
		var response model.Plant
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Plant", *response.Nickname)
		assert.Equal(t, 14, response.WateringFrequencyDays)
		
		mockDB.AssertExpectations(t)
	})

	t.Run("invalid plant ID", func(t *testing.T) {
		mockDB := new(MockDatabase)

		body := map[string]string{"nickname": "Updated Plant"}
		req := createTestRequest("PUT", "/plants/invalid", body)
		req = mux.SetURLVars(req, map[string]string{"plantID": "invalid"})

		w := httptest.NewRecorder()
		handler := UpdatePlant(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid plant ID")
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockDB := new(MockDatabase)

		req := httptest.NewRequest("PATCH", "/plants/123", bytes.NewReader([]byte("invalid json")))
		req = req.WithContext(context.WithValue(req.Context(), "user", map[string]interface{}{"uid": "1"}))
		req = mux.SetURLVars(req, map[string]string{"plantID": "123"})

		w := httptest.NewRecorder()
		handler := UpdatePlant(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request body")
	})

	t.Run("unauthorized - no user context", func(t *testing.T) {
		mockDB := new(MockDatabase)

		body := map[string]string{"nickname": "Updated Plant"}
		req := createTestRequest("PATCH", "/plants/123", body)
		req = mux.SetURLVars(req, map[string]string{"plantID": "123"})

		w := httptest.NewRecorder()
		handler := UpdatePlant(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("database error", func(t *testing.T) {
		mockDB := new(MockDatabase)
		
		// Mock GetPlantByID to return existing plant
		existingPlant := &model.Plant{
			ID:                    123,
			UserID:                1,
			SpeciesID:             1,
			Nickname:              stringPtr("Original Plant"),
			WateringFrequencyDays: 7,
		}
		mockDB.On("GetPlantByID", mock.Anything, int64(1), int64(123)).Return(existingPlant, nil)
		mockDB.On("UpdatePlant", mock.Anything, mock.AnythingOfType("*model.Plant")).Return(errors.New("db error"))

		body := map[string]string{"nickname": "Updated Plant"}
		req := createTestRequest("PATCH", "/plants/123", body)
		req = req.WithContext(context.WithValue(req.Context(), "user", map[string]interface{}{"uid": "1"}))
		req = mux.SetURLVars(req, map[string]string{"plantID": "123"})

		w := httptest.NewRecorder()
		handler := UpdatePlant(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to update plant")
		mockDB.AssertExpectations(t)
	})

	t.Run("update with partial data", func(t *testing.T) {
		mockDB := new(MockDatabase)
		
		// Mock GetPlantByID to return existing plant
		existingPlant := &model.Plant{
			ID:                    123,
			UserID:                1,
			SpeciesID:             1,
			Nickname:              stringPtr("Original Plant"),
			WateringFrequencyDays: 7,
		}
		mockDB.On("GetPlantByID", mock.Anything, int64(1), int64(123)).Return(existingPlant, nil)
		mockDB.On("UpdatePlant", mock.Anything, mock.AnythingOfType("*model.Plant")).Return(nil)

		body := map[string]interface{}{
			"nickname": "Only nickname updated",
		}
		req := createTestRequest("PATCH", "/plants/123", body)
		req = req.WithContext(context.WithValue(req.Context(), "user", map[string]interface{}{"uid": "1"}))
		req = mux.SetURLVars(req, map[string]string{"plantID": "123"})

		w := httptest.NewRecorder()
		handler := UpdatePlant(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		// Verify response contains updated plant with only nickname changed
		var response model.Plant
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Only nickname updated", *response.Nickname)
		assert.Equal(t, 7, response.WateringFrequencyDays) // Should remain unchanged
		assert.Equal(t, int64(1), response.SpeciesID)      // Should remain unchanged
		
		mockDB.AssertExpectations(t)
	})
}

func TestDeletePlant(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		mockDB := new(MockDatabase)
		mockDB.On("DeletePlant", mock.Anything, int64(1), int64(123)).Return(nil)

		req := httptest.NewRequest("DELETE", "/plants/123", nil)
		req = req.WithContext(context.WithValue(req.Context(), "user", map[string]interface{}{"uid": "1"}))
		req = mux.SetURLVars(req, map[string]string{"plantID": "123"})

		w := httptest.NewRecorder()
		handler := DeletePlant(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Empty(t, w.Body.String())
		mockDB.AssertExpectations(t)
	})

	t.Run("invalid plant ID", func(t *testing.T) {
		mockDB := new(MockDatabase)

		req := httptest.NewRequest("DELETE", "/plants/invalid", nil)
		req = mux.SetURLVars(req, map[string]string{"plantID": "invalid"})

		w := httptest.NewRecorder()
		handler := DeletePlant(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid plant ID")
	})

	t.Run("unauthorized - no user context", func(t *testing.T) {
		mockDB := new(MockDatabase)

		req := httptest.NewRequest("DELETE", "/plants/123", nil)
		req = mux.SetURLVars(req, map[string]string{"plantID": "123"})

		w := httptest.NewRecorder()
		handler := DeletePlant(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("database error", func(t *testing.T) {
		mockDB := new(MockDatabase)
		mockDB.On("DeletePlant", mock.Anything, int64(1), int64(123)).Return(errors.New("db error"))

		req := httptest.NewRequest("DELETE", "/plants/123", nil)
		req = req.WithContext(context.WithValue(req.Context(), "user", map[string]interface{}{"uid": "1"}))
		req = mux.SetURLVars(req, map[string]string{"plantID": "123"})

		w := httptest.NewRecorder()
		handler := DeletePlant(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to delete plant")
		mockDB.AssertExpectations(t)
	})

	t.Run("zero plant ID", func(t *testing.T) {
		mockDB := new(MockDatabase)

		req := httptest.NewRequest("DELETE", "/plants/0", nil)
		req = mux.SetURLVars(req, map[string]string{"plantID": "0"})

		w := httptest.NewRecorder()
		handler := DeletePlant(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestListPlants(t *testing.T) {
	t.Run("successful list with default pagination", func(t *testing.T) {
		mockDB := new(MockDatabase)
		plants := []*model.Plant{
			{ID: 1, UserID: 1, SpeciesID: 123},
			{ID: 2, UserID: 1, SpeciesID: 456},
		}
		mockDB.On("ListPlants", mock.Anything, int64(1), 20, 0).Return(plants, nil)

		req := httptest.NewRequest("GET", "/plants", nil)
		req = req.WithContext(context.WithValue(req.Context(), "user", map[string]interface{}{"uid": "1"}))

		w := httptest.NewRecorder()
		handler := ListPlants(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response []*model.Plant
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)

		mockDB.AssertExpectations(t)
	})

	t.Run("successful list with custom pagination", func(t *testing.T) {
		mockDB := new(MockDatabase)
		plants := []*model.Plant{}
		mockDB.On("ListPlants", mock.Anything, int64(1), 10, 5).Return(plants, nil)

		req := httptest.NewRequest("GET", "/plants?limit=10&offset=5", nil)
		req = req.WithContext(context.WithValue(req.Context(), "user", map[string]interface{}{"uid": "1"}))

		w := httptest.NewRecorder()
		handler := ListPlants(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockDB.AssertExpectations(t)
	})

	t.Run("invalid pagination parameters ignored", func(t *testing.T) {
		mockDB := new(MockDatabase)
		plants := []*model.Plant{}
		mockDB.On("ListPlants", mock.Anything, int64(1), 20, 0).Return(plants, nil)

		req := httptest.NewRequest("GET", "/plants?limit=invalid&offset=invalid", nil)
		req = req.WithContext(context.WithValue(req.Context(), "user", map[string]interface{}{"uid": "1"}))

		w := httptest.NewRecorder()
		handler := ListPlants(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockDB.AssertExpectations(t)
	})

	t.Run("unauthorized - no user context", func(t *testing.T) {
		mockDB := new(MockDatabase)

		req := httptest.NewRequest("GET", "/plants", nil)

		w := httptest.NewRecorder()
		handler := ListPlants(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("database error", func(t *testing.T) {
		mockDB := new(MockDatabase)
		mockDB.On("ListPlants", mock.Anything, int64(1), 20, 0).Return([]*model.Plant{}, errors.New("db error"))

		req := httptest.NewRequest("GET", "/plants", nil)
		req = req.WithContext(context.WithValue(req.Context(), "user", map[string]interface{}{"uid": "1"}))

		w := httptest.NewRecorder()
		handler := ListPlants(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to list plants")
		mockDB.AssertExpectations(t)
	})

	t.Run("large limit parameter", func(t *testing.T) {
		mockDB := new(MockDatabase)
		plants := []*model.Plant{}
		mockDB.On("ListPlants", mock.Anything, int64(1), 1000, 0).Return(plants, nil)

		req := httptest.NewRequest("GET", "/plants?limit=1000", nil)
		req = req.WithContext(context.WithValue(req.Context(), "user", map[string]interface{}{"uid": "1"}))

		w := httptest.NewRecorder()
		handler := ListPlants(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockDB.AssertExpectations(t)
	})

	t.Run("negative pagination parameters", func(t *testing.T) {
		mockDB := new(MockDatabase)
		plants := []*model.Plant{}
		mockDB.On("ListPlants", mock.Anything, int64(1), 20, 0).Return(plants, nil)

		req := httptest.NewRequest("GET", "/plants?limit=-5&offset=-10", nil)
		req = req.WithContext(context.WithValue(req.Context(), "user", map[string]interface{}{"uid": "1"}))

		w := httptest.NewRecorder()
		handler := ListPlants(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockDB.AssertExpectations(t)
	})

	t.Run("empty plant list", func(t *testing.T) {
		mockDB := new(MockDatabase)
		plants := []*model.Plant{}
		mockDB.On("ListPlants", mock.Anything, int64(1), 20, 0).Return(plants, nil)

		req := httptest.NewRequest("GET", "/plants", nil)
		req = req.WithContext(context.WithValue(req.Context(), "user", map[string]interface{}{"uid": "1"}))

		w := httptest.NewRecorder()
		handler := ListPlants(mockDB)
		handler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []*model.Plant
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Len(t, response, 0)

		mockDB.AssertExpectations(t)
	})
}
