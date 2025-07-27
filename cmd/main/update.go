package main

import (
	textinput "github.com/aykleo/ion/ui/text-input"
	tea "github.com/charmbracelet/bubbletea"
)

func (m terminal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		default:
			updatedModel, inputCmd := m.input.Update(msg)
			*m.input = updatedModel.(textinput.Input)
			return m, inputCmd
		}
	case editorFinishedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.input.SetWidth(m.width)
		return m, nil
	}

	updatedModel, cmd := m.input.Update(msg)
	*m.input = updatedModel.(textinput.Input)
	return m, cmd
}
