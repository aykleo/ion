package textinput

import (
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyEnter:
			if m.input.Value() != "" {
				currentMessage := m.input.Value()
				inputCmd, err := m.doCommand(currentMessage)
				if err != nil {
					m.err = err
					return m, nil
				}
				return m, inputCmd
			}
		}
		switch msg.String() {

		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > cursor.CursorHide {
				m.cursorMode = cursor.CursorBlink
			}
			return m, m.input.Cursor.SetMode(m.cursorMode)
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *Input) doCommand(msg string) (tea.Cmd, error) {
	defer m.input.Reset()
	isIonCommand := strings.HasPrefix(msg, "ion")

	parts := strings.Fields(msg)
	var args []string

	if len(parts) > 1 {
		args = parts[1:]
	}

	if isIonCommand {
		args = parts
	}

	var commandName string
	if len(parts) > 0 {
		commandName = parts[0]
	}

	alias := m.input.Value()
	command := CommandMsg{
		Alias:        alias,
		Command:      commandName,
		Args:         args,
		IsIonCommand: isIonCommand,
	}

	return func() tea.Msg {
		return command
	}, nil
}
