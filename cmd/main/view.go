package main

import (
	"strings"

	"github.com/aykleo/ion/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

func (m terminal) View() string {
	if m.err != nil {
		return "Error: " + m.err.Error() + "\n"
	}
	pager := lipgloss.NewStyle().Padding(1, 0, 0, 0).Render(m.pager.View())
	content := styles.JoinVertical(m.UserInput(), pager)
	content = styles.TerminalStyle.Render(content)
	return content
}

func (m terminal) UserInput() string {
	var b strings.Builder
	username := m.storage.GetUser().Username
	b.WriteString(styles.MainTheme.Render("~/"))
	b.WriteString(styles.MainTheme.Render(username))
	b.WriteString(styles.MainTheme.Render(" "))
	content := styles.JoinHorizontal(b.String(), m.input.View())

	return content
}
