package sqlite

import "database/sql"

func GetUser(db *sql.DB) (User, error) {
	var username string
	err := db.QueryRow("SELECT username FROM users LIMIT 1").Scan(&username)
	if err == sql.ErrNoRows {
		return User{Username: ""}, nil
	}
	if err != nil {
		return User{}, err
	}
	return User{Username: username}, nil
}

func SetUser(db *sql.DB, username string) error {
	_, err := db.Exec(`INSERT OR REPLACE INTO users (username) VALUES (?)`, username)
	return err
}
