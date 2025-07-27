package storage

import (
	"time"
)

type Storage struct {
	User      User
	Passwords []Password
	Aliases   []Alias
}

func InitStorage() Storage {
	return Storage{
		User: User{
			Username: "aykleo",
		},
		Passwords: []Password{},
		Aliases:   []Alias{},
	}
}

type Password struct {
	ID        string
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

type IStorage interface {
	GetUser() User
	GetPasswords() []Password
	GetAliases() []Alias
}

func (s *Storage) GetUser() User {
	return s.User
}

func (s *Storage) GetPasswords() []Password {
	return s.Passwords
}

func (s *Storage) GetAliases() []Alias {
	return s.Aliases
}

func NewStorage() IStorage {
	storage := InitStorage()
	return &storage
}
