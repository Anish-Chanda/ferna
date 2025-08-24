package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/anish-chanda/ferna/model"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockDB implements db.Database for testing all handlers
type mockDB struct {
	// Species methods
	SearchSpeciesFunc  func(ctx context.Context, query string, limit, offset int) ([]*model.Species, error)
	GetSpeciesByIDFunc func(ctx context.Context, speciesID int64) (*model.Species, error)

	// Location methods
	CreateLocationFunc    func(ctx context.Context, location *model.Location) (int64, error)
	ListLocationsFunc     func(ctx context.Context, userID int64) ([]*model.Location, error)
	GetLocationByIDFunc   func(ctx context.Context, userID, locationID int64) (*model.Location, error)
	UpdateLocationFunc    func(ctx context.Context, location *model.Location) error
	DeleteLocationFunc    func(ctx context.Context, userID, locationID int64) error

	// User plant methods
	CreateUserPlantFunc func(ctx context.Context, plant *model.UserPlant) (int64, error)
	GetUserPlantByIDFunc func(ctx context.Context, userID, plantID int64) (*model.UserPlant, error)
	ListUserPlantsFunc  func(ctx context.Context, userID int64, limit, offset int) ([]*model.UserPlant, error)
	UpdateUserPlantFunc func(ctx context.Context, plant *model.UserPlant) error
	DeleteUserPlantFunc func(ctx context.Context, userID, plantID int64) error

	// Plant task methods
	CreatePlantTaskFunc       func(ctx context.Context, task *model.PlantTask) (int64, error)
	GetPlantTasksByPlantIDFunc func(ctx context.Context, userID, plantID int64) ([]*model.PlantTask, error)
	GetPlantTaskByIDFunc      func(ctx context.Context, userID, taskID int64) (*model.PlantTask, error)
	UpdatePlantTaskFunc       func(ctx context.Context, task *model.PlantTask) error
	DeletePlantTaskFunc       func(ctx context.Context, userID, taskID int64) error
	GetOverdueTasksFunc       func(ctx context.Context, userID int64) ([]*model.PlantTask, error)

	// Care event methods
	CreateCareEventFunc        func(ctx context.Context, event *model.CareEvent) (int64, error)
	GetCareEventsByPlantIDFunc func(ctx context.Context, userID, plantID int64, limit, offset int) ([]*model.CareEvent, error)
	GetCareEventByIDFunc       func(ctx context.Context, userID, eventID int64) (*model.CareEvent, error)
	UpdateCareEventFunc        func(ctx context.Context, event *model.CareEvent) error
	DeleteCareEventFunc        func(ctx context.Context, userID, eventID int64) error
}

// Species methods
func (m *mockDB) SearchSpecies(ctx context.Context, query string, limit, offset int) ([]*model.Species, error) {
	if m.SearchSpeciesFunc != nil {
		return m.SearchSpeciesFunc(ctx, query, limit, offset)
	}
	return nil, nil
}

func (m *mockDB) GetSpeciesByID(ctx context.Context, speciesID int64) (*model.Species, error) {
	if m.GetSpeciesByIDFunc != nil {
		return m.GetSpeciesByIDFunc(ctx, speciesID)
	}
	return nil, nil
}

// Location methods
func (m *mockDB) CreateLocation(ctx context.Context, location *model.Location) (int64, error) {
	if m.CreateLocationFunc != nil {
		return m.CreateLocationFunc(ctx, location)
	}
	return 0, nil
}

func (m *mockDB) ListLocations(ctx context.Context, userID int64) ([]*model.Location, error) {
	if m.ListLocationsFunc != nil {
		return m.ListLocationsFunc(ctx, userID)
	}
	return nil, nil
}

func (m *mockDB) GetLocationByID(ctx context.Context, userID, locationID int64) (*model.Location, error) {
	if m.GetLocationByIDFunc != nil {
		return m.GetLocationByIDFunc(ctx, userID, locationID)
	}
	return nil, nil
}

