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
	db, err := sqlite.InitSQLite()
	if err != nil {
		fmt.Printf("Error initializing SQLite: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	input := textinput.NewTextInput()
	data := data.NewData()
	data.SetDB(db)
	data.GetInitialData(config.GetPath())
	folder := getFolderFromOs()
	pager := pager.NewPager()
	m := terminal{
		input:         input,
		data:          data,
		db:            db,
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
