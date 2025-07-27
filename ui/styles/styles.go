package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	MainTheme     = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	TerminalStyle = lipgloss.NewStyle().Padding(1, 2)
	NoStyle       = lipgloss.NewStyle()
)

func JoinHorizontal(strs ...string) string {
	return lipgloss.JoinHorizontal(lipgloss.Center, strs...)
}
