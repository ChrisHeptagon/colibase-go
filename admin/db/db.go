package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s.sqlite?cache=shared&mode=rwc&_journal_mode=WAL&_synchronous=NORMAL&_foreign_keys=ON", os.Getenv("DB_NAME")))
	if err != nil {
		fmt.Println("Error opening database:", err)
		return nil, err
	}
	db.Exec(fmt.Sprintf("CREATE DATABASE %s;", os.Getenv("DB_NAME")))
	return db, nil
}
