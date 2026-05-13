package database

import (
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dbURL string) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("failed to get working directory:", err)
	}

	migrationsPath := filepath.Join(cwd, "migrations")
	m, err := migrate.New(
		"file://"+migrationsPath,
		dbURL,
	)
	if err != nil {
		log.Fatal("migration init error:", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("migration up error:", err)
	}

	log.Println("migrations applied successfully")
}
