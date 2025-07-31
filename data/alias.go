package data

import (
	"errors"
	"strings"
	"time"

	"github.com/aykleo/ion/data/sqlite"
)

func (s *Data) AddAlias(args []string, path string) error {
	if s.db == nil {
		return errors.New("database connection not available")
	}

	if len(args) < 1 {
		return errors.New("invalid arguments, use ion alias add <alias-name>=<command>")
	}

	aliasName, command, err := s.extractAliasArgs(args)
	if err != nil {
		return err
	}

	exists, _ := s.checkIfAliasExists(aliasName)
	if exists {
		return errors.New("alias already exists, use ion alias update <alias-name> <new-command> to change it")
	}

	commandString := strings.Join(command, " ")

	sqliteAlias := sqlite.Alias{
		ID:        sqlite.GenerateUUID(),
		Name:      aliasName,
		Command:   commandString,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := sqlite.AddAlias(s.db, sqliteAlias); err != nil {
		return err
	}

	dataAlias := Alias{
		ID:        sqliteAlias.ID,
		Name:      sqliteAlias.Name,
		Command:   sqliteAlias.Command,
		CreatedAt: sqliteAlias.CreatedAt,
		UpdatedAt: sqliteAlias.UpdatedAt,
	}
	s.Aliases = append(s.Aliases, dataAlias)
	s.addToAliasIndex(aliasName, len(s.Aliases)-1)

	return nil
}
