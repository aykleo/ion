package main

import (
	"fmt"
	"os"

	"github.com/aykleo/ion/config"
	"github.com/aykleo/ion/data"
	"github.com/aykleo/ion/data/sqlite"
	pager "github.com/aykleo/ion/ui/pager"
	textinput "github.com/aykleo/ion/ui/text-input"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	config := config.Init()
	sqlite.InitSQLite()
	input := textinput.NewTextInput()
	data := data.NewData()
	dataFields, exists := data.GetOrCreateDataFields(config.GetPath())
	if exists {
		data = dataFields
	}
	folder := getFolderFromOs()
	pager := pager.NewPager()
	m := terminal{
		input:         input,
		data:          data,
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
