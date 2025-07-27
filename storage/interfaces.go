package storage

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
