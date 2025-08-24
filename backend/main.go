package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	myAuth "github.com/anish-chanda/ferna/auth"
	"github.com/anish-chanda/ferna/db"
	"github.com/anish-chanda/ferna/db/sqlite3"
	"github.com/anish-chanda/ferna/handlers"
	"github.com/go-pkgz/auth/v2"
	"github.com/go-pkgz/auth/v2/avatar"
	"github.com/go-pkgz/auth/v2/provider"
	"github.com/go-pkgz/auth/v2/token"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// load env vars
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	baseUrl := os.Getenv("BASE_URL")
	if baseUrl == "" {
		log.Println("BASE_URL environment variable is not set, defaulting to http://localhost:8080")
		baseUrl = "http://localhost:8080"
	}

	// Prepare data directory
	dataDir := "./data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Initialize database
	dbPath := filepath.Join(dataDir, "ferna.db")
	dsn := fmt.Sprintf("file:%s?_foreign_keys=ON&parseTime=true", dbPath)

	// Connect & migrate
	fmt.Println("Migrating database...")
	var database db.Database = sqlite3.NewSQLiteDB()
	if err := database.Connect(dsn); err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer database.Close()

	if err := database.Migrate(); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	fmt.Println("Database connected and migrated successfully!")

	// setup auth options
	authOptions := auth.Opts{
		SecretReader: token.SecretFunc(func(id string) (string, error) { // secret key for JWT
			return jwtSecret, nil
		}),
		TokenDuration:  time.Minute * 5, // token expires in 5 minutes
		CookieDuration: time.Hour * 24,  // cookie expires in 1 day and will enforce re-login
		Issuer:         "ferna",
		URL:            baseUrl,
		DisableXSRF:    true,
		ClaimsUpd: token.ClaimsUpdFunc(func(cl token.Claims) token.Claims {
			if cl.User.Name == "" {
				return cl
			}
			u, err := database.GetUserByEmail(context.TODO(), cl.User.Name)
			if err != nil || u == nil {
				return cl
			}
			cl.User.SetStrAttr("uid", fmt.Sprint(u.ID))
			return cl
		}),
		AvatarStore: avatar.NewLocalFS("/tmp"),
	}

	// create auth service with providers
	authService := auth.NewService(authOptions)
	authService.AddDirectProvider("local", provider.CredCheckerFunc(func(user, password string) (ok bool, err error) {
		return myAuth.HandleLogin(database, user, password)
	}))

	r := mux.NewRouter()

	// register routes
	r.HandleFunc("/auth/local/signup", func(w http.ResponseWriter, r *http.Request) {
		myAuth.HandleSignup(database, w, r)
	}).Methods("POST")

	// setup auth routes
	authRoutes, avaRoutes := authService.Handlers()
	r.PathPrefix("/auth").Handler(authRoutes)
	r.PathPrefix("/avatar").Handler(avaRoutes)

	// create middleware and mount api endpoints
	authMiddleware := authService.Middleware()
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.Use(authMiddleware.Auth)

	// species routes
	apiRouter.HandleFunc("/species", handlers.SearchSpecies(database)).Methods("GET")

	// location routes
	apiRouter.HandleFunc("/locations", handlers.CreateLocation(database)).Methods("POST")
	apiRouter.HandleFunc("/locations", handlers.ListLocations(database)).Methods("GET")
	apiRouter.HandleFunc("/locations/{id}", handlers.GetLocation(database)).Methods("GET")
	apiRouter.HandleFunc("/locations/{id}", handlers.UpdateLocation(database)).Methods("PATCH")
	apiRouter.HandleFunc("/locations/{id}", handlers.DeleteLocation(database)).Methods("DELETE")

	// user plant routes
	apiRouter.HandleFunc("/plants", handlers.CreateUserPlant(database)).Methods("POST")
	apiRouter.HandleFunc("/plants", handlers.ListUserPlants(database)).Methods("GET")
	apiRouter.HandleFunc("/plants/{id}", handlers.GetUserPlant(database)).Methods("GET")
	apiRouter.HandleFunc("/plants/{id}", handlers.UpdateUserPlant(database)).Methods("PATCH")
	apiRouter.HandleFunc("/plants/{id}", handlers.DeleteUserPlant(database)).Methods("DELETE")

	// plant task routes
	apiRouter.HandleFunc("/plants/{plantId}/tasks", handlers.CreatePlantTask(database)).Methods("POST")
	apiRouter.HandleFunc("/plants/{plantId}/tasks", handlers.GetPlantTasksByPlantID(database)).Methods("GET")
	apiRouter.HandleFunc("/tasks/{id}", handlers.GetPlantTask(database)).Methods("GET")
	apiRouter.HandleFunc("/tasks/{id}", handlers.UpdatePlantTask(database)).Methods("PATCH")
	apiRouter.HandleFunc("/tasks/{id}", handlers.DeletePlantTask(database)).Methods("DELETE")
	apiRouter.HandleFunc("/tasks/overdue", handlers.GetOverdueTasks(database)).Methods("GET")

	// care event routes
	apiRouter.HandleFunc("/plants/{plantId}/events", handlers.CreateCareEvent(database)).Methods("POST")
	apiRouter.HandleFunc("/plants/{plantId}/events", handlers.GetCareEventsByPlantID(database)).Methods("GET")
	apiRouter.HandleFunc("/events/{id}", handlers.GetCareEvent(database)).Methods("GET")
	apiRouter.HandleFunc("/events/{id}", handlers.UpdateCareEvent(database)).Methods("PATCH")
	apiRouter.HandleFunc("/events/{id}", handlers.DeleteCareEvent(database)).Methods("DELETE")

	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", r)
}
