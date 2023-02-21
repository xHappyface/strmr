package database

import (
	"errors"
	"os"

	"github.com/jmoiron/sqlx"
)

type Database struct {
	db *sqlx.DB
}

func GetDB(name string) (*sqlx.DB, error) {
	db_name := "./" + name + ".db"
	if _, err := os.Stat(db_name); err != nil {
		msg := "Could not get database file: " + err.Error()
		return nil, errors.New(msg)
	}
	database, err := sqlx.Open("sqlite3", db_name+"?_foreign_keys=on")
	if err != nil {
		msg := "Could not open database: " + err.Error()
		return nil, errors.New(msg)
	}
	err = database.Ping()
	if err != nil {
		msg := "Could not ping database: " + err.Error()
		return nil, errors.New(msg)
	}
	return database, nil
}

func New(database *sqlx.DB) *Database {
	return &Database{
		db: database,
	}
}
