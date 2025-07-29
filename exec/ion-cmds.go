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
				Err:     errors.New("ion user set <username>"),
				Command: strings.Join(args, " "),
				Output:  "ion user needs one argument, try [ ion user set <username> ]",
				NewDir:  currentDir,
			}
		}
	}
	username := args[0]
	dataRef.SetUsername(username, configPath)
	return func() tea.Msg {
		return CommandFinishedMsg{
			Command: "ion user set " + username,
			Output:  "username set to " + username,
			NewDir:  currentDir,
		}
	}
}

func addSecret(args []string, configPath string, dataRef data.IData) tea.Cmd {
	if len(args) < 2 {
		return func() tea.Msg {
			var b strings.Builder
			b.WriteString("ion secret add accepts two optionals flags and two arguments\n\n")
			b.WriteString("          ion secret add <name> <value> \n")
			b.WriteString("          ion secret add -s <salt> <name> <value> \n")
			b.WriteString("          ion secret add -t <tag1> <tag2> <name> <value>\n\n")
			b.WriteString(" if no salt is provided, a random salt will be generated\n\n")
			b.WriteString(" if no tags are provided, the secret will not have any tags\n\n")
			b.WriteString(" example: ion secret add -s mysalt -t tag1 tag2 name value")
			return CommandFinishedMsg{
				Err:     errors.New("ion secret add error"),
				Command: strings.Join(args, " "),
				Output:  b.String(),
				NewDir:  currentDir,
			}
		}
	}

	err := dataRef.AddSecret(args, configPath)
	if err != nil {
		return func() tea.Msg {
			return CommandFinishedMsg{
				Err:     err,
				Command: strings.Join(args, " "),
				Output:  err.Error(),
				NewDir:  currentDir,
			}
		}
	}
	return func() tea.Msg {
		var b strings.Builder
		b.WriteString(args[(len(args) - 2)])
		b.WriteString(" added to secrets")
		return CommandFinishedMsg{
			Command: "ion secret add ",
			Output:  b.String(),
			NewDir:  currentDir,
		}
	}
}

func updateSecret(args []string, configPath string, dataRef data.IData) tea.Cmd {
	if len(args) < 2 {
		return func() tea.Msg {
			var b strings.Builder
			b.WriteString("ion secret update accepts two arguments\n\n")
			b.WriteString("          ion secret update <name> <new-value> \n")
			b.WriteString(" example: ion secret update cool-name cool-value")
			return CommandFinishedMsg{
				Err:     errors.New("ion secret update error"),
				Command: strings.Join(args, " "),
				Output:  b.String(),
				NewDir:  currentDir,
			}
		}
	}

	err := dataRef.UpdateSecret(args, configPath)
	if err != nil {
		return func() tea.Msg {
			return CommandFinishedMsg{
				Err:     err,
				Command: strings.Join(args, " "),
				Output:  err.Error(),
				NewDir:  currentDir,
			}
		}
	}
	return func() tea.Msg {
		var b strings.Builder
		b.WriteString(args[(len(args) - 2)])
		b.WriteString(" updated")
		return CommandFinishedMsg{
			Command: "ion secret update ",
			Output:  b.String(),
			NewDir:  currentDir,
		}
	}
}
