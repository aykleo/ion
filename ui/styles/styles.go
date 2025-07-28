package styles

import (
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

const (
	MainThemeColor   = "135"
	FolderColor      = "2"
	SuccessColor     = "2"
	ErrorColor       = "1"
	CommandColor     = "4"
	OutputColor      = "7"
	TableHeaderColor = "6"
	TableBorderColor = "8"
)

var (
	MainTheme     = lipgloss.NewStyle().Foreground(lipgloss.Color(MainThemeColor))
	Placeholder   = MainTheme.Italic(true).Faint(true)
	FolderStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color(FolderColor)).Bold(false).Padding(1).Background(lipgloss.Color(MainThemeColor))
	TerminalStyle = lipgloss.NewStyle().Padding(1, 2)
	NoStyle       = lipgloss.NewStyle()
)

var (
	PagerTitleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return MainTheme.BorderStyle(b).BorderForeground(lipgloss.Color(MainThemeColor)).Padding(0, 1)
	}()

	PagerInfoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return PagerTitleStyle.BorderStyle(b).BorderForeground(lipgloss.Color(MainThemeColor))
	}()

	CommandPromptStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(CommandColor)).
				Bold(true).
				Padding(0, 1)

	CommandOutputStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(OutputColor)).
				Padding(0, 2)

	SuccessMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(SuccessColor)).
				Bold(true).
				Padding(0, 1)

	ErrorMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ErrorColor)).
				Bold(true).
				Padding(0, 1)

	PagerContentStyle = lipgloss.NewStyle().
				Padding(0, 2)

	TableHeaderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(TableHeaderColor)).
				Bold(true).
				Padding(0, 1)

	TableRowStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(OutputColor)).
			Padding(0, 1)

	TableBorderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(TableBorderColor))

	DirectoryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(FolderColor)).
			Bold(true)

	FileStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(OutputColor))
)

func JoinHorizontal(strs ...string) string {
	return lipgloss.JoinHorizontal(lipgloss.Center, strs...)
}

func JoinVertical(strs ...string) string {
	return lipgloss.JoinVertical(lipgloss.Left, strs...)
}

func FormatCommandPrompt(command, username string) string {
	var b strings.Builder
	date := time.Now().Format("2006-01-02 15:04:05")
	b.WriteString(MainTheme.Render(date))
	b.WriteString(MainTheme.Render(" "))
	b.WriteString(MainTheme.Render("~/"))
	b.WriteString(MainTheme.Render(username))
	b.WriteString(MainTheme.Render(" "))
	b.WriteString(MainTheme.Render("> "))
	b.WriteString(MainTheme.Render(command))
	return CommandPromptStyle.Render(b.String())
}

func FormatCommandOutput(output string) string {
	if output == "" {
		return ""
	}

	trimmedOutput := strings.TrimSpace(output)
	if trimmedOutput == "" {
		return ""
	}

	if isTableOutput(trimmedOutput) {
		return formatTable(trimmedOutput)
	}

	return SuccessMessageStyle.Render(trimmedOutput)
}

func FormatSuccessMessage(message string) string {
	return SuccessMessageStyle.Render("✓ " + message)
}

func FormatErrorMessage(message string) string {
	return ErrorMessageStyle.Render("✗ " + message)
}
