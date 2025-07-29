package data

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func (s *Data) SetUsername(name, path string) {
	s.User.Username = name

	dataPath := filepath.Join(path, "data.json")
	jsonData, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile(dataPath, jsonData, 0644); err != nil {
		panic(err)
	}
}
