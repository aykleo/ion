package exec

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aykleo/ion/data"
	"github.com/aykleo/ion/ui/styles"
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
			b.WriteString("ion secret add accepts two optionals flags and two arguments\n\n")
			b.WriteString("          ion secret add <name> <value>\n")
			b.WriteString("          ion secret add -s <salt> <name> <value>\n")
			b.WriteString("          ion secret add -t <tag1> <tag2> <name> <value>\n\n")
			b.WriteString(" if no salt is not provided, a random salt will be generated\n")
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
			b.WriteString("ion secret update accepts two arguments\n\n")
			b.WriteString("          ion secret update <name> <new-value> \n\n")
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
			b.WriteString("ion secret rename accepts two arguments\n\n")
			b.WriteString("          ion secret rename <name> <new-name> \n\n")
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
			b.WriteString("ion secret tag accepts many arguments but the the last one should always be the name of the secret\n\n")
			b.WriteString("          ion secret tag <tag1> <tag2> <name> \n\n")
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

func listSecrets(args []string, configPath string, dataRef data.IData) tea.Cmd {
	secrets, isJson, err := dataRef.ListSecrets(args, configPath)
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
	if len(secrets) == 0 {
		return func() tea.Msg {
			return CommandFinishedMsg{
				Err:     errors.New("no secrets found"),
				Command: "ion secret list",
				Output:  "no secrets found, try adding one with ion secret add <name> <value>",
				NewDir:  currentDir,
			}
		}
	}
	if isJson {
		return func() tea.Msg {
			return CommandFinishedMsg{
				Command: "ion secret list",
				Output:  styles.FormatSecretsAsJSON(secrets),
				NewDir:  currentDir,
			}
		}
	}
	return func() tea.Msg {
		var b strings.Builder
		b.WriteString("Name             Value         Salt         Tags                    Updated\n")
		b.WriteString("----             -----         ----         ----                    -------\n")
		for _, secret := range secrets {
			name := secret.ID
			if len(name) > 15 {
				name = name[:12] + "..."
			}
			b.WriteString(fmt.Sprintf("%-17s", name))

			value := secret.Value
			if len(value) > 10 {
				value = value[:10] + "..."
			}
			b.WriteString(fmt.Sprintf("%-14s", value))

			salt := secret.Salt
			if len(salt) > 12 {
				salt = salt[:9] + "..."
			}
			b.WriteString(fmt.Sprintf("%-13s", salt))

			tags := strings.Join(secret.Tags, ",")
			if len(tags) > 20 {
				tags = tags[:17] + "..."
			}
			b.WriteString(fmt.Sprintf("%-23s", tags))

			updatedAt := secret.UpdatedAt.Format("2006-01-02")
			b.WriteString(updatedAt)

			b.WriteString("\n")
		}

		return CommandFinishedMsg{
			Command: "ion secret list",
			Output:  b.String(),
			NewDir:  currentDir,
		}
	}
}

func searchSecret(args []string, configPath string, dataRef data.IData) tea.Cmd {
	secret, err := dataRef.SearchSecret(args)
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
			Command: "ion secret search",
			Output:  styles.FormatSecretsAsJSON(secret),
			NewDir:  currentDir,
		}
	}
}

func removeSecret(args []string, configPath string, dataRef data.IData) tea.Cmd {
	err := dataRef.RemoveSecret(args, configPath)
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
		b.WriteString(args[(len(args) - 1)])
		b.WriteString(" removed")
		return CommandFinishedMsg{
			Command: "ion secret remove",
			Output:  b.String(),
			NewDir:  currentDir,
		}
	}
}

func copySecretToClipboard(args []string, configPath string, dataRef data.IData) tea.Cmd {
	err := dataRef.CopySecretToClipboard(args, configPath)
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
		b.WriteString(args[(len(args) - 1)])
		b.WriteString(" copied to clipboard")
		return CommandFinishedMsg{
			Command: "ion secret copy",
			Output:  b.String(),
			NewDir:  currentDir,
		}
	}
}
