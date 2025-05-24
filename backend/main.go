package main

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/anish-chanda/ferna/db"
	"github.com/anish-chanda/ferna/db/sqlite3"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/argon2"
)

const (
	argonTime    = 1
	argonMemory  = 64 * 1024
	argonThreads = 4
	argonSaltLen = 16
	argonKeyLen  = 32
)

func main() {
	// Prepare data directory
	dataDir := "./data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Initialize database
	dbPath := filepath.Join(dataDir, "ferna.db")
	dsn := fmt.Sprintf("file:%s?_foreign_keys=ON&parseTime=true", dbPath)

	// Connect & migrate
	var database db.Database = sqlite3.NewSQLiteDB()
	if err := database.Connect(dsn); err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer database.Close()

	if err := database.Migrate(); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	fmt.Println("Database connected and migrated successfully!")

	r := mux.NewRouter()

	// register routes
	r.HandleFunc("/auth/local/signup", func(w http.ResponseWriter, r *http.Request) {
		HandleSignup(database, w, r)
	})

	http.ListenAndServe(":8080", r)
}

func HandleSignup(db db.Database, w http.ResponseWriter, r *http.Request) {
	type SignupRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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
	//TODO: hash passoword
	passHash, err := hashPassword(password)
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
	//return success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"success":true, "user_id":%d}`, userID)
}

// hashPassword applies Argon2id with OWASP‐recommended params and returns
// a single string in the standard “$argon2id$v=19$m=…,t=…,p=…$salt$hash” format.
func hashPassword(password string) (string, error) {
	salt := make([]byte, argonSaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		argonTime,
		argonMemory,
		argonThreads,
		argonKeyLen,
	)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	parts := []string{
		"argon2id",
		fmt.Sprintf("v=%d", argon2.Version),
		fmt.Sprintf("m=%d,t=%d,p=%d", argonMemory, argonTime, argonThreads),
		b64Salt,
		b64Hash,
	}
	return "$" + strings.Join(parts, "$"), nil
}

// verifyPassword parses and verifies an encoded Argon2id hash.
func verifyPassword(password, encoded string) (bool, error) {
	// encoded: $argon2id$v=19$m=...,t=...,p=...$<salt>$<hash>
	fields := strings.Split(encoded, "$")
	if len(fields) != 6 || fields[1] != "argon2id" {
		return false, fmt.Errorf("invalid hash format")
	}
	var memory, timeParam, threads uint32
	if _, err := fmt.Sscanf(fields[3], "m=%d,t=%d,p=%d", &memory, &timeParam, &threads); err != nil {
		return false, err
	}
	salt, err := base64.RawStdEncoding.DecodeString(fields[4])
	if err != nil {
		return false, err
	}
	hash, err := base64.RawStdEncoding.DecodeString(fields[5])
	if err != nil {
		return false, err
	}

	computed := argon2.IDKey([]byte(password), salt, timeParam, memory, uint8(threads), uint32(len(hash)))
	// constant-time compare
	if subtle.ConstantTimeCompare(computed, hash) == 1 {
		return true, nil
	}
	return false, nil
}
