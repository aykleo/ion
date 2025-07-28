package data

func InitData() Data {
	return Data{
		User: User{
			Username: "",
		},
		Secrets: []Secret{},
		Aliases: []Alias{},
	}
}
