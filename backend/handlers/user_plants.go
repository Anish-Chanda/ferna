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

// Creates a user plant for the user. Requires userID from JWT token in request context
func CreateUserPlant(database db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var plant model.UserPlant
		if err := json.NewDecoder(r.Body).Decode(&plant); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Get user ID
		userID, err := GetUserIDFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		plant.UserID = userID

		// Validate required fields
		if plant.SpeciesID == 0 {
			http.Error(w, "species_id is required", http.StatusBadRequest)
			return
		}

		// Get species to validate it exists
		species, err := database.GetSpeciesByID(r.Context(), plant.SpeciesID)
		if err != nil {
			http.Error(w, "Failed to get species: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if species == nil {
			http.Error(w, "Species not found", http.StatusBadRequest)
			return
		}

		// Set timestamps
		now := time.Now()
		plant.CreatedAt = now
		plant.UpdatedAt = now

		// Set default nickname if not provided
		if plant.Nickname == nil || *plant.Nickname == "" {
			plant.Nickname = &species.CommonName
		}

		id, err := database.CreateUserPlant(r.Context(), &plant)
		if err != nil {
			http.Error(w, "Failed to create plant", http.StatusInternalServerError)
			return
		}

		plant.ID = id
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(plant)
	}
}

func GetUserPlant(database db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		plantID, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			http.Error(w, "Invalid plant ID", http.StatusBadRequest)
			return
		}

		userID, err := GetUserIDFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		plant, err := database.GetUserPlantByID(r.Context(), userID, plantID)
		if err != nil {
			http.Error(w, "Failed to get plant: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if plant == nil {
			http.Error(w, "Plant not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(plant)
	}
}

func UpdateUserPlant(database db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		plantID, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			http.Error(w, "Invalid plant ID", http.StatusBadRequest)
			return
		}

		userID, err := GetUserIDFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		var plant model.UserPlant
		if err := json.NewDecoder(r.Body).Decode(&plant); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Get existing plant to verify ownership and get current values
		existingPlant, err := database.GetUserPlantByID(r.Context(), userID, plantID)
		if err != nil {
			http.Error(w, "Failed to get plant: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if existingPlant == nil {
			http.Error(w, "Plant not found", http.StatusNotFound)
			return
		}

		// Set immutable fields
		plant.ID = plantID
		plant.UserID = userID
		plant.CreatedAt = existingPlant.CreatedAt
		plant.UpdatedAt = time.Now()

		// Keep existing values for fields not provided
		if plant.SpeciesID == 0 {
			plant.SpeciesID = existingPlant.SpeciesID
		}

		err = database.UpdateUserPlant(r.Context(), &plant)
		if err != nil {
			http.Error(w, "Failed to update plant: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(plant)
	}
}

func DeleteUserPlant(database db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		plantID, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			http.Error(w, "Invalid plant ID", http.StatusBadRequest)
			return
		}

		userID, err := GetUserIDFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		err = database.DeleteUserPlant(r.Context(), userID, plantID)
		if err != nil {
			http.Error(w, "Failed to delete plant: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func ListUserPlants(database db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		plants, err := database.ListUserPlants(r.Context(), userID, limit, offset)
		if err != nil {
			http.Error(w, "Failed to list plants: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Return empty array if user doesn't have any plants
		if plants == nil {
			plants = []*model.UserPlant{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(plants)
	}
}
