package data

import (
	"errors"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/aykleo/ion/data/sqlite"
)

func (s *Data) AddSecret(args []string, path string) error {
	if s.db == nil {
		return errors.New("database connection not available")
	}

	if len(args) < 2 {
		return errors.New("invalid arguments, use ion secret add <name> <value> with optional -s <salt> and -t <tag1> <tag2>")
	}

	exists, _ := s.checkIfSecretExists(args[len(args)-2])
	if exists {
		return errors.New("secret already exists, use ion secret update <name> <new-value> to change it")
	}

	name, value, salt, tgs, err := s.extractSecretArgs(args, true)
	if err != nil {
		return err
	}

	sqliteSecret := sqlite.Secret{
		ID:        sqlite.GenerateUUID(),
		Name:      name,
		Salt:      salt,
		Value:     encrypt(salt, value),
		Tags:      tgs,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := sqlite.AddSecret(s.db, sqliteSecret); err != nil {
		return err
	}

	dataSecret := Secret{
		ID:        sqliteSecret.ID,
		Name:      sqliteSecret.Name,
		Salt:      sqliteSecret.Salt,
		Value:     sqliteSecret.Value,
		Tags:      sqliteSecret.Tags,
		CreatedAt: sqliteSecret.CreatedAt,
		UpdatedAt: sqliteSecret.UpdatedAt,
	}
	s.Secrets = append(s.Secrets, dataSecret)
	s.addToSecretIndex(name, len(s.Secrets)-1)

	return nil
}

func (s *Data) UpdateSecretValue(args []string, path string) error {
	if s.db == nil {
		return errors.New("database connection not available")
	}

	if len(args) != 2 {
		return errors.New("invalid arguments, use ion secret update <name> <new-value>")
	}

	var b strings.Builder
	secretName := args[len(args)-2]
	b.WriteString("secret ")
	b.WriteString(secretName)
	b.WriteString(" was not found")
	exists, index := s.checkIfSecretExists(secretName)
	if !exists {
		return errors.New(b.String())
	}

	name, value, _, _, err := s.extractSecretArgs(args, false)
	if err != nil {
		return err
	}

	currentSalt := s.Secrets[index].Salt
	encryptedValue := encrypt(currentSalt, value)

	if err := sqlite.UpdateSecretValueByName(s.db, name, encryptedValue); err != nil {
		return err
	}

	s.ensureSecretIndex()
	if index, exists := s.secretIndex[name]; exists {
		s.Secrets[index].Value = encryptedValue
		s.Secrets[index].UpdatedAt = time.Now()
	}

	return nil
}

func (s *Data) UpdateSecretName(args []string, path string) error {
	if s.db == nil {
		return errors.New("database connection not available")
	}

	if len(args) != 2 {
		return errors.New("invalid arguments, use ion secret rename <name> <new-name>")
	}

	var b strings.Builder
	secretName := args[len(args)-2]
	b.WriteString("secret ")
	b.WriteString(secretName)
	b.WriteString(" was not found")
	exists, _ := s.checkIfSecretExists(secretName)
	if !exists {
		return errors.New(b.String())
	}

	_, newName, _, _, err := s.extractSecretArgs(args, false)
	if err != nil {
		return err
	}

	if err := sqlite.UpdateSecretName(s.db, secretName, newName); err != nil {
		return err
	}

	s.ensureSecretIndex()
	if index, exists := s.secretIndex[secretName]; exists {
		s.Secrets[index].Name = newName
		s.Secrets[index].UpdatedAt = time.Now()
		s.updateSecretIndex(secretName, newName, index)
	}

	return nil
}

func (s *Data) UpdateSecretTags(args []string, path string) error {
	if s.db == nil {
		return errors.New("database connection not available")
	}

	if len(args) < 2 {
		return errors.New("invalid arguments, use ion secret tag <tag1> <tag2> <name>")
	}

	var b strings.Builder
	secretName := args[len(args)-1]
	b.WriteString("secret ")
	b.WriteString(secretName)
	b.WriteString(" was not found")
	exists, _ := s.checkIfSecretExists(secretName)
	if !exists {
		return errors.New(b.String())
	}

	tagArgs := args[:len(args)-1]
	for i, tag := range tagArgs {
		tagArgs[i] = strings.ToUpper(tag)
	}

	if err := sqlite.UpdateSecretTags(s.db, secretName, tagArgs); err != nil {
		return err
	}

	s.ensureSecretIndex()
	if index, exists := s.secretIndex[secretName]; exists {
		s.Secrets[index].Tags = tagArgs
		s.Secrets[index].UpdatedAt = time.Now()
	}

	return nil
}

func (s *Data) ListSecrets(args []string, path string) ([]Secret, bool, error) {
	if s.db == nil {
		return nil, false, errors.New("database connection not available")
	}

	if err := s.loadDataFromDB(); err != nil {
		return nil, false, err
	}

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
				Name:      secret.Name,
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

func (s *Data) SearchSecret(args []string) ([]Secret, error) {
	if s.db == nil {
		return nil, errors.New("database connection not available")
	}

	if len(args) != 1 {
		return nil, errors.New("invalid arguments, use ion secret search <name>")
	}

	if err := s.loadDataFromDB(); err != nil {
		return nil, err
	}

	secretName := args[len(args)-1]

	index, err := s.fuzzySearchSecret(secretName)
	if err != nil {
		return nil, err
	}

	return []Secret{s.Secrets[index]}, nil
}

func (s *Data) RemoveSecret(args []string, path string) error {
	if s.db == nil {
		return errors.New("database connection not available")
	}

	if len(args) != 1 {
		return errors.New("invalid arguments, use ion secret remove <name>")
	}

	secretName := args[len(args)-1]
	exists, _ := s.checkIfSecretExists(secretName)
	if !exists {
		return errors.New("secret was not found")
	}

	if err := sqlite.RemoveSecretByName(s.db, secretName); err != nil {
		return err
	}

	s.ensureSecretIndex()
	if index, exists := s.secretIndex[secretName]; exists {
		s.Secrets = append(s.Secrets[:index], s.Secrets[index+1:]...)
		s.removeFromSecretIndex(secretName, index)
	}

	return nil
}

func (s *Data) CopySecretToClipboard(args []string, path string) error {
	if s.db == nil {
		return errors.New("database connection not available")
	}

	if len(args) != 1 {
		return errors.New("invalid arguments, use ion secret copy <name>")
	}

	if err := s.loadDataFromDB(); err != nil {
		return err
	}

	secretName := args[len(args)-1]
	exists, index := s.checkIfSecretExists(secretName)
	if !exists {
		return errors.New("secret was not found")
	}

	secret := s.Secrets[index]
	decryptedValue := decrypt(secret.Salt, secret.Value)
	clipboard.WriteAll(decryptedValue)

	return nil
}
