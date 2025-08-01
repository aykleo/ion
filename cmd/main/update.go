package main

import (
	"strings"

	"github.com/aykleo/ion/exec"
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
		var fullCmd string
		if msg.IsIonCommand {
			fullCmd = strings.Join(msg.Args, " ")
		} else {
			if len(msg.Args) > 0 {
				fullCmd = msg.Command + " " + strings.Join(msg.Args, " ")
			} else {
				fullCmd = msg.Command
			}
		}
		pipeCmds, isPipeCmd := checkForPipeCmds([]string{fullCmd})
		if isPipeCmd {
			return m.executePipeCmds(pipeCmds)
		}
		return m.execute(msg)

	case exec.CommandFinishedMsg:
		return m.finishCommand(msg)
	}

	_, inputCmd := m.input.Update(msg)
	_, pagerCmd := m.pager.Update(msg)
	return m, tea.Batch(inputCmd, pagerCmd)
}
