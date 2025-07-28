package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/aykleo/ion/storage"
	pager "github.com/aykleo/ion/ui/pager"
	textinput "github.com/aykleo/ion/ui/text-input"
	tea "github.com/charmbracelet/bubbletea"
)

func openEditor() tea.Cmd {

	c := exec.Command("powershell", "-NoExit", "-Command", "Set-Location C:\\dev\\ion")
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err}
	})
}

func main() {
	input := textinput.NewTextInput()
	storage := storage.NewStorage()
	folder := getFolderFromOs()
	pager := pager.NewPager()
	m := terminal{
		input:         input,
		storage:       storage,
		currentFolder: folder,
		pager:         pager,
	}
	pager.SetCurrentPath(folder)

	if _, err := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
