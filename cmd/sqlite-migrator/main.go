package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var (
		storagePath, migrationsPath, migrationsTable string
		isDownMigration                              bool
	)

	flag.StringVar(&storagePath, "storage-path", "", "Path to the storage database file")
	flag.StringVar(&migrationsPath, "migrations-path", "", "Path to the migrations directory")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "Name of the migrations table")
	flag.BoolVar(&isDownMigration, "down", false, "Set to true to apply down migrations instead of up migrations")
	flag.Parse()

	if storagePath == "" {
		panic("storage-path is required")
	}

	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable),
	)
	if err != nil {
		panic(err)
	}

	if isDownMigration {
		if err := m.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no migrations to revert")
			} else {
				panic(err)
			}
		}
		fmt.Println("migrations reverted successfully")
		return
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
		} else {
			panic(err)
		}
	}


	fmt.Println("migrations applied successfully")
}
