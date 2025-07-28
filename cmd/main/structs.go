package main

import (
	"os"

	"github.com/aykleo/ion/storage"
	pager "github.com/aykleo/ion/ui/pager"
	textinput "github.com/aykleo/ion/ui/text-input"
	tea "github.com/charmbracelet/bubbletea"
)

type editorFinishedMsg struct{ err error }

type terminal struct {
	width         int
	height        int
	err           error
	currentFolder string
	storage       storage.IStorage
	input         textinput.ITextInput
	pager         pager.IPager
}

func (m terminal) Init() tea.Cmd {
	return m.input.Init()
}

func getFolderFromOs() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}

func (m *terminal) setCurrentFolder() {
	m.currentFolder = getFolderFromOs()
}

func (m *terminal) getCurrentFolder() string {
	return m.currentFolder
}
