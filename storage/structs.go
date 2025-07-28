package storage

import (
	"time"
)

type Storage struct {
	User      User
	Passwords []Password
	Aliases   []Alias
}

type Password struct {
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
