package data

func InitData() Data {
	data := Data{
		User: User{
			Username: "",
		},
		Secrets:     []Secret{},
		Aliases:     []Alias{},
		secretIndex: make(map[string]int),
	}
	data.buildSecretIndex()
	return data
}
