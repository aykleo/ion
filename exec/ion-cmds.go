package exec

import (
	"errors"
	"strings"

	"github.com/aykleo/ion/data"
	tea "github.com/charmbracelet/bubbletea"
)

func changeUsername(args []string, configPath string, dataRef data.IData) tea.Cmd {
	if len(args) != 1 {
		return func() tea.Msg {
			return CommandFinishedMsg{
				Err:     errors.New("ion user <username>"),
				Command: strings.Join(args, " "),
				Output:  "ion user needs one argument, try [ ion user <username> ]",
				NewDir:  currentDir,
			}
		}
	}
	username := args[0]
	dataRef.SetUsername(username, configPath)
	return func() tea.Msg {
		return CommandFinishedMsg{
			Command: "ion set-username " + username,
			Output:  "username set to " + username,
			NewDir:  currentDir,
		}
	}
}
