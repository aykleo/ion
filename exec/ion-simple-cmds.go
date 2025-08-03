package exec

import (
	tea "github.com/charmbracelet/bubbletea"
)

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
