package data

import (
	"time"
)

type Data struct {
	User        User
	Secrets     []Secret
	Aliases     []Alias
	secretIndex map[string]int
}

type Secret struct {
	ID        string
	Salt      string
	Value     string
	Tags      []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Alias struct {
	ID        string
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	Username string
}
