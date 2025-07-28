package main

import (
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
	case editorFinishedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.input.SetWidth(m.width)
		_, pagerCmd := m.pager.Update(msg)
		return m, pagerCmd
	}

	_, inputCmd := m.input.Update(msg)
	_, pagerCmd := m.pager.Update(msg)
	return m, tea.Batch(inputCmd, pagerCmd)
}
