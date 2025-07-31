package main

import (
	"strings"

	"github.com/aykleo/ion/data"
	"github.com/aykleo/ion/exec"
	"github.com/aykleo/ion/ui/styles"
	textinput "github.com/aykleo/ion/ui/text-input"
	tea "github.com/charmbracelet/bubbletea"
)

func (m terminal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		default:
			_, pagerCmd := m.pager.Update(msg)
			_, inputCmd := m.input.Update(msg)
			return m, tea.Batch(inputCmd, pagerCmd)
		}

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.input.SetWidth(m.width)
		_, pagerCmd := m.pager.Update(msg)
		return m, pagerCmd

	case textinput.CommandMsg:
		isIonCommand := msg.IsIonCommand
		fullCommand := msg.Command
		aliasArgs := checkForAlias(msg.Command, m.data)
		if aliasArgs != nil {
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
		formattedCommand := styles.FormatCommandPrompt(fullCommand, m.data.GetUser().Username)
		_, pagerCmd := m.pager.AppendCommand(formattedCommand)
		if isIonCommand {
			ionCmd := exec.ExecIonCommand(msg.Args, m.data)
			return m, tea.Batch(pagerCmd, ionCmd)
		} else {
			execCmd := exec.ExecSysCommand(msg.Command, msg.Args)
			return m, tea.Batch(pagerCmd, execCmd)
		}

	case exec.CommandFinishedMsg:
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
			formattedOutput := styles.FormatCommandOutput(msg.Output)
			_, pagerCmdOutput := m.pager.AppendCommand(formattedOutput)
			return m, pagerCmdOutput
		}

	}

	_, inputCmd := m.input.Update(msg)
	_, pagerCmd := m.pager.Update(msg)
	return m, tea.Batch(inputCmd, pagerCmd)
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
