package main

import (
	"strings"

	"github.com/aykleo/ion/ui/styles"
)

func (m terminal) View() string {
	if m.err != nil {
		return "Error: " + m.err.Error() + "\n"
	}

	return m.UserInput()
}

func (m terminal) UserInput() string {
	var b strings.Builder
	username := m.storage.GetUser().Username
	b.WriteString(styles.MainTheme.Render("~/"))
	b.WriteString(styles.MainTheme.Render(username))
	b.WriteString(styles.MainTheme.Render(" "))
	content := styles.JoinHorizontal(b.String(), m.input.View())

	content = styles.TerminalStyle.Render(content)
	return content
}
