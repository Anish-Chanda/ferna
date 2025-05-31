package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/anish-chanda/ferna/db"
)

func HandleLogin(db db.Database, user, password string) (ok bool, err error) {
	// set timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//get user by email
	normalized := strings.TrimSpace(strings.ToLower(user))
	userRecord, err := db.GetUserByEmail(ctx, normalized)
	if err != nil {
		fmt.Println("failed to get user by email:", err)
		return false, err
	}
	if userRecord == nil {
		// No such user
		return false, nil
	}

	// verify password
	ok, err = VerifyPassword(password, userRecord.PassHash)
	if err != nil {
		fmt.Println("failed to verify password:", err)
		return false, fmt.Errorf("failed to verify password: %w", err)
	}

	if !ok {
		return false, fmt.Errorf("invalid credentials")
	}
	// if password is correct, return true
	return ok, err
}

func HandleSignup(db db.Database, w http.ResponseWriter, r *http.Request) {
	type SignupRequest struct {
		Email    string `json:"user"`
		Password string `json:"passwd"`
	}

	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	email := req.Email
	password := req.Password

	if email == "" || password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	email = strings.TrimSpace(strings.ToLower(email))

	//set timeout ctx
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	//check if user with email already exists
	exists, err := db.CheckIfEmailExists(ctx, email)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("failed to check email existence: %v", err)
		return
	}
	if exists {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}

	passHash, err := HashPassword(password)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("failed to hash password: %v", err)
		return
	}

	//create user in database
	userID, err := db.CreateUser(ctx, email, passHash)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("failed to create user: %v", err)
		return
	}
	fmt.Println("User created with ID:", userID)

	//return success
	payload := struct {
		Success bool  `json:"success"`
		UserID  int64 `json:"user_id"`
	}{
		Success: true,
		UserID:  userID,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to write JSON response: %v", err)
	}
}
