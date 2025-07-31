package sqlite

import "database/sql"

func AddAlias(db *sql.DB, alias Alias) error {
	_, err := db.Exec(`
		INSERT INTO aliases (id, name, command, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`, alias.ID, alias.Name, alias.Command, alias.CreatedAt, alias.UpdatedAt)
	return err
}

func GetAliases(db *sql.DB) ([]Alias, error) {
	rows, err := db.Query(`
		SELECT id, name, command, created_at, updated_at
		FROM aliases
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var aliases []Alias
	for rows.Next() {
		var alias Alias
		err := rows.Scan(&alias.ID, &alias.Name, &alias.Command, &alias.CreatedAt, &alias.UpdatedAt)
		if err != nil {
			return nil, err
		}
		aliases = append(aliases, alias)
	}

	return aliases, nil
}

func UpdateAlias(db *sql.DB, alias Alias) error {
	_, err := db.Exec(`
		UPDATE aliases 
		SET name = ?, command = ?, updated_at = ?
		WHERE id = ?
	`, alias.Name, alias.Command, alias.UpdatedAt, alias.ID)
	return err
}

func RemoveAlias(db *sql.DB, id string) error {
	_, err := db.Exec(`
		DELETE FROM aliases WHERE id = ?
	`, id)
	return err
}
