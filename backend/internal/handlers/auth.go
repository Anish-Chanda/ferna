package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/anish-chanda/ferna/internal/auth"
	"github.com/anish-chanda/ferna/internal/db"
	"github.com/anish-chanda/ferna/internal/logger"
	"github.com/anish-chanda/ferna/model"
	"github.com/google/uuid"
)

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Timezone string `json:"timezone"`
}

type SignupResponse struct {
	Success bool      `json:"success"`
	UserID  uuid.UUID `json:"user_id"`
	Message string    `json:"message"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool      `json:"success"`
	UserID  uuid.UUID `json:"user_id,omitempty"`
	Message string    `json:"message"`
}

// SignupHandler handles user registration for local auth provider
func SignupHandler(database *db.PostgresDB, logger *logger.ServiceLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		// Parse request body
		var req SignupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Debugf("Invalid JSON in signup request: %v", err)
			writeErrorResponse(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		// Validate required fields
		if err := validateSignupRequest(req); err != nil {
			logger.Debugf("Signup validation failed: %v", err)
			writeErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Normalize email
		email := strings.TrimSpace(strings.ToLower(req.Email))

		// Check if email already exists
		exists, err := database.CheckIfEmailExists(ctx, email)
		if err != nil {
			logger.Debugf("Database error checking email existence: %v", err)
			writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if exists {
			writeErrorResponse(w, "Email already registered", http.StatusConflict)
			return
		}

		// Hash password
		hashedPassword, err := auth.HashPassword(req.Password)
		if err != nil {
			logger.Debugf("Password hashing failed: %v", err)
			writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Generate user ID
		userID, err := auth.GenerateUserID()
		if err != nil {
			logger.Debugf("User ID generation failed: %v", err)
			writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Create user
		// TODO: Consider adding avatar URL handling in the future
		user := &model.User{
			ID:           userID,
			Email:        email,
			FullName:     strings.TrimSpace(req.FullName),
			PasswordHash: &hashedPassword,
			AuthProvider: model.AuthProviderLocal,
			Timezone:     strings.TrimSpace(req.Timezone),
		}

		createdUserID, err := database.CreateUser(ctx, user)
		if err != nil {
			logger.Debugf("User creation failed: %v", err)
			writeErrorResponse(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		logger.Infof("User registered successfully: %s", email)

		// Return success response
		response := SignupResponse{
			Success: true,
			UserID:  createdUserID,
			Message: "User registered successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

// LoginHandler handles user authentication for local auth provider
func LoginHandler(database *db.PostgresDB, logger *logger.ServiceLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		// Parse request body
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Debugf("Invalid JSON in login request: %v", err)
			writeErrorResponse(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		// Validate required fields
		if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
			writeErrorResponse(w, "Email and password are required", http.StatusBadRequest)
			return
		}

		// Normalize email
		email := strings.TrimSpace(strings.ToLower(req.Email))

		// Get user from database
		user, err := database.GetUserByEmail(ctx, email)
		if err != nil {
			logger.Debugf("Database error getting user: %v", err)
			writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if user == nil {
			logger.Debugf("Login attempt for non-existent user: %s", email)
			writeErrorResponse(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// Check if user is local auth provider and has password
		if user.AuthProvider != model.AuthProviderLocal || user.PasswordHash == nil {
			logger.Debugf("Login attempt for non-local user: %s", email)
			writeErrorResponse(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// Verify password
		valid, err := auth.VerifyPassword(req.Password, *user.PasswordHash)
		if err != nil {
			logger.Debugf("Password verification error: %v", err)
			writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if !valid {
			logger.Debugf("Invalid password for user: %s", email)
			writeErrorResponse(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		logger.Infof("User logged in successfully: %s", email)

		// Return success response
		response := LoginResponse{
			Success: true,
			UserID:  user.ID,
			Message: "Login successful",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// Helper functions

func validateSignupRequest(req SignupRequest) error {
	if strings.TrimSpace(req.Email) == "" {
		return fmt.Errorf("email is required")
	}
	if strings.TrimSpace(req.Password) == "" {
		return fmt.Errorf("password is required")
	}
	if len(strings.TrimSpace(req.Password)) < 6 {
		return fmt.Errorf("password must be at least 6 characters long")
	}
	if strings.TrimSpace(req.FullName) == "" {
		return fmt.Errorf("full name is required")
	}
	if strings.TrimSpace(req.Timezone) == "" {
		return fmt.Errorf("timezone is required")
	}

	// Basic email validation
	if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

func writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]interface{}{
		"success": false,
		"message": message,
	}

	json.NewEncoder(w).Encode(response)
}
