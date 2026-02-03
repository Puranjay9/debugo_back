package db

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

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
	if err != nil {
		return err
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS history (
			id UUID PRIMARY KEY,
			project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
			error TEXT,
			error_source TEXT,
			code_diff TEXT,
			timestamp TIMESTAMPTZ DEFAULT NOW()
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
