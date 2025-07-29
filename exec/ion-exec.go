package exec

import (
	"errors"

	"github.com/aykleo/ion/config"
	"github.com/aykleo/ion/data"
	tea "github.com/charmbracelet/bubbletea"
)

type IonCommandHandler func(args []string, configPath string, dataRef data.IData) tea.Cmd

var routes = map[string]map[string]IonCommandHandler{
	"user": {
		"rename": changeUsername,
	},
	"secret": {
		"add":    addSecret,
		"update": updateSecretValue,
		"rename": updateSecretName,
		"tags":   updateSecretTags,
		"list":   listSecrets,
		"search": searchSecret,
		"remove": removeSecret,
		"use":    copySecretToClipboard,
	},
}

func ExecIonCommand(args []string, dataRef data.IData) tea.Cmd {
	configPath := config.GetConfigPath()

	if args[0] == "ionize" && len(args) == 1 {
		return func() tea.Msg {
			return CommandFinishedMsg{
				Command: "ionize",
				Output:  "ionize",
			}
		}
	}
	if len(args) < 3 {
		return func() tea.Msg {
			return CommandFinishedMsg{
				Err:     errors.New("ion [category] [action] <args>"),
				Command: "ion",
				Output:  "usage: ion [category] [action] <args>",
				NewDir:  currentDir,
			}
		}
	}

	category := args[1]
	action := args[2]
	args = args[3:]

	if categoryHandlers, categoryExists := routes[category]; categoryExists {
		if handler, actionExists := categoryHandlers[action]; actionExists {
			return handler(args, configPath, dataRef)
		}
	}

	return func() tea.Msg {
		return CommandFinishedMsg{
			Err:     errors.New("command not found, please use ion help"),
			Command: "ion",
			Output:  "command not found",
			NewDir:  currentDir,
		}
	}
}
