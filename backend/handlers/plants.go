package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/anish-chanda/ferna/db"
	"github.com/anish-chanda/ferna/model"
	"github.com/gorilla/mux"
)

// Creates a plant for the user. Requires userID from JWT token in request context
func CreatePlant(database db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p model.Plant
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Get user ID
		userID, err := GetUserIDFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		p.UserID = userID

		// validate required fielda
		if p.SpeciesID == 0 {
			http.Error(w, "species_id is required", http.StatusBadRequest)
			return
		}

		// get species
		species, err := database.GetSpeciesByID(r.Context(), p.SpeciesID)
		if err != nil {
			http.Error(w, "Failed to get species: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if species == nil {
			http.Error(w, "Species not found", http.StatusBadRequest)
			return
		}

		// set nickname and wattering frequency if not provided
		if p.Nickname == nil || *p.Nickname == "" {
			p.Nickname = &species.CommonName
		}
		if p.WateringFrequencyDays <= 0 {
			p.WateringFrequencyDays = species.DefaultWateringFrequency
		}

		id, err := database.CreatePlant(r.Context(), &p)
		if err != nil {
			http.Error(w, "Failed to create plant", http.StatusInternalServerError)
			return
		}

		p.ID = id
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(p)
	}
}

func GetPlant(database db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		plantID, err := strconv.ParseInt(mux.Vars(r)["plantID"], 10, 64)
		if err != nil {
			http.Error(w, "Invalid plant ID", http.StatusBadRequest)
			return
		}

		userID, err := GetUserIDFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		plant, err := database.GetPlantByID(r.Context(), userID, plantID)
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

func UpdatePlant(database db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		plantID, err := strconv.ParseInt(mux.Vars(r)["plantID"], 10, 64)
		if err != nil {
			http.Error(w, "Invalid plant ID", http.StatusBadRequest)
			return
		}

		// Parse the partial update request first
		var updateData model.Plant
		if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		userID, err := GetUserIDFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// First, fetch the existing plant
		existingPlant, err := database.GetPlantByID(r.Context(), userID, plantID)
		if err != nil {
			http.Error(w, "Failed to get plant: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if existingPlant == nil {
			http.Error(w, "Plant not found", http.StatusNotFound)
			return
		}

		// Merge changes into existing plant (only update non-zero fields)
		if updateData.SpeciesID != 0 {
			existingPlant.SpeciesID = updateData.SpeciesID
		}
		if updateData.Nickname != nil {
			existingPlant.Nickname = updateData.Nickname
		}
		if updateData.ImageURL != nil {
			existingPlant.ImageURL = updateData.ImageURL
		}
		if updateData.WateringFrequencyDays != 0 {
			existingPlant.WateringFrequencyDays = updateData.WateringFrequencyDays
		}
		if updateData.LastWateredAt != nil {
			existingPlant.LastWateredAt = updateData.LastWateredAt
		}
		if updateData.Note != nil {
			existingPlant.Note = updateData.Note
		}

		// Update the plant
		if err := database.UpdatePlant(r.Context(), existingPlant); err != nil {
			http.Error(w, "Failed to update plant: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Return the updated plant
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(existingPlant)
	}
}

func DeletePlant(database db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		plantID, err := strconv.ParseInt(mux.Vars(r)["plantID"], 10, 64)
		if err != nil {
			http.Error(w, "Invalid plant ID", http.StatusBadRequest)
			return
		}

		userID, err := GetUserIDFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if err := database.DeletePlant(r.Context(), userID, plantID); err != nil {
			http.Error(w, "Failed to delete plant: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func ListPlants(database db.Database) http.HandlerFunc {
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

		plants, err := database.ListPlants(r.Context(), userID, limit, offset)
		if err != nil {
			http.Error(w, "Failed to list plants: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// return empty array if user doesnt have any plants
		if plants == nil {
			plants = []*model.Plant{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(plants)
	}
}