func (m *mockDB) UpdateLocation(ctx context.Context, location *model.Location) error {
	if m.UpdateLocationFunc != nil {
		return m.UpdateLocationFunc(ctx, location)
	}
	return nil
}

func (m *mockDB) DeleteLocation(ctx context.Context, userID, locationID int64) error {
	if m.DeleteLocationFunc != nil {
		return m.DeleteLocationFunc(ctx, userID, locationID)
	}
	return nil
}

// User plant methods
func (m *mockDB) CreateUserPlant(ctx context.Context, plant *model.UserPlant) (int64, error) {
	if m.CreateUserPlantFunc != nil {
		return m.CreateUserPlantFunc(ctx, plant)
	}
	return 0, nil
}

func (m *mockDB) GetUserPlantByID(ctx context.Context, userID, plantID int64) (*model.UserPlant, error) {
	if m.GetUserPlantByIDFunc != nil {
		return m.GetUserPlantByIDFunc(ctx, userID, plantID)
	}
	return nil, nil
}

func (m *mockDB) ListUserPlants(ctx context.Context, userID int64, limit, offset int) ([]*model.UserPlant, error) {
	if m.ListUserPlantsFunc != nil {
		return m.ListUserPlantsFunc(ctx, userID, limit, offset)
	}
	return nil, nil
}

func (m *mockDB) UpdateUserPlant(ctx context.Context, plant *model.UserPlant) error {
	if m.UpdateUserPlantFunc != nil {
		return m.UpdateUserPlantFunc(ctx, plant)
	}
	return nil
}

func (m *mockDB) DeleteUserPlant(ctx context.Context, userID, plantID int64) error {
	if m.DeleteUserPlantFunc != nil {
		return m.DeleteUserPlantFunc(ctx, userID, plantID)
	}
	return nil
}

// Plant task methods
func (m *mockDB) CreatePlantTask(ctx context.Context, task *model.PlantTask) (int64, error) {
	if m.CreatePlantTaskFunc != nil {
		return m.CreatePlantTaskFunc(ctx, task)
	}
	return 0, nil
}

func (m *mockDB) GetPlantTasksByPlantID(ctx context.Context, userID, plantID int64) ([]*model.PlantTask, error) {
	if m.GetPlantTasksByPlantIDFunc != nil {
		return m.GetPlantTasksByPlantIDFunc(ctx, userID, plantID)
	}
	return nil, nil
}

func (m *mockDB) GetPlantTaskByID(ctx context.Context, userID, taskID int64) (*model.PlantTask, error) {
	if m.GetPlantTaskByIDFunc != nil {
		return m.GetPlantTaskByIDFunc(ctx, userID, taskID)
	}
	return nil, nil
}

func (m *mockDB) UpdatePlantTask(ctx context.Context, task *model.PlantTask) error {
	if m.UpdatePlantTaskFunc != nil {
		return m.UpdatePlantTaskFunc(ctx, task)
	}
	return nil
}

func (m *mockDB) DeletePlantTask(ctx context.Context, userID, taskID int64) error {
	if m.DeletePlantTaskFunc != nil {
		return m.DeletePlantTaskFunc(ctx, userID, taskID)
	}
	return nil
}

func (m *mockDB) GetOverdueTasks(ctx context.Context, userID int64) ([]*model.PlantTask, error) {
	if m.GetOverdueTasksFunc != nil {
		return m.GetOverdueTasksFunc(ctx, userID)
	}
	return nil, nil
}

// Care event methods
func (m *mockDB) CreateCareEvent(ctx context.Context, event *model.CareEvent) (int64, error) {
	if m.CreateCareEventFunc != nil {
		return m.CreateCareEventFunc(ctx, event)
	}
	return 0, nil
}

func (m *mockDB) GetCareEventsByPlantID(ctx context.Context, userID, plantID int64, limit, offset int) ([]*model.CareEvent, error) {
	if m.GetCareEventsByPlantIDFunc != nil {
		return m.GetCareEventsByPlantIDFunc(ctx, userID, plantID, limit, offset)
	}
	return nil, nil
}

