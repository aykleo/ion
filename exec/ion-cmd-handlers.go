package exec

import (
	"github.com/aykleo/ion/data"
	tea "github.com/charmbracelet/bubbletea"
)

type IonCompositeCommandHandler func(args []string, configPath string, dataRef data.IData) tea.Cmd

var routes = map[string]map[string]IonCompositeCommandHandler{
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
	"alias": {
		"add":    addAlias,
		"update": updateAlias,
		"rename": renameAlias,
		"remove": removeAlias,
		"list":   listAliases,
		"search": searchAliases,
	},
}

type IonSimpleCommandHandler func(args []string) tea.Cmd

var simpleCmds = map[string]IonSimpleCommandHandler{
	"zen": toggleZenMode,
}
