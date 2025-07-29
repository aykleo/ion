package data

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type IData interface {
	GetUser() User
	GetData() IData

	GetOrCreateDataFields(path string) (IData, bool)

	SetUsername(name, path string)

	AddSecret(args []string, path string) error
	UpdateSecretValue(args []string, path string) error
	UpdateSecretName(args []string, path string) error
	UpdateSecretTags(args []string, path string) error
	ListSecrets(args []string, path string) ([]Secret, bool, error)
	SearchSecret(args []string) ([]Secret, error)
	RemoveSecret(args []string, path string) error
}

func (s *Data) GetUser() User {
	return s.User
}

func NewData() IData {
	data := InitData()
	return &data
}
func (s *Data) GetData() IData {
	return s
}

func (s *Data) GetOrCreateDataFields(path string) (IData, bool) {
	dataPath := filepath.Join(path, "data.json")

	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		emptyData := s.GetData()

		if err := os.MkdirAll(filepath.Dir(dataPath), 0755); err != nil {
			return nil, false
		}

		jsonData, err := json.MarshalIndent(emptyData, "", "  ")
		if err != nil {
			return nil, false
		}

		if err := os.WriteFile(dataPath, jsonData, 0644); err != nil {
			return nil, false
		}

		return emptyData, false
	}

	fileData, err := os.ReadFile(dataPath)
	if err != nil {
		return nil, false
	}

	var data Data
	if err := json.Unmarshal(fileData, &data); err != nil {
		return nil, false
	}

	return &data, true
}