func (m *mockDB) GetCareEventByID(ctx context.Context, userID, eventID int64) (*model.CareEvent, error) {
	if m.GetCareEventByIDFunc != nil {
		return m.GetCareEventByIDFunc(ctx, userID, eventID)
	}
	return nil, nil
}

func (m *mockDB) UpdateCareEvent(ctx context.Context, event *model.CareEvent) error {
	if m.UpdateCareEventFunc != nil {
		return m.UpdateCareEventFunc(ctx, event)
	}
	return nil
}

func (m *mockDB) DeleteCareEvent(ctx context.Context, userID, eventID int64) error {
	if m.DeleteCareEventFunc != nil {
		return m.DeleteCareEventFunc(ctx, userID, eventID)
	}
	return nil
}

// Unused auth methods to satisfy interface
func (m *mockDB) CheckIfEmailExists(ctx context.Context, email string) (bool, error) { return false, nil }
func (m *mockDB) CreateUser(ctx context.Context, user *model.User) (int64, error)     { return 0, nil }
func (m *mockDB) GetUserByEmail(ctx context.Context, email string) (*model.User, error) { return nil, nil }
func (m *mockDB) GetUserByID(ctx context.Context, userID int64) (*model.User, error) { return nil, nil }

// Infrastructure methods
func (m *mockDB) Connect(dsn string) error { return nil }
func (m *mockDB) Close() error             { return nil }
func (m *mockDB) Migrate() error           { return nil }

// Helper functions for tests
func stringPtr(s string) *string {
	return &s
}

func createRequestWithUserContext(method, url string, body []byte, userID int64) *http.Request {
	req := httptest.NewRequest(method, url, bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), "user", map[string]interface{}{
		"uid": strconv.FormatInt(userID, 10),
	})
	return req.WithContext(ctx)
}

// Species handler tests
func TestSearchSpecies(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		limit          string
		offset         string
		mockSpecies    []*model.Species
		mockError      error
		expectedStatus int
		expectedCount  int
	}{
		{
			name:  "successful search",
			query: "Rose",
			mockSpecies: []*model.Species{
				{ID: 1, CommonName: "Rose", ScientificName: stringPtr("Rosa rubiginosa")},
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:           "empty search",
			query:          "",
			mockSpecies:    []*model.Species{},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name:           "database error",
			query:          "Rose",
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "with limit and offset",
			query:  "Plant",
			limit:  "10",
			offset: "5",
			mockSpecies: []*model.Species{
				{ID: 1, CommonName: "Plant 1"},
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &mockDB{
				SearchSpeciesFunc: func(ctx context.Context, query string, limit, offset int) ([]*model.Species, error) {
					assert.Equal(t, tt.query, query)
					if tt.limit != "" {
						expectedLimit, _ := strconv.Atoi(tt.limit)
						assert.Equal(t, expectedLimit, limit)
					} else {
						assert.Equal(t, 20, limit) // default
					}
					if tt.offset != "" {
						expectedOffset, _ := strconv.Atoi(tt.offset)
						assert.Equal(t, expectedOffset, offset)
					} else {
						assert.Equal(t, 0, offset) // default
					}
					return tt.mockSpecies, tt.mockError
				},
			}

			url := "/api/species?query=" + tt.query
			if tt.limit != "" {
				url += "&limit=" + tt.limit
			}
			if tt.offset != "" {
				url += "&offset=" + tt.offset
			}

			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			SearchSpecies(mockDB)(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var species []*model.Species
				err := json.NewDecoder(w.Body).Decode(&species)
				require.NoError(t, err)
				assert.Len(t, species, tt.expectedCount)
			}
		})
	}
}

