package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func getConnection() (db *sql.DB, err error) {

	DB_USER := os.Getenv("DB_USER")
	DB_PASS := os.Getenv("DB_PASS")
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_NAME := os.Getenv("DB_NAME")
	uri := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_NAME)

	db, err = sql.Open("postgres", uri)

	return
}
