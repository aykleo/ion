package sqlite

import "time"

type User struct {
	Username string
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
