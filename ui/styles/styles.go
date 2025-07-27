package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	MainTheme     = lipgloss.NewStyle().Foreground(lipgloss.Color("135"))
	Placeholder   = MainTheme.Italic(true).Faint(true)
	TerminalStyle = lipgloss.NewStyle().Padding(1, 2)
	NoStyle       = lipgloss.NewStyle()
)

func JoinHorizontal(strs ...string) string {
	return lipgloss.JoinHorizontal(lipgloss.Center, strs...)
}
