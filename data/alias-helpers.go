package data

import (
	"errors"
	"strings"
)

func (s *Data) buildAliasIndex() {
	if s.aliasIndex == nil {
		s.aliasIndex = make(map[string]int)
	} else {
		for k := range s.aliasIndex {
			delete(s.aliasIndex, k)
		}
	}

	for i, alias := range s.Aliases {
		s.aliasIndex[alias.Name] = i
	}
}

func (s *Data) ensureAliasIndex() {
	if s.aliasIndex == nil || len(s.aliasIndex) != len(s.Aliases) {
		s.buildAliasIndex()
	}
}

func (s *Data) addToAliasIndex(aliasName string, index int) {
	if s.aliasIndex == nil {
		s.aliasIndex = make(map[string]int)
	}
	s.aliasIndex[aliasName] = index
}

func (s *Data) removeFromAliasIndex(aliasName string, removedIndex int) {
	if s.aliasIndex == nil {
		return
	}

	delete(s.aliasIndex, aliasName)

	for name, idx := range s.aliasIndex {
		if idx > removedIndex {
			s.aliasIndex[name] = idx - 1
		}
	}
}

func (s *Data) updateAliasIndex(oldName, newName string, index int) {
	if s.aliasIndex == nil {
		s.aliasIndex = make(map[string]int)
	}

	delete(s.aliasIndex, oldName)
	s.aliasIndex[newName] = index
}

func (s *Data) checkIfAliasExists(name string) (bool, int) {
	s.ensureAliasIndex()

	if index, exists := s.aliasIndex[name]; exists {
		return true, index
	}
	return false, -1
}

func (s *Data) extractAliasArgs(args []string) (string, []string, error) {
	separator := "="

	for i, arg := range args {
		if strings.Contains(arg, separator) {
			parts := strings.SplitN(arg, separator, 2)
			if len(parts) == 2 && parts[0] != "" {
				aliasName := strings.TrimSpace(parts[0])
				command := strings.TrimSpace(parts[1])
				remainingArgs := args[i+1:]

				if command == "" && len(remainingArgs) == 0 {
					return "", nil, errors.New("invalid arguments, command cannot be empty")
				}

				if command == "" && len(remainingArgs) > 0 {
					return aliasName, remainingArgs, nil
				}

				allCommandParts := append([]string{command}, remainingArgs...)
				return aliasName, allCommandParts, nil
			}
		}
	}

	for i := 0; i < len(args)-1; i++ {
		if i+1 < len(args) && args[i+1] == separator {
			aliasName := strings.TrimSpace(args[i])
			if aliasName != "" && i+2 < len(args) {
				remainingArgs := args[i+2:]
				return aliasName, remainingArgs, nil
			}
		}
	}

	return "", nil, errors.New("invalid arguments, use ion alias add <alias-name>=<command> with a '=' in between")
}
