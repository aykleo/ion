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

func (s *Data) UpdateAlias(args []string, path string) error {
	if s.db == nil {
		return errors.New("database connection not available")
	}

	if len(args) < 1 {
		return errors.New("invalid arguments, use ion alias update <alias-name>=<new-command>")
	}

	aliasName, command, err := s.extractAliasArgs(args)
	if err != nil {
		return err
	}

	exists, index := s.checkIfAliasExists(aliasName)
	if !exists {
		return errors.New("alias does not exist, use ion alias add <alias-name>=<command> to create it")
	}

	commandString := strings.Join(command, " ")

	alias := &s.Aliases[index]
	alias.Command = commandString
	alias.UpdatedAt = time.Now()

	sqliteAlias := sqlite.Alias{
		ID:        alias.ID,
		Name:      alias.Name,
		Command:   alias.Command,
		CreatedAt: alias.CreatedAt,
		UpdatedAt: alias.UpdatedAt,
	}

	if err := sqlite.UpdateAlias(s.db, sqliteAlias); err != nil {
		return err
	}

	return nil
}

func (s *Data) RenameAlias(args []string, path string) error {
	if s.db == nil {
		return errors.New("database connection not available")
	}

	if len(args) < 2 {
		return errors.New("invalid arguments, use ion alias rename <old-name> <new-name>")
	}

	oldName := args[0]
	newName := args[1]

	exists, index := s.checkIfAliasExists(oldName)
	if !exists {
		return errors.New("alias does not exist")
	}

	newExists, _ := s.checkIfAliasExists(newName)
	if newExists {
		return errors.New("alias with new name already exists")
	}

	alias := &s.Aliases[index]
	alias.Name = newName
	alias.UpdatedAt = time.Now()

	sqliteAlias := sqlite.Alias{
		ID:        alias.ID,
		Name:      alias.Name,
		Command:   alias.Command,
		CreatedAt: alias.CreatedAt,
		UpdatedAt: alias.UpdatedAt,
	}

	if err := sqlite.UpdateAlias(s.db, sqliteAlias); err != nil {
		return err
	}

	s.updateAliasIndex(oldName, newName, index)

	return nil
}

func (s *Data) RemoveAlias(args []string, path string) error {
	if s.db == nil {
		return errors.New("database connection not available")
	}

	if len(args) < 1 {
		return errors.New("invalid arguments, use ion alias remove <alias-name>")
	}

	aliasName := args[0]

	exists, index := s.checkIfAliasExists(aliasName)
	if !exists {
		return errors.New("alias does not exist")
	}

	alias := s.Aliases[index]

	if err := sqlite.RemoveAlias(s.db, alias.ID); err != nil {
		return err
	}

	s.Aliases = append(s.Aliases[:index], s.Aliases[index+1:]...)
	s.removeFromAliasIndex(aliasName, index)

	return nil
}

func (s *Data) ListAliases(args []string, path string) ([]Alias, bool, error) {
	if s.db == nil {
		return nil, false, errors.New("database connection not available")
	}

	isJson := false
	for _, arg := range args {
		if arg == "-j" {
			isJson = true
			break
		}
	}

	return s.Aliases, isJson, nil
}

func (s *Data) SearchAlias(args []string) ([]Alias, error) {
	if s.db == nil {
		return nil, errors.New("database connection not available")
	}

	if len(args) != 1 {
		return nil, errors.New("invalid arguments, use ion secret alias <name>")
	}

	if err := s.loadDataFromDB(); err != nil {
		return nil, err
	}

	aliasName := args[len(args)-1]

	index, err := s.fuzzySearchAlias(aliasName)
	if err != nil {
		return nil, err
	}

	return []Alias{s.Aliases[index]}, nil
}
