package data

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (s *Data) AddSecret(args []string, path string) error {
	if len(args) < 2 {
		return errors.New("invalid arguments, use ion secret add <name> <value> with optional -s <salt> and -t <tag1> <tag2>")
	}

	exists, _ := s.checkIfSecretExists(args[len(args)-2], path)
	if exists {
		return errors.New("secret already exists, use ion secret update <name> <new-value> to change it")
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

func (s *Data) UpdateSecretValue(args []string, path string) error {
	if len(args) != 2 {
		return errors.New("invalid arguments, use ion secret update <name> <new-value>")
	}
	var b strings.Builder
	secretName := args[len(args)-2]
	b.WriteString("secret ")
	b.WriteString(secretName)
	b.WriteString(" was not found")
	exists, index := s.checkIfSecretExists(secretName, path)
	if !exists {
		return errors.New(b.String())
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

func (s *Data) UpdateSecretName(args []string, path string) error {
	if len(args) != 2 {
		return errors.New("invalid arguments, use ion secret rename <name> <new-name>")
	}
	var b strings.Builder
	secretName := args[len(args)-2]
	b.WriteString("secret ")
	b.WriteString(secretName)
	b.WriteString(" was not found")
	exists, index := s.checkIfSecretExists(secretName, path)
	if !exists {
		return errors.New(b.String())
	}

	_, value, _, _, err := s.extractArgs(args, false)
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
	currentValue := data.Secrets[index].Value

	data.Secrets[index] = Secret{
		ID:        value,
		Salt:      currentSalt,
		Value:     currentValue,
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

func (s *Data) UpdateSecretTags(args []string, path string) error {
	if len(args) < 2 {
		return errors.New("invalid arguments, use ion secret tag <tag1> <tag2> <name>")
	}
	var b strings.Builder
	secretName := args[len(args)-1]
	b.WriteString("secret ")
	b.WriteString(secretName)
	b.WriteString(" was not found")
	exists, index := s.checkIfSecretExists(secretName, path)
	if !exists {
		return errors.New(b.String())
	}

	tagArgs := args[:len(args)-1]
	for i, tag := range tagArgs {
		tagArgs[i] = strings.ToUpper(tag)
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

	currentSalt := data.Secrets[index].Salt
	currentValue := data.Secrets[index].Value
	currentName := data.Secrets[index].ID

	data.Secrets[index] = Secret{
		ID:        currentName,
		Salt:      currentSalt,
		Value:     currentValue,
		Tags:      tagArgs,
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

func (s *Data) ListSecrets(args []string, path string) ([]Secret, bool, error) {
	if len(args) == 0 {
		return s.Secrets, false, nil
	}

	var b strings.Builder
	b.WriteString("ion secret list accepts two optionals flags, -d and -j\n\n")
	b.WriteString(" -d or --decrypt: decrypt the secrets showing their real values\n")
	b.WriteString(" -j or --json: show the secrets in json format\n\n")
	b.WriteString(" example: ion secret list -j -d\n")

	if len(args) > 2 {
		return nil, false, errors.New(b.String())
	}

	var hasDecrypt, hasJSON bool

	for _, arg := range args {
		switch arg {
		case "-d", "--decrypt":
			hasDecrypt = true
		case "-j", "--json":
			hasJSON = true
		default:
			return nil, false, errors.New(b.String())
		}
	}

	secrets := s.Secrets
	if hasDecrypt {
		decryptedSecrets := []Secret{}
		for _, secret := range s.Secrets {
			decryptedSecrets = append(decryptedSecrets, Secret{
				ID:        secret.ID,
				Salt:      secret.Salt,
				Value:     decrypt(secret.Salt, secret.Value),
				Tags:      secret.Tags,
				CreatedAt: secret.CreatedAt,
				UpdatedAt: secret.UpdatedAt,
			})
		}
		secrets = decryptedSecrets
	}

	return secrets, hasJSON, nil
}
