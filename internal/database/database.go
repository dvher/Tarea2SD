package database

import (
	"database/sql"
	"log"
)

func New() *sql.DB {
	db, err := getConnection()

	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}

	log.Println("Connected to database")

	return db

}
