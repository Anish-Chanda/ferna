package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/anish-chanda/ferna/db"
	"github.com/anish-chanda/ferna/db/sqlite3"
	_ "github.com/mattn/go-sqlite3"
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

}
