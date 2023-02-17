package app

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"os"
)

var (
	db *sql.DB
)

func connectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
}

// InitDatabase For Vanilla SQL
func InitDatabase() (*sql.DB, error) {
	localDB, err := sql.Open(os.Getenv("DB_DRIVER"), connectionString())
	if err != nil {
		return nil, err
	}
	log.Info("Database connection was created")

	db = localDB

	return localDB, nil
}

func RunMigrations(rootDir ...string) error {

	basePath := "."
	if len(rootDir) != 0 {
		basePath = rootDir[0]
	}

	err := goose.Up(db, basePath+"/migrations")
	if err != nil {
		return err
	}
	return nil
}
