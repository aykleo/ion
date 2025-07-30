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
		INSERT INTO secrets (id, salt, value, tags, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, secret.ID, secret.Salt, secret.Value, string(tagsJSON), secret.CreatedAt, secret.UpdatedAt)

	return err
}

func GetSecrets(db *sql.DB) ([]Secret, error) {
	rows, err := db.Query(`
		SELECT id, salt, value, tags, created_at, updated_at
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
		err := rows.Scan(&secret.ID, &secret.Salt, &secret.Value, &tagsJSON, &secret.CreatedAt, &secret.UpdatedAt)
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

func RemoveSecret(db *sql.DB, secretID string) error {
	_, err := db.Exec("DELETE FROM secrets WHERE id = ?", secretID)
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

func UpdateSecretName(db *sql.DB, oldID, newID string) error {
	_, err := db.Exec(`
		UPDATE secrets 
		SET id = ?, updated_at = ?
		WHERE id = ?
	`, newID, time.Now(), oldID)
	return err
}

func UpdateSecretTags(db *sql.DB, secretID string, tags []string) error {
	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		UPDATE secrets 
		SET tags = ?, updated_at = ?
		WHERE id = ?
	`, string(tagsJSON), time.Now(), secretID)
	return err
}
