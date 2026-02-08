package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"sort"
)

func RunMigrations(db *sql.DB) error {
	// Try multiple possible paths for migrations directory
	possiblePaths := []string{
		"migrations",
		"./migrations",
		"../migrations",
	}

	var migrationDir string
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			migrationDir = path
			break
		}
	}

	if migrationDir == "" {
		return filepath.ErrBadPattern
	}

	files, err := filepath.Glob(filepath.Join(migrationDir, "*.sql"))
	if err != nil {
		return err
	}

	if len(files) == 0 {
		log.Println("No migration files found in", migrationDir)
		return nil
	}

	sort.Strings(files)

	log.Println("Found", len(files), "migration file(s)")

	for _, file := range files {
		log.Println("Applying migration:", filepath.Base(file))

		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			log.Printf("Error reading migration file %s: %v", file, err)
			return err
		}

		_, err = db.Exec(string(sqlBytes))
		if err != nil {
			log.Printf("Error executing migration %s: %v", file, err)
			return err
		}

		log.Println("✓ Successfully applied:", filepath.Base(file))
	}

	log.Println("All migrations completed successfully!")
	return nil
}
