package database

import (
	"database/sql"
	"log"
)

var db *sql.DB

func New() *sql.DB{
    db, err := getConnection()

    if err != nil {
        log.Panic(err)
    }

    if err = db.Ping(); err != nil {
        log.Panic(err)
    }
    
    return db

}

func Close() error {

    err := db.Close()

    return err
}
