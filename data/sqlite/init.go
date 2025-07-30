package sqlite

import (
	"database/sql"
	"path/filepath"

	"github.com/aykleo/ion/config"

	_ "github.com/glebarez/go-sqlite"
)

func InitSQLite() (*sql.DB, error) {
	path := config.GetConfigPath()
	dbPath := filepath.Join(path, "data.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := CreateAllTables(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func CreateUsersTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		username TEXT PRIMARY KEY
	)`
	_, err := db.Exec(query)
	return err
}

func CreateSecretsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS secrets (
		id TEXT PRIMARY KEY,
		salt TEXT NOT NULL,
		value TEXT NOT NULL,
		tags TEXT, -- JSON array stored as text
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	)`
	_, err := db.Exec(query)
	return err
}

func CreateAliasesTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS aliases (
		id TEXT PRIMARY KEY,
		value TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	)`
	_, err := db.Exec(query)
	return err
}

func CreateAllTables(db *sql.DB) error {
	if err := CreateUsersTable(db); err != nil {
		return err
	}
	if err := CreateSecretsTable(db); err != nil {
		return err
	}
	if err := CreateAliasesTable(db); err != nil {
		return err
	}
	return nil
}
