package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type IStorage interface {
	GetUser() User
	GetOrCreateStorageFields(path string) (IStorage, bool)
	GetStorage() IStorage
}

func (s *Storage) GetUser() User {
	return s.User
}

func NewStorage() IStorage {
	storage := InitStorage()
	return &storage
}
func (s *Storage) GetStorage() IStorage {
	return s
}

func (s *Storage) GetOrCreateStorageFields(path string) (IStorage, bool) {
	dataPath := filepath.Join(path, "data.json")

	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		emptyStorage := s.GetStorage()

		if err := os.MkdirAll(filepath.Dir(dataPath), 0755); err != nil {
			return nil, false
		}

		jsonData, err := json.MarshalIndent(emptyStorage, "", "  ")
		if err != nil {
			return nil, false
		}

		if err := os.WriteFile(dataPath, jsonData, 0644); err != nil {
			return nil, false
		}

		return emptyStorage, false
	}

	fileData, err := os.ReadFile(dataPath)
	if err != nil {
		return nil, false
	}

	var storage Storage
	if err := json.Unmarshal(fileData, &storage); err != nil {
		return nil, false
	}

	return &storage, true
}
