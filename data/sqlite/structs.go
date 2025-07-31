package sqlite

import "time"

type User struct {
	Username string
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
