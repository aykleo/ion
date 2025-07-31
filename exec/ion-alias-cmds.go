package exec

import (
	"errors"
	"strings"

	"github.com/aykleo/ion/data"
	tea "github.com/charmbracelet/bubbletea"
)

func addAlias(args []string, configPath string, dataRef data.IData) tea.Cmd {
	if len(args) < 1 {
		return func() tea.Msg {
			var b strings.Builder
			b.WriteString("ion alias add accepts arguments with a '=' in between\n\n")
			b.WriteString("          ion alias add <alias-name>=<command>\n")
			b.WriteString("      or: ion alias add <alias-name> = <command>\n")
			b.WriteString(" example: ion alias add new-alias=ion secret add <name> <value>\n")
			return CommandFinishedMsg{
				Err:     errors.New("ion alias add error"),
				Command: strings.Join(args, " "),
				Output:  b.String(),
				NewDir:  currentDir,
			}
		}
	}

	err := dataRef.AddAlias(args, configPath)
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
	aliasName := args[0]
	if strings.Contains(aliasName, "=") {
		aliasName = strings.Split(aliasName, "=")[0]
	}
	return func() tea.Msg {
		var b strings.Builder
		b.WriteString(aliasName)
		b.WriteString(" added to aliases")
		return CommandFinishedMsg{
			Command: "ion alias add ",
			Output:  b.String(),
			NewDir:  currentDir,
		}
	}
}