// Location handler tests
func TestCreateLocation(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    model.Location
		userID         int64
		mockError      error
		expectedStatus int
		expectedID     int64
	}{
		{
			name: "successful creation",
			requestBody: model.Location{
				Name: "Living Room",
			},
			userID:         123,
			expectedStatus: http.StatusCreated,
			expectedID:     1,
		},
		{
			name: "missing name",
			requestBody: model.Location{
				// No name provided
			},
			userID:         123,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "database error",
			requestBody: model.Location{
				Name: "Bedroom",
			},
			userID:         123,
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &mockDB{
				CreateLocationFunc: func(ctx context.Context, location *model.Location) (int64, error) {
					assert.Equal(t, tt.userID, location.UserID)
					assert.Equal(t, tt.requestBody.Name, location.Name)
					assert.NotZero(t, location.CreatedAt)
					assert.NotZero(t, location.UpdatedAt)
					return tt.expectedID, tt.mockError
				},
			}

			body, _ := json.Marshal(tt.requestBody)
			req := createRequestWithUserContext("POST", "/api/locations", body, tt.userID)
			w := httptest.NewRecorder()

			CreateLocation(mockDB)(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				var location model.Location
				err := json.NewDecoder(w.Body).Decode(&location)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedID, location.ID)
				assert.Equal(t, tt.userID, location.UserID)
			}
		})
	}
}

func TestListLocations(t *testing.T) {
	userID := int64(123)
	
	tests := []struct {
		name           string
		mockLocations  []*model.Location
		mockError      error
		expectedStatus int
		expectedCount  int
	}{
		{
			name: "successful list",
			mockLocations: []*model.Location{
				{ID: 1, Name: "Living Room", UserID: userID},
				{ID: 2, Name: "Bedroom", UserID: userID},
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "empty list",
			mockLocations:  nil,
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name:           "database error",
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &mockDB{
				ListLocationsFunc: func(ctx context.Context, uid int64) ([]*model.Location, error) {
					assert.Equal(t, userID, uid)
					return tt.mockLocations, tt.mockError
				},
			}

			req := createRequestWithUserContext("GET", "/api/locations", nil, userID)
			w := httptest.NewRecorder()

			ListLocations(mockDB)(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var locations []*model.Location
				err := json.NewDecoder(w.Body).Decode(&locations)
				require.NoError(t, err)
				assert.Len(t, locations, tt.expectedCount)
			}
		})
	}
}

func TestGetLocation(t *testing.T) {
	userID := int64(123)
	locationID := int64(456)

	tests := []struct {
		name           string
		mockLocation   *model.Location
		mockError      error
		expectedStatus int
	}{
		{
			name: "successful get",
			mockLocation: &model.Location{
				ID:     locationID,
				UserID: userID,
				Name:   "Living Room",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "location not found",
			mockLocation:   nil,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "database error",
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &mockDB{
				GetLocationByIDFunc: func(ctx context.Context, uid, lid int64) (*model.Location, error) {
					assert.Equal(t, userID, uid)
					assert.Equal(t, locationID, lid)
					return tt.mockLocation, tt.mockError
				},
			}

			req := createRequestWithUserContext("GET", "/api/locations/456", nil, userID)
			req = mux.SetURLVars(req, map[string]string{"id": "456"})
			w := httptest.NewRecorder()

			GetLocation(mockDB)(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var location model.Location
				err := json.NewDecoder(w.Body).Decode(&location)
				require.NoError(t, err)
				assert.Equal(t, locationID, location.ID)
			}
		})
	}
}

