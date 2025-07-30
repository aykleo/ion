package data

import (
	"github.com/aykleo/ion/data/sqlite"
)

func (s *Data) SetUsername(name, path string) {
	if s.db == nil {
		panic("database connection not available")
	}

	if err := sqlite.SetUser(s.db, name); err != nil {
		panic(err)
	}

	s.User.Username = name
}

func (s *Data) GetUser() User {
	return s.User
}
