package utils

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

var MigrationFiles embed.FS

// InitializeDB creates the database file and runs migrations if needed.
func InitializeDB() (*sql.DB, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	aiDir := filepath.Join(homeDir, ".ai")
	dbPath := filepath.Join(aiDir, "history.db")

	// Check if the .ai directory exists; create it if it doesn't
	if _, err := os.Stat(aiDir); os.IsNotExist(err) {
		if err := os.MkdirAll(aiDir, 0755); err != nil { // Create with appropriate permissions
			return nil, fmt.Errorf("failed to create .ai directory: %w", err)
		}
	}

	// Check if the database file exists
	dbExists := true
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		dbExists = false
	}

	// Open the database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// If the database didn't exist, run migrations
	if !dbExists {
		if err := runMigrations(db); err != nil {
			db.Close() // Close the connection on error
			return nil, fmt.Errorf("failed to run migrations: %w", err)
		}
	}

	return db, nil
}

// runMigrations applies the embedded SQL migrations to the database.
func runMigrations(db *sql.DB) error {
	entries, err := MigrationFiles.ReadDir("db/migration")
	if err != nil {
		return fmt.Errorf("failed to read migration directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".sql" {
			filePath := filepath.Join("db/migration", entry.Name())
			migrationBytes, err := MigrationFiles.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read migration file %s: %w", filePath, err)
			}

			// Execute the migration SQL
			_, err = db.Exec(string(migrationBytes))
			if err != nil {
				return fmt.Errorf("failed to execute migration %s: %w", filePath, err)
			}
			fmt.Printf("Applied migration: %s\n", filePath)
		}
	}
	return nil
}