// User plant handler tests
func TestCreateUserPlant(t *testing.T) {
	userID := int64(123)
	speciesID := int64(456)

	tests := []struct {
		name           string
		requestBody    model.UserPlant
		mockSpecies    *model.Species
		mockError      error
		expectedStatus int
		expectedID     int64
	}{
		{
			name: "successful creation",
			requestBody: model.UserPlant{
				SpeciesID: speciesID,
				Nickname:  stringPtr("My Rose"),
			},
			mockSpecies: &model.Species{
				ID:         speciesID,
				CommonName: "Rose",
			},
			expectedStatus: http.StatusCreated,
			expectedID:     1,
		},
		{
			name: "missing species_id",
			requestBody: model.UserPlant{
				Nickname: stringPtr("Plant"),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "species not found",
			requestBody: model.UserPlant{
				SpeciesID: speciesID,
			},
			mockSpecies:    nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "database error",
			requestBody: model.UserPlant{
				SpeciesID: speciesID,
			},
			mockSpecies: &model.Species{
				ID:         speciesID,
				CommonName: "Rose",
			},
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &mockDB{
				GetSpeciesByIDFunc: func(ctx context.Context, sid int64) (*model.Species, error) {
					if tt.requestBody.SpeciesID != 0 {
						assert.Equal(t, speciesID, sid)
					}
					return tt.mockSpecies, nil
				},
				CreateUserPlantFunc: func(ctx context.Context, plant *model.UserPlant) (int64, error) {
					assert.Equal(t, userID, plant.UserID)
					assert.Equal(t, tt.requestBody.SpeciesID, plant.SpeciesID)
					assert.NotZero(t, plant.CreatedAt)
					assert.NotZero(t, plant.UpdatedAt)
					return tt.expectedID, tt.mockError
				},
			}

			body, _ := json.Marshal(tt.requestBody)
			req := createRequestWithUserContext("POST", "/api/plants", body, userID)
			w := httptest.NewRecorder()

			CreateUserPlant(mockDB)(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				var plant model.UserPlant
				err := json.NewDecoder(w.Body).Decode(&plant)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedID, plant.ID)
				assert.Equal(t, userID, plant.UserID)
			}
		})
	}
}

func TestGetUserPlant(t *testing.T) {
	userID := int64(123)
	plantID := int64(456)

	tests := []struct {
		name           string
		mockPlant      *model.UserPlant
		mockError      error
		expectedStatus int
	}{
		{
			name: "successful get",
			mockPlant: &model.UserPlant{
				ID:        plantID,
				UserID:    userID,
				SpeciesID: 789,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "plant not found",
			mockPlant:      nil,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "database error",
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &mockDB{
				GetUserPlantByIDFunc: func(ctx context.Context, uid, pid int64) (*model.UserPlant, error) {
					assert.Equal(t, userID, uid)
					assert.Equal(t, plantID, pid)
					return tt.mockPlant, tt.mockError
				},
			}

			req := createRequestWithUserContext("GET", "/api/plants/456", nil, userID)
			req = mux.SetURLVars(req, map[string]string{"id": "456"})
			w := httptest.NewRecorder()

			GetUserPlant(mockDB)(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var plant model.UserPlant
				err := json.NewDecoder(w.Body).Decode(&plant)
				require.NoError(t, err)
				assert.Equal(t, plantID, plant.ID)
			}
		})
	}
}

