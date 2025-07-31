package data

import (
	"database/sql"
	"time"
)

type Data struct {
	User        User
	Secrets     []Secret
	Aliases     []Alias
	secretIndex map[string]int
	aliasIndex  map[string]int
	db          *sql.DB
}

type Secret struct {
	ID        string
	Name      string
	Salt      string
	Value     string
	Tags      []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Alias struct {
	ID        string
	Name      string
	Command   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	Username string
}
