package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	MainTheme     = lipgloss.NewStyle().Foreground(lipgloss.Color("135"))
	Placeholder   = MainTheme.Italic(true).Faint(true)
	FolderStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(false).Padding(1).Background(lipgloss.Color("135"))
	TerminalStyle = lipgloss.NewStyle().Padding(1, 2)
	NoStyle       = lipgloss.NewStyle()
)

var (
	PagerTitleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	PagerInfoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return PagerTitleStyle.BorderStyle(b)
	}()
)

func JoinHorizontal(strs ...string) string {
	return lipgloss.JoinHorizontal(lipgloss.Center, strs...)
}

func JoinVertical(strs ...string) string {
	return lipgloss.JoinVertical(lipgloss.Left, strs...)
}
