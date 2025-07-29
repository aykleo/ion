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
				Err:     errors.New("ion user rename <username>"),
				Command: strings.Join(args, " "),
				Output:  "ion user needs one argument, try [ ion user rename <username> ]",
				NewDir:  currentDir,
			}
		}
	}
	username := args[0]
	currentUsername := dataRef.GetUser().Username
	dataRef.SetUsername(username, configPath)
	return func() tea.Msg {
		return CommandFinishedMsg{
			Command: "ion user rename " + username,
			Output:  "username " + currentUsername + " renamed to " + username,
			NewDir:  currentDir,
		}
	}
}

func addSecret(args []string, configPath string, dataRef data.IData) tea.Cmd {
	if len(args) < 2 {
		return func() tea.Msg {
			var b strings.Builder
			b.WriteString("ion secret add accepts two optionals flags and two arguments\n")
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

func updateSecretValue(args []string, configPath string, dataRef data.IData) tea.Cmd {
	if len(args) < 2 {
		return func() tea.Msg {
			var b strings.Builder
			b.WriteString("ion secret update accepts two arguments\n")
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

	err := dataRef.UpdateSecretValue(args, configPath)
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
		b.WriteString("secret ")
		b.WriteString(args[(len(args) - 2)])
		b.WriteString(" updated with a new value")
		return CommandFinishedMsg{
			Command: "ion secret update",
			Output:  b.String(),
			NewDir:  currentDir,
		}
	}
}

func updateSecretName(args []string, configPath string, dataRef data.IData) tea.Cmd {
	if len(args) < 2 {
		return func() tea.Msg {
			var b strings.Builder
			b.WriteString("ion secret rename accepts two arguments\n")
			b.WriteString("          ion secret rename <name> <new-name> \n")
			b.WriteString(" example: ion secret rename name cooler-name")
			return CommandFinishedMsg{
				Err:     errors.New("ion secret update error"),
				Command: strings.Join(args, " "),
				Output:  b.String(),
				NewDir:  currentDir,
			}
		}
	}

	err := dataRef.UpdateSecretName(args, configPath)
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
		b.WriteString("secret ")
		b.WriteString(args[(len(args) - 2)])
		b.WriteString(" renamed to ")
		b.WriteString(args[(len(args) - 1)])
		return CommandFinishedMsg{
			Command: "ion secret rename",
			Output:  b.String(),
			NewDir:  currentDir,
		}
	}
}

func updateSecretTags(args []string, configPath string, dataRef data.IData) tea.Cmd {
	if len(args) < 2 {
		return func() tea.Msg {
			var b strings.Builder
			b.WriteString("ion secret tag accepts many arguments but the the last one should always be the name of the secret\n")
			b.WriteString("          ion secret tag <tag1> <tag2> <name> \n")
			b.WriteString(" example: ion secret tag tag1 tag2 name")
			return CommandFinishedMsg{
				Err:     errors.New("ion secret tag error"),
				Command: strings.Join(args, " "),
				Output:  b.String(),
				NewDir:  currentDir,
			}
		}
	}
	err := dataRef.UpdateSecretTags(args, configPath)
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
		b.WriteString("tags for secret ")
		b.WriteString(args[(len(args) - 1)])
		b.WriteString(" updated")
		return CommandFinishedMsg{
			Command: "ion secret tags",
			Output:  b.String(),
			NewDir:  currentDir,
		}
	}
}
