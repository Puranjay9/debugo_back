package db

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

type Project struct {
	ID   int
	Name string
}

var DB *sql.DB

func ConnectToDB() error {
	var err error
	DB, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS projects (
			id UUID PRIMARY KEY, 
			name TEXT
		)
	`)

	return err
}

func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
