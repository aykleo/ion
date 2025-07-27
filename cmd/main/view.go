package main

import "github.com/charmbracelet/lipgloss"

func (m terminal) View() string {
	if m.err != nil {
		return "Error: " + m.err.Error() + "\n"
	}
	terminalStyle := lipgloss.NewStyle().Padding(1, 2)
	input := terminalStyle.Render(m.input.View())
	return input
}
