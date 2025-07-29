package data

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

func (s *Data) AddSecret(args []string, path string) error {
	if len(args) < 2 {
		return errors.New("ion secret add <name> <value>")
	}

	exists, _ := s.checkIfSecretExists(args[len(args)-2], path)
	if exists {
		return errors.New("secret already exists, use ion secret update <name> <value> to change it")
	}

	name, value, salt, tgs, err := s.extractArgs(args, true)
	if err != nil {
		return err
	}

	dataPath := filepath.Join(path, "data.json")

	var data Data
	fileData, err := os.ReadFile(dataPath)
	if err != nil {
		data = *s
	} else {
		if err := json.Unmarshal(fileData, &data); err != nil {
			return err
		}
	}

	data.Secrets = append(data.Secrets, Secret{
		ID:        name,
		Salt:      salt,
		Value:     encrypt(salt, value),
		Tags:      tgs,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	s.Secrets = data.Secrets

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(dataPath, jsonData, 0644); err != nil {
		return err
	}
	return nil
}

func (s *Data) UpdateSecret(args []string, path string) error {
	if len(args) != 2 {
		return errors.New("ion secret update <name> <new-value>")
	}
	exists, index := s.checkIfSecretExists(args[len(args)-2], path)
	if !exists {
		return errors.New(args[len(args)-2] + " was not found")
	}

	name, value, _, _, err := s.extractArgs(args, false)
	if err != nil {
		return err
	}

	dataPath := filepath.Join(path, "data.json")

	var data Data
	fileData, err := os.ReadFile(dataPath)
	if err != nil {
		data = *s
	} else {
		if err := json.Unmarshal(fileData, &data); err != nil {
			return err
		}
	}

	currentTags := data.Secrets[index].Tags
	currentSalt := data.Secrets[index].Salt

	data.Secrets[index] = Secret{
		ID:        name,
		Salt:      currentSalt,
		Value:     encrypt(currentSalt, value),
		Tags:      currentTags,
		CreatedAt: data.Secrets[index].CreatedAt,
		UpdatedAt: time.Now(),
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(dataPath, jsonData, 0644); err != nil {
		return err
	}

	*s = data

	return nil
}
