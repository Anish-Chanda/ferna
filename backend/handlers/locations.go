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

func CreateLocation(database db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var location model.Location
		if err := json.NewDecoder(r.Body).Decode(&location); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		userID, err := GetUserIDFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		location.UserID = userID

		if location.Name == "" {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}

		now := time.Now()
		location.CreatedAt = now
		location.UpdatedAt = now

		id, err := database.CreateLocation(r.Context(), &location)
		if err != nil {
			http.Error(w, "Failed to create location", http.StatusInternalServerError)
			return
		}

		location.ID = id
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(location)
	}
}

func ListLocations(database db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := GetUserIDFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		locations, err := database.ListLocations(r.Context(), userID)
		if err != nil {
			http.Error(w, "Failed to list locations: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if locations == nil {
			locations = []*model.Location{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(locations)
	}
}

func GetLocation(database db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		locationID, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			http.Error(w, "Invalid location ID", http.StatusBadRequest)
			return
		}

		userID, err := GetUserIDFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		location, err := database.GetLocationByID(r.Context(), userID, locationID)
		if err != nil {
			http.Error(w, "Failed to get location: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if location == nil {
			http.Error(w, "Location not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(location)
	}
}

func UpdateLocation(database db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		locationID, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			http.Error(w, "Invalid location ID", http.StatusBadRequest)
			return
		}

		userID, err := GetUserIDFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		var location model.Location
		if err := json.NewDecoder(r.Body).Decode(&location); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Get existing location to verify ownership
		existing, err := database.GetLocationByID(r.Context(), userID, locationID)
		if err != nil {
			http.Error(w, "Failed to get location: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if existing == nil {
			http.Error(w, "Location not found", http.StatusNotFound)
			return
		}

		// Set immutable fields
		location.ID = locationID
		location.UserID = userID
		location.CreatedAt = existing.CreatedAt
		location.UpdatedAt = time.Now()

		if location.Name == "" {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}

		err = database.UpdateLocation(r.Context(), &location)
		if err != nil {
			http.Error(w, "Failed to update location: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(location)
	}
}

func DeleteLocation(database db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		locationID, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			http.Error(w, "Invalid location ID", http.StatusBadRequest)
			return
		}

		userID, err := GetUserIDFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		err = database.DeleteLocation(r.Context(), userID, locationID)
		if err != nil {
			http.Error(w, "Failed to delete location: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
