package data

func InitData() Data {
	return Data{
		User: User{
			Username: "",
		},
		Passwords: []Password{},
		Aliases:   []Alias{},
	}
}
