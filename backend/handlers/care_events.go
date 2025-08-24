package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/anish-chanda/ferna/db"
	"github.com/anish-chanda/ferna/model"
	"github.com/gorilla/mux"
)

func CreateCareEvent(database db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var event model.CareEvent
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		userID, err := GetUserIDFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Validate required fields
		if event.PlantID == 0 {
			http.Error(w, "plant_id is required", http.StatusBadRequest)
			return
		}
		if event.EventType == "" {
			http.Error(w, "event_type is required", http.StatusBadRequest)
			return
		}

		// Verify the plant belongs to the user
		plant, err := database.GetUserPlantByID(r.Context(), userID, event.PlantID)
		if err != nil {
			http.Error(w, "Failed to get plant: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if plant == nil {
			http.Error(w, "Plant not found", http.StatusNotFound)
			return
		}

		// If task_id is provided, verify it belongs to the plant
		if event.TaskID != nil {
			task, err := database.GetPlantTaskByID(r.Context(), userID, *event.TaskID)
			if err != nil {
				http.Error(w, "Failed to get task: "+err.Error(), http.StatusInternalServerError)
				return
			}
			if task == nil || task.PlantID != event.PlantID {
				http.Error(w, "Task not found or doesn't belong to this plant", http.StatusBadRequest)
				return
			}
		}

		now := time.Now()
		event.CreatedAt = now
		event.UpdatedAt = now

		// Set default happened_at if not provided
		if event.HappenedAt.IsZero() {
			event.HappenedAt = now
		}

		id, err := database.CreateCareEvent(r.Context(), &event)
		if err != nil {
			http.Error(w, "Failed to create care event", http.StatusInternalServerError)
			return
		}

		event.ID = id
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(event)
	}
}

func GetCareEventsByPlantID(database db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		plantID, err := strconv.ParseInt(vars["plantId"], 10, 64)
		if err != nil {
			http.Error(w, "Invalid plant ID", http.StatusBadRequest)
			return
		}

		userID, err := GetUserIDFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		limit := 20
		if l := r.URL.Query().Get("limit"); l != "" {
			if v, err := strconv.Atoi(l); err == nil && v > 0 {
				limit = v
			}
		}

		offset := 0
		if o := r.URL.Query().Get("offset"); o != "" {
			if v, err := strconv.Atoi(o); err == nil && v >= 0 {
				offset = v
			}
		}

		events, err := database.GetCareEventsByPlantID(r.Context(), userID, plantID, limit, offset)
		if err != nil {
			http.Error(w, "Failed to get care events: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if events == nil {
			events = []*model.CareEvent{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(events)
	}
}

func GetCareEvent(database db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		eventID, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			http.Error(w, "Invalid event ID", http.StatusBadRequest)
			return
		}

		userID, err := GetUserIDFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		event, err := database.GetCareEventByID(r.Context(), userID, eventID)
		if err != nil {
			http.Error(w, "Failed to get care event: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if event == nil {
			http.Error(w, "Care event not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(event)
	}
}

func UpdateCareEvent(database db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		eventID, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			http.Error(w, "Invalid event ID", http.StatusBadRequest)
			return
		}

		userID, err := GetUserIDFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		var event model.CareEvent
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Get existing event to verify ownership
		existing, err := database.GetCareEventByID(r.Context(), userID, eventID)
		if err != nil {
			http.Error(w, "Failed to get care event: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if existing == nil {
			http.Error(w, "Care event not found", http.StatusNotFound)
			return
		}

		// Set immutable fields
		event.ID = eventID
		event.PlantID = existing.PlantID
		event.CreatedAt = existing.CreatedAt
		event.UpdatedAt = time.Now()

		// Keep existing values for fields not provided
		if event.EventType == "" {
			event.EventType = existing.EventType
		}
		if event.HappenedAt.IsZero() {
			event.HappenedAt = existing.HappenedAt
		}

		// If task_id is provided, verify it belongs to the plant
		if event.TaskID != nil {
			task, err := database.GetPlantTaskByID(r.Context(), userID, *event.TaskID)
			if err != nil {
				http.Error(w, "Failed to get task: "+err.Error(), http.StatusInternalServerError)
				return
			}
			if task == nil || task.PlantID != event.PlantID {
				http.Error(w, "Task not found or doesn't belong to this plant", http.StatusBadRequest)
				return
			}
		}

		err = database.UpdateCareEvent(r.Context(), &event)
		if err != nil {
			http.Error(w, "Failed to update care event: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(event)
	}
}

func DeleteCareEvent(database db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		eventID, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			http.Error(w, "Invalid event ID", http.StatusBadRequest)
			return
		}

		userID, err := GetUserIDFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		err = database.DeleteCareEvent(r.Context(), userID, eventID)
		if err != nil {
			http.Error(w, "Failed to delete care event: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
