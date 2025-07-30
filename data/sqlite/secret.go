package sqlite

import (
	"database/sql"
	"encoding/json"
	"time"
)

func AddSecret(db *sql.DB, secret Secret) error {
	tagsJSON, err := json.Marshal(secret.Tags)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO secrets (id, name, salt, value, tags, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, secret.ID, secret.Name, secret.Salt, secret.Value, string(tagsJSON), secret.CreatedAt, secret.UpdatedAt)

	return err
}

func GetSecrets(db *sql.DB) ([]Secret, error) {
	rows, err := db.Query(`
		SELECT id, name, salt, value, tags, created_at, updated_at
		FROM secrets
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var secrets []Secret
	for rows.Next() {
		var secret Secret
		var tagsJSON string
		err := rows.Scan(&secret.ID, &secret.Name, &secret.Salt, &secret.Value, &tagsJSON, &secret.CreatedAt, &secret.UpdatedAt)
		if err != nil {
			return nil, err
		}

		if tagsJSON != "" {
			err = json.Unmarshal([]byte(tagsJSON), &secret.Tags)
			if err != nil {
				return nil, err
			}
		}

		secrets = append(secrets, secret)
	}

	return secrets, nil
}

// GetSecretByName finds a secret by its user-friendly name
func GetSecretByName(db *sql.DB, name string) (Secret, error) {
	var secret Secret
	var tagsJSON string

	err := db.QueryRow(`
		SELECT id, name, salt, value, tags, created_at, updated_at
		FROM secrets 
		WHERE name = ?
	`, name).Scan(&secret.ID, &secret.Name, &secret.Salt, &secret.Value, &tagsJSON, &secret.CreatedAt, &secret.UpdatedAt)

	if err != nil {
		return Secret{}, err
	}

	if tagsJSON != "" {
		err = json.Unmarshal([]byte(tagsJSON), &secret.Tags)
		if err != nil {
			return Secret{}, err
		}
	}

	return secret, nil
}

func RemoveSecret(db *sql.DB, secretID string) error {
	_, err := db.Exec("DELETE FROM secrets WHERE id = ?", secretID)
	return err
}

func RemoveSecretByName(db *sql.DB, name string) error {
	_, err := db.Exec("DELETE FROM secrets WHERE name = ?", name)
	return err
}

func UpdateSecretValue(db *sql.DB, secretID, newValue string) error {
	_, err := db.Exec(`
		UPDATE secrets 
		SET value = ?, updated_at = ?
		WHERE id = ?
	`, newValue, time.Now(), secretID)
	return err
}

func UpdateSecretValueByName(db *sql.DB, name, newValue string) error {
	_, err := db.Exec(`
		UPDATE secrets 
		SET value = ?, updated_at = ?
		WHERE name = ?
	`, newValue, time.Now(), name)
	return err
}

func UpdateSecretName(db *sql.DB, oldName, newName string) error {
	_, err := db.Exec(`
		UPDATE secrets 
		SET name = ?, updated_at = ?
		WHERE name = ?
	`, newName, time.Now(), oldName)
	return err
}

func UpdateSecretTags(db *sql.DB, name string, tags []string) error {
	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		UPDATE secrets 
		SET tags = ?, updated_at = ?
		WHERE name = ?
	`, string(tagsJSON), time.Now(), name)
	return err
}