// Plant task handler tests
func TestCreatePlantTask(t *testing.T) {
	userID := int64(123)
	plantID := int64(456)

	tests := []struct {
		name           string
		requestBody    model.PlantTask
		mockPlant      *model.UserPlant
		mockError      error
		expectedStatus int
		expectedID     int64
	}{
		{
			name: "successful creation",
			requestBody: model.PlantTask{
				PlantID:        plantID,
				TaskType:       model.TaskTypeWatering,
				IntervalDays:   7,
				ToleranceDays:  2,
			},
			mockPlant: &model.UserPlant{
				ID:     plantID,
				UserID: userID,
			},
			expectedStatus: http.StatusCreated,
			expectedID:     1,
		},
		{
			name: "missing plant_id",
			requestBody: model.PlantTask{
				TaskType:     model.TaskTypeWatering,
				IntervalDays: 7,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing task_type",
			requestBody: model.PlantTask{
				PlantID:      plantID,
				IntervalDays: 7,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid interval_days",
			requestBody: model.PlantTask{
				PlantID:      plantID,
				TaskType:     model.TaskTypeWatering,
				IntervalDays: 0,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "plant not found",
			requestBody: model.PlantTask{
				PlantID:      plantID,
				TaskType:     model.TaskTypeWatering,
				IntervalDays: 7,
			},
			mockPlant:      nil,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &mockDB{
				GetUserPlantByIDFunc: func(ctx context.Context, uid, pid int64) (*model.UserPlant, error) {
					if tt.requestBody.PlantID != 0 {
						assert.Equal(t, userID, uid)
						assert.Equal(t, plantID, pid)
					}
					return tt.mockPlant, nil
				},
				CreatePlantTaskFunc: func(ctx context.Context, task *model.PlantTask) (int64, error) {
					assert.Equal(t, tt.requestBody.PlantID, task.PlantID)
					assert.Equal(t, tt.requestBody.TaskType, task.TaskType)
					assert.Equal(t, tt.requestBody.IntervalDays, task.IntervalDays)
					assert.NotZero(t, task.CreatedAt)
					assert.NotZero(t, task.UpdatedAt)
					return tt.expectedID, tt.mockError
				},
			}

			body, _ := json.Marshal(tt.requestBody)
			req := createRequestWithUserContext("POST", "/api/plants/"+strconv.FormatInt(plantID, 10)+"/tasks", body, userID)
			w := httptest.NewRecorder()

			CreatePlantTask(mockDB)(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				var task model.PlantTask
				err := json.NewDecoder(w.Body).Decode(&task)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedID, task.ID)
			}
		})
	}
}

func TestGetPlantTasksByPlantID(t *testing.T) {
	userID := int64(123)
	plantID := int64(456)

	tests := []struct {
		name           string
		mockTasks      []*model.PlantTask
		mockError      error
		expectedStatus int
		expectedCount  int
	}{
		{
			name: "successful get",
			mockTasks: []*model.PlantTask{
				{ID: 1, PlantID: plantID, TaskType: model.TaskTypeWatering},
				{ID: 2, PlantID: plantID, TaskType: model.TaskTypeFertilizer},
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "empty list",
			mockTasks:      []*model.PlantTask{},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name:           "database error",
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &mockDB{
				GetPlantTasksByPlantIDFunc: func(ctx context.Context, uid, pid int64) ([]*model.PlantTask, error) {
					assert.Equal(t, userID, uid)
					assert.Equal(t, plantID, pid)
					return tt.mockTasks, tt.mockError
				},
			}

			req := createRequestWithUserContext("GET", "/api/plants/456/tasks", nil, userID)
			req = mux.SetURLVars(req, map[string]string{"plantId": "456"})
			w := httptest.NewRecorder()

			GetPlantTasksByPlantID(mockDB)(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var tasks []*model.PlantTask
				err := json.NewDecoder(w.Body).Decode(&tasks)
				require.NoError(t, err)
				assert.Len(t, tasks, tt.expectedCount)
			}
		})
	}
}

// Care event handler tests  
func TestCreateCareEvent(t *testing.T) {
	userID := int64(123)
	plantID := int64(456)

	tests := []struct {
		name           string
		requestBody    model.CareEvent
		mockPlant      *model.UserPlant
		mockError      error
		expectedStatus int
		expectedID     int64
	}{
		{
			name: "successful creation",
			requestBody: model.CareEvent{
				PlantID:   plantID,
				EventType: model.EventTypeWatering,
				Notes:     stringPtr("Watered the plant"),
			},
			mockPlant: &model.UserPlant{
				ID:     plantID,
				UserID: userID,
			},
			expectedStatus: http.StatusCreated,
			expectedID:     1,
		},
		{
			name: "missing plant_id",
			requestBody: model.CareEvent{
				EventType: model.EventTypeWatering,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing event_type",
			requestBody: model.CareEvent{
				PlantID: plantID,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "plant not found",
			requestBody: model.CareEvent{
				PlantID:   plantID,
				EventType: model.EventTypeWatering,
			},
			mockPlant:      nil,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &mockDB{
				GetUserPlantByIDFunc: func(ctx context.Context, uid, pid int64) (*model.UserPlant, error) {
					if tt.requestBody.PlantID != 0 {
						assert.Equal(t, userID, uid)
						assert.Equal(t, plantID, pid)
					}
					return tt.mockPlant, nil
				},
				CreateCareEventFunc: func(ctx context.Context, event *model.CareEvent) (int64, error) {
					assert.Equal(t, tt.requestBody.PlantID, event.PlantID)
					assert.Equal(t, tt.requestBody.EventType, event.EventType)
					assert.NotZero(t, event.CreatedAt)
					assert.NotZero(t, event.UpdatedAt)
					return tt.expectedID, tt.mockError
				},
			}

			body, _ := json.Marshal(tt.requestBody)
			req := createRequestWithUserContext("POST", "/api/plants/"+strconv.FormatInt(plantID, 10)+"/events", body, userID)
			w := httptest.NewRecorder()

			CreateCareEvent(mockDB)(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				var event model.CareEvent
				err := json.NewDecoder(w.Body).Decode(&event)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedID, event.ID)
			}
		})
	}
}

// Additional handler tests to increase coverage

// Update and Delete Location tests
func TestUpdateLocation(t *testing.T) {
	userID := int64(123)
	locationID := int64(456)

	tests := []struct {
		name           string
		requestBody    model.Location
		mockLocation   *model.Location
		mockError      error
		expectedStatus int
	}{
		{
			name: "successful update",
			requestBody: model.Location{
				Name: "Updated Living Room",
			},
			mockLocation: &model.Location{
				ID:     locationID,
				UserID: userID,
				Name:   "Living Room",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "missing name",
			requestBody: model.Location{
				Name: "",
			},
			mockLocation: &model.Location{
				ID:     locationID,
				UserID: userID,
				Name:   "Living Room",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "location not found",
			requestBody: model.Location{
				Name: "Updated Room",
			},
			mockLocation:   nil,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &mockDB{
				GetLocationByIDFunc: func(ctx context.Context, uid, lid int64) (*model.Location, error) {
					assert.Equal(t, userID, uid)
					assert.Equal(t, locationID, lid)
					return tt.mockLocation, nil
				},
				UpdateLocationFunc: func(ctx context.Context, location *model.Location) error {
					return tt.mockError
				},
			}

			body, _ := json.Marshal(tt.requestBody)
			req := createRequestWithUserContext("PUT", "/api/locations/456", body, userID)
			req = mux.SetURLVars(req, map[string]string{"id": "456"})
			w := httptest.NewRecorder()

			UpdateLocation(mockDB)(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestDeleteLocation(t *testing.T) {
	userID := int64(123)
	locationID := int64(456)

	tests := []struct {
		name           string
		mockError      error
		expectedStatus int
	}{
		{
			name:           "successful delete",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "database error",
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &mockDB{
				DeleteLocationFunc: func(ctx context.Context, uid, lid int64) error {
					assert.Equal(t, userID, uid)
					assert.Equal(t, locationID, lid)
					return tt.mockError
				},
			}

			req := createRequestWithUserContext("DELETE", "/api/locations/456", nil, userID)
			req = mux.SetURLVars(req, map[string]string{"id": "456"})
			w := httptest.NewRecorder()

			DeleteLocation(mockDB)(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

// Helper function tests
func TestGetUserIDFromRequest(t *testing.T) {
	tests := []struct {
		name        string
		userContext interface{}
		expectError bool
		expectedID  int64
	}{
		{
			name: "valid user context",
			userContext: map[string]interface{}{
				"uid": "123",
			},
			expectError: false,
			expectedID:  123,
		},
		{
			name:        "no user context",
			userContext: nil,
			expectError: true,
		},
		{
			name: "invalid user context type",
			userContext: "invalid",
			expectError: true,
		},
		{
			name: "missing uid in context",
			userContext: map[string]interface{}{
				"other": "value",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.userContext != nil {
				ctx := context.WithValue(req.Context(), "user", tt.userContext)
				req = req.WithContext(ctx)
			}

			userID, err := GetUserIDFromRequest(req)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, userID)
			}
		})
	}
}
