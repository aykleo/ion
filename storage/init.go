package storage

func InitStorage() Storage {
	return Storage{
		User: User{
			Username: "aykleo",
		},
		Passwords: []Password{},
		Aliases:   []Alias{},
	}
}
