package storage

func InitStorage() Storage {
	return Storage{
		User: User{
			Username: "",
		},
		Passwords: []Password{},
		Aliases:   []Alias{},
	}
}
