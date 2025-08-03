package exec

import (
	"errors"

	"github.com/aykleo/ion/config"
	"github.com/aykleo/ion/data"
	tea "github.com/charmbracelet/bubbletea"
)

func ExecIonCommand(args []string, dataRef data.IData) tea.Cmd {
	configPath := config.GetConfigPath()

	if len(args) > 0 && args[0] == "ionize" && len(args) == 1 {
		return func() tea.Msg {
			return CommandFinishedMsg{
				Command: "ionize",
				Output:  "ionize",
			}
		}
	}

	if len(args) == 2 {
		action := args[1]

		if handler, actionExists := simpleCmds[action]; actionExists {
			return handler(args)
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
			Err:     errors.New("command not found"),
			Command: "ion",
			Output:  "command not found, please use ion help",
			NewDir:  currentDir,
		}
	}
}
