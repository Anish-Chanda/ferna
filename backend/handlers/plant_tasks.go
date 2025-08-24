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

func CreatePlantTask(database db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task model.PlantTask
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		userID, err := GetUserIDFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Validate required fields
		if task.PlantID == 0 {
			http.Error(w, "plant_id is required", http.StatusBadRequest)
			return
		}
		if task.TaskType == "" {
			http.Error(w, "task_type is required", http.StatusBadRequest)
			return
		}
		if task.IntervalDays <= 0 {
			http.Error(w, "interval_days must be greater than 0", http.StatusBadRequest)
			return
		}
		if task.ToleranceDays < 0 {
			http.Error(w, "tolerance_days must be >= 0", http.StatusBadRequest)
			return
		}

		// Verify the plant belongs to the user
		plant, err := database.GetUserPlantByID(r.Context(), userID, task.PlantID)
		if err != nil {
			http.Error(w, "Failed to get plant: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if plant == nil {
			http.Error(w, "Plant not found", http.StatusNotFound)
			return
		}

		now := time.Now()
		task.CreatedAt = now
		task.UpdatedAt = now

		// Set default next due date if not provided
		if task.NextDueAt == nil {
			nextDue := now.AddDate(0, 0, task.IntervalDays)
			task.NextDueAt = &nextDue
		}

		id, err := database.CreatePlantTask(r.Context(), &task)
		if err != nil {
			http.Error(w, "Failed to create task", http.StatusInternalServerError)
			return
		}

		task.ID = id
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(task)
	}
}

func GetPlantTasksByPlantID(database db.Database) http.HandlerFunc {
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

		tasks, err := database.GetPlantTasksByPlantID(r.Context(), userID, plantID)
		if err != nil {
			http.Error(w, "Failed to get tasks: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if tasks == nil {
			tasks = []*model.PlantTask{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	}
}

func GetPlantTask(database db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		taskID, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		userID, err := GetUserIDFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		task, err := database.GetPlantTaskByID(r.Context(), userID, taskID)
		if err != nil {
			http.Error(w, "Failed to get task: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if task == nil {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
	}
}

func UpdatePlantTask(database db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		taskID, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		userID, err := GetUserIDFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		var task model.PlantTask
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Get existing task to verify ownership
		existing, err := database.GetPlantTaskByID(r.Context(), userID, taskID)
		if err != nil {
			http.Error(w, "Failed to get task: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if existing == nil {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		// Set immutable fields
		task.ID = taskID
		task.PlantID = existing.PlantID
		task.CreatedAt = existing.CreatedAt
		task.UpdatedAt = time.Now()

		// Validate required fields
		if task.TaskType == "" {
			task.TaskType = existing.TaskType
		}
		if task.IntervalDays <= 0 {
			task.IntervalDays = existing.IntervalDays
		}
		if task.ToleranceDays < 0 {
			task.ToleranceDays = existing.ToleranceDays
		}

		err = database.UpdatePlantTask(r.Context(), &task)
		if err != nil {
			http.Error(w, "Failed to update task: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
	}
}

func DeletePlantTask(database db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		taskID, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		userID, err := GetUserIDFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		err = database.DeletePlantTask(r.Context(), userID, taskID)
		if err != nil {
			http.Error(w, "Failed to delete task: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func GetOverdueTasks(database db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := GetUserIDFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		tasks, err := database.GetOverdueTasks(r.Context(), userID)
		if err != nil {
			http.Error(w, "Failed to get overdue tasks: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if tasks == nil {
			tasks = []*model.PlantTask{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	}
}
