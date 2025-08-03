package main

import (
	"fmt"
	"strings"

	"github.com/aykleo/ion/data"
	"github.com/aykleo/ion/exec"
	"github.com/aykleo/ion/ui/styles"
	textinput "github.com/aykleo/ion/ui/text-input"
	tea "github.com/charmbracelet/bubbletea"
)

const CustomAliasPlaceholder = "${ion}"

func checkForPipeCmds(cmd []string) ([]textinput.CommandMsg, bool) {
	separator := "&&"
	fullCommand := strings.Join(cmd, " ")

	if !strings.Contains(fullCommand, separator) {
		return nil, false
	}

	commands := strings.Split(fullCommand, separator)
	var cleanedCommands []textinput.CommandMsg
	ion := "ion"
	for _, cmd := range commands {
		trimmed := strings.TrimSpace(cmd)
		if trimmed != "" {
			isIonCommand := strings.HasPrefix(trimmed, ion)
			parts := strings.Fields(trimmed)
			var args []string
			var commandName string

			if len(parts) > 0 {
				commandName = parts[0]
				if len(parts) > 1 {
					args = parts[1:]
				}
			}

			if isIonCommand {
				args = parts
			}

			cleanedCommands = append(cleanedCommands, textinput.CommandMsg{
				Alias:        trimmed,
				Command:      commandName,
				Args:         args,
				IsIonCommand: isIonCommand,
			})
		}
	}

	return cleanedCommands, true
}

func checkForAlias(command string, data data.IData) []string {
	aliases, _, _ := data.ListAliases([]string{}, "")
	for _, alias := range aliases {
		if alias.Name == command {
			return strings.Split(alias.Command, " ")
		}
	}
	return nil
}

func (m *terminal) executePipeCmds(pipeCmds []textinput.CommandMsg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	for _, pipeCmd := range pipeCmds {
		_, cmd := m.execute(pipeCmd)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return m, tea.Sequence(cmds...)
}

func (m *terminal) execute(msg textinput.CommandMsg) (tea.Model, tea.Cmd) {
	isIonCommand := msg.IsIonCommand
	fullCommand := msg.Command
	aliasArgs := checkForAlias(msg.Command, m.data)
	if aliasArgs != nil {
		if containsDynamicPlaceholders(aliasArgs) {
			requiredArgs := countDynamicPlaceholders(aliasArgs)
			if len(msg.Args) < requiredArgs {
				return m.handleInsufficientArgsError(msg.Command, requiredArgs, len(msg.Args))
			}
			aliasArgs = translateDynamicAlias(aliasArgs, msg.Args)
			msg.Args = msg.Args[requiredArgs:]
		}

		msg.Command = aliasArgs[0]
		if len(aliasArgs) > 1 {
			msg.Args = append(aliasArgs[1:], msg.Args...)
			if aliasArgs[0] == "ion" {
				msg.Args = append([]string{"ion"}, msg.Args...)
				isIonCommand = true
			}
		}
	}
	if len(msg.Args) > 0 {
		if isIonCommand {
			fullCommand = strings.Join(msg.Args, " ")
		} else {
			fullCommand = msg.Command + " " + strings.Join(msg.Args, " ")
		}
	}
	cmdString := fullCommand
	if aliasArgs != nil && msg.Alias != "" {
		var b strings.Builder
		b.WriteString(msg.Alias)
		b.WriteString(" ")
		b.WriteString(styles.MiscStyle.Render("("))
		b.WriteString(styles.FadedStyle.Render(strings.Join(aliasArgs, " ")))
		b.WriteString(styles.MiscStyle.Render(")"))
		cmdString = b.String()
	}
	formattedCommand := styles.FormatCommandPrompt(cmdString, m.data.GetUser().Username)
	_, pagerCmd := m.pager.AppendCommand(formattedCommand)
	if isIonCommand {
		ionCmd := exec.ExecIonCommand(msg.Args, m.data)
		return m, tea.Batch(pagerCmd, ionCmd)
	} else {
		execCmd := exec.ExecSysCommand(msg.Command, msg.Args)
		return m, tea.Batch(pagerCmd, execCmd)
	}
}

func (m *terminal) finishCommand(msg exec.CommandFinishedMsg) (tea.Model, tea.Cmd) {
	if msg.NewDir != m.currentFolder {
		m.currentFolder = msg.NewDir
		m.pager.SetCurrentFolder(msg.NewDir)
	}

	if msg.Err != nil {
		if msg.Output != "" {
			formattedError := styles.FormatErrorMessage(strings.TrimSpace(msg.Output))
			_, pagerCmdErr := m.pager.AppendCommand(formattedError)
			return m, pagerCmdErr
		}
		formattedError := styles.FormatErrorMessage(msg.Err.Error())
		_, pagerCmdErr := m.pager.AppendCommand(formattedError)
		return m, pagerCmdErr
	}

	if msg.Output != "" {
		formattedOutput := styles.FormatCommandOutput(msg.Output, msg.IsSystemCmd)
		_, pagerCmdOutput := m.pager.AppendCommand(formattedOutput)
		return m, pagerCmdOutput
	}
	return m, nil
}

func containsDynamicPlaceholders(command []string) bool {
	dynamicValue := CustomAliasPlaceholder
	for _, arg := range command {
		if arg == dynamicValue {
			return true
		}
	}
	return false
}

func countDynamicPlaceholders(command []string) int {
	dynamicValue := CustomAliasPlaceholder
	count := 0
	for _, arg := range command {
		if arg == dynamicValue {
			count++
		}
	}
	return count
}

func translateDynamicAlias(originalCommand []string, newValues []string) []string {
	dynamicValue := CustomAliasPlaceholder
	translatedCommand := make([]string, len(originalCommand))
	valueIndex := 0

	for i, arg := range originalCommand {
		if arg == dynamicValue && valueIndex < len(newValues) {
			translatedCommand[i] = newValues[valueIndex]
			valueIndex++
		} else {
			translatedCommand[i] = arg
		}
	}
	return translatedCommand
}

func (m *terminal) handleInsufficientArgsError(aliasName string, required, provided int) (tea.Model, tea.Cmd) {
	errorMsg := styles.FormatErrorMessage(
		fmt.Sprintf("alias '%s' requires %d argument(s), but %d provided",
			aliasName, required, provided))
	_, pagerCmd := m.pager.AppendCommand(errorMsg)
	return m, pagerCmd
}
