package migrations

import (
	"embed"
	"fmt"
)

//go:embed sqlite3/*sql
var SQLiteMigrations embed.FS

func GetMigrationsFS(dbType string) (embed.FS, string, error) {
	switch dbType {
	case "sqlite3":
		return SQLiteMigrations, "sqlite3", nil
	default:
		return embed.FS{}, "", fmt.Errorf("unsupported database type: %s", dbType)
	}
	// add case for postgres later

}
