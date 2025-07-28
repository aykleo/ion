package main

import (
	"fmt"
	"os"

	"github.com/aykleo/ion/config"
	"github.com/aykleo/ion/storage"
	pager "github.com/aykleo/ion/ui/pager"
	textinput "github.com/aykleo/ion/ui/text-input"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	config := config.Init()
	input := textinput.NewTextInput()
	storage := storage.NewStorage()
	folder := getFolderFromOs()
	pager := pager.NewPager()
	m := terminal{
		input:         input,
		storage:       storage,
		currentFolder: folder,
		pager:         pager,
		config:        config,
	}
	pager.SetCurrentFolder(folder)

	if _, err := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}
