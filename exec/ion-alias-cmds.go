package exec

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aykleo/ion/data"
	"github.com/aykleo/ion/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

func addAlias(args []string, configPath string, dataRef data.IData) tea.Cmd {
	if len(args) < 1 {
		return func() tea.Msg {
			var b strings.Builder
			b.WriteString("ion alias add accepts arguments with a '=' in between\n\n")
			b.WriteString("          ion alias add <alias-name>=<command>\n")
			b.WriteString("      or: ion alias add <alias-name> = <command>\n\n")
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

func updateAlias(args []string, configPath string, dataRef data.IData) tea.Cmd {
	if len(args) < 1 {
		return func() tea.Msg {
			var b strings.Builder
			b.WriteString("ion alias update accepts arguments with a '=' in between\n\n")
			b.WriteString("          ion alias update <alias-name>=<new-command>\n")
			b.WriteString("      or: ion alias update <alias-name> = <new-command>\n\n")
			b.WriteString(" example: ion alias update my-alias=ion secret list\n")
			return CommandFinishedMsg{
				Err:     errors.New("ion alias update error"),
				Command: strings.Join(args, " "),
				Output:  b.String(),
				NewDir:  currentDir,
			}
		}
	}

	err := dataRef.UpdateAlias(args, configPath)
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
		b.WriteString("alias ")
		b.WriteString(aliasName)
		b.WriteString(" updated with a new command")
		return CommandFinishedMsg{
			Command: "ion alias update",
			Output:  b.String(),
			NewDir:  currentDir,
		}
	}
}

func renameAlias(args []string, configPath string, dataRef data.IData) tea.Cmd {
	if len(args) < 2 {
		return func() tea.Msg {
			var b strings.Builder
			b.WriteString("ion alias rename accepts two arguments\n\n")
			b.WriteString("          ion alias rename <old-name> <new-name> \n\n")
			b.WriteString(" example: ion alias rename old-alias new-alias")
			return CommandFinishedMsg{
				Err:     errors.New("ion alias rename error"),
				Command: strings.Join(args, " "),
				Output:  b.String(),
				NewDir:  currentDir,
			}
		}
	}

	err := dataRef.RenameAlias(args, configPath)
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
		b.WriteString("alias ")
		b.WriteString(args[0])
		b.WriteString(" renamed to ")
		b.WriteString(args[1])
		return CommandFinishedMsg{
			Command: "ion alias rename",
			Output:  b.String(),
			NewDir:  currentDir,
		}
	}
}

func removeAlias(args []string, configPath string, dataRef data.IData) tea.Cmd {
	if len(args) < 1 {
		return func() tea.Msg {
			var b strings.Builder
			b.WriteString("ion alias remove accepts one argument\n\n")
			b.WriteString("          ion alias remove <alias-name> \n\n")
			b.WriteString(" example: ion alias remove my-alias")
			return CommandFinishedMsg{
				Err:     errors.New("ion alias remove error"),
				Command: strings.Join(args, " "),
				Output:  b.String(),
				NewDir:  currentDir,
			}
		}
	}

	err := dataRef.RemoveAlias(args, configPath)
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
		b.WriteString("alias ")
		b.WriteString(args[0])
		b.WriteString(" removed")
		return CommandFinishedMsg{
			Command: "ion alias remove",
			Output:  b.String(),
			NewDir:  currentDir,
		}
	}
}

func listAliases(args []string, configPath string, dataRef data.IData) tea.Cmd {
	aliases, isJson, err := dataRef.ListAliases(args, configPath)
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
	if len(aliases) == 0 {
		return func() tea.Msg {
			return CommandFinishedMsg{
				Err:     errors.New("no aliases found"),
				Command: "ion alias list",
				Output:  "no aliases found, try adding one with ion alias add <name>=<command>",
				NewDir:  currentDir,
			}
		}
	}
	if isJson {
		return func() tea.Msg {
			return CommandFinishedMsg{
				Command: "ion alias list",
				Output:  styles.FormatAliasesAsJSON(aliases),
				NewDir:  currentDir,
			}
		}
	}
	return func() tea.Msg {
		var b strings.Builder
		b.WriteString("Name             Command                          Updated\n")
		b.WriteString("----             -------                          -------\n")
		for _, alias := range aliases {
			name := alias.Name
			if len(name) > 15 {
				name = name[:12] + "..."
			}
			b.WriteString(fmt.Sprintf("%-17s", name))

			command := alias.Command
			if len(command) > 30 {
				command = command[:27] + "..."
			}
			b.WriteString(fmt.Sprintf("%-33s", command))

			updatedAt := alias.UpdatedAt.Format("2006-01-02")
			b.WriteString(updatedAt)

			b.WriteString("\n")
		}

		return CommandFinishedMsg{
			Command: "ion alias list",
			Output:  b.String(),
			NewDir:  currentDir,
		}
	}
}

func searchAliases(args []string, configPath string, dataRef data.IData) tea.Cmd {
	alias, err := dataRef.SearchAlias(args)
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
		return CommandFinishedMsg{
			IsSystemCmd: true,
			Command:     "ion alias search",
			Output:      styles.FormatSearchedAlias(alias[0]),
			NewDir:      currentDir,
		}
	}
}

func helpAlias(args []string, configPath string, dataRef data.IData) tea.Cmd {
	var b strings.Builder
	b.WriteString("ion alias help\n\n")
	b.WriteString("          ion alias add <alias-name>=<command> - add a new alias\n")
	b.WriteString("          ion alias update <alias-name>=<new-command> - update the command of an alias\n")
	b.WriteString("          ion alias rename <old-name> <new-name> - rename an alias\n")
	b.WriteString("          ion alias remove <alias-name> - remove an alias\n")
	b.WriteString("          ion alias list - list all aliases, you can use -j to get the output as json\n")
	b.WriteString("          ion alias search <alias-name> - search for an alias\n\n")
	b.WriteString("          commands can be transformed dynamically with the use of ${ion} on the place of the arg you want to make dynamic\n")
	b.WriteString("          creation example: ion alias add new-alias=echo ${ion} said ${ion}\n")
	b.WriteString("          usage example: new-alias Bob Hello\n")
	b.WriteString("          output: Bob said Hello\n\n")
	return func() tea.Msg {
		return CommandFinishedMsg{
			IsSystemCmd: true,
			Command:     "ion alias help",
			Output:      b.String(),
			NewDir:      currentDir,
		}
	}
}
