package exec

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func help(args []string) tea.Cmd {
	var b strings.Builder
	b.WriteString("ion help\n\n")
	b.WriteString("          ion zen - toggle zen mode\n")
	b.WriteString("          ion help - show this help\n")
	b.WriteString("          ion alias help - show the help for the alias command\n")
	b.WriteString("          ion secret help - show the help for the secret command\n")

	return func() tea.Msg {
		return CommandFinishedMsg{
			IsSystemCmd: true,
			Command:     "ion help",
			Output:      b.String(),
			NewDir:      currentDir,
		}
	}
}

func toggleZenMode(args []string) tea.Cmd {
	zenMsg := ToggleZenModeMsg{}
	cmdFinishedMsg := CommandFinishedMsg{
		Command: "ion zen",
		Output:  "changed zen mode",
		NewDir:  currentDir,
	}
	return tea.Batch(
		func() tea.Msg { return zenMsg },
		func() tea.Msg { return cmdFinishedMsg },
	)
}
