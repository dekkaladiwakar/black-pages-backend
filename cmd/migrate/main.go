package main

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Migration struct {
	Version int
	Name    string
	Path    string
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run cmd/migrate/main.go [up|down|status]")
	}

	command := os.Args[1]

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Database connection
	dsn := getEnv("DATABASE_URL", "postgresql://postgres:postgres@localhost:5432/black_pages?sslmode=disable")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Create migrations table if it doesn't exist
	if err := createMigrationsTable(db); err != nil {
		log.Fatal("Failed to create migrations table:", err)
	}

	switch command {
	case "up":
		if err := migrateUp(db); err != nil {
			log.Fatal("Migration up failed:", err)
		}
	case "down":
		if err := migrateDown(db); err != nil {
			log.Fatal("Migration down failed:", err)
		}
	case "status":
		if err := showStatus(db); err != nil {
			log.Fatal("Failed to show status:", err)
		}
	default:
		log.Fatal("Invalid command. Use: up, down, or status")
	}
}

func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := db.Exec(query)
	return err
}

func migrateUp(db *sql.DB) error {
	migrations, err := loadMigrations()
	if err != nil {
		return err
	}

	appliedVersions, err := getAppliedVersions(db)
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		if _, applied := appliedVersions[migration.Version]; applied {
			continue
		}

		fmt.Printf("Applying migration %d: %s\n", migration.Version, migration.Name)

		content, err := os.ReadFile(migration.Path)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %v", migration.Path, err)
		}

		// Execute migration
		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %d: %v", migration.Version, err)
		}

		// Record migration
		if _, err := db.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", migration.Version); err != nil {
			return fmt.Errorf("failed to record migration %d: %v", migration.Version, err)
		}

		fmt.Printf("✅ Applied migration %d successfully\n", migration.Version)
	}

	fmt.Println("All migrations applied successfully!")
	return nil
}

func migrateDown(db *sql.DB) error {
	// Get latest applied migration
	var latestVersion int
	err := db.QueryRow("SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1").Scan(&latestVersion)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No migrations to rollback")
			return nil
		}
		return err
	}

	fmt.Printf("Rolling back migration %d\n", latestVersion)

	// For now, we'll just remove the record (no down migrations implemented)
	if _, err := db.Exec("DELETE FROM schema_migrations WHERE version = $1", latestVersion); err != nil {
		return fmt.Errorf("failed to rollback migration %d: %v", latestVersion, err)
	}

	fmt.Printf("✅ Rolled back migration %d\n", latestVersion)
	return nil
}

func showStatus(db *sql.DB) error {
	migrations, err := loadMigrations()
	if err != nil {
		return err
	}

	appliedVersions, err := getAppliedVersions(db)
	if err != nil {
		return err
	}

	fmt.Println("Migration Status:")
	fmt.Println("================")

	for _, migration := range migrations {
		status := "❌ Pending"
		if _, applied := appliedVersions[migration.Version]; applied {
			status = "✅ Applied"
		}
		fmt.Printf("Version %d: %s - %s\n", migration.Version, migration.Name, status)
	}

	return nil
}

func loadMigrations() ([]Migration, error) {
	var migrations []Migration

	err := filepath.WalkDir("migrations", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".sql") {
			return nil
		}

		// Extract version from filename (e.g., "001_create_users_table.sql")
		filename := d.Name()
		parts := strings.SplitN(filename, "_", 2)
		if len(parts) < 2 {
			return fmt.Errorf("invalid migration filename format: %s", filename)
		}

		version, err := strconv.Atoi(parts[0])
		if err != nil {
			return fmt.Errorf("invalid version number in filename: %s", filename)
		}

		name := strings.TrimSuffix(parts[1], ".sql")
		name = strings.ReplaceAll(name, "_", " ")

		migrations = append(migrations, Migration{
			Version: version,
			Name:    name,
			Path:    path,
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

func getAppliedVersions(db *sql.DB) (map[int]bool, error) {
	versions := make(map[int]bool)

	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		versions[version] = true
	}

	return versions, rows.Err()
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}