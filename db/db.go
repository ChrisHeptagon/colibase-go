package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {

	db, err := sql.Open("sqlite3", "./"+os.Getenv("DB_NAME")+".db")

	if err != nil {
		return nil, err
	}
	db.Exec(fmt.Sprintf("CREATE DATABASE %s;", os.Getenv("DB_NAME")))

	return db, nil
}
