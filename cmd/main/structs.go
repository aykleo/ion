package main

import (
	"github.com/aykleo/ion/storage"
	textinput "github.com/aykleo/ion/ui/text-input"
	tea "github.com/charmbracelet/bubbletea"
)

type editorFinishedMsg struct{ err error }

type terminal struct {
	width   int
	height  int
	err     error
	storage *storage.Storage
	input   *textinput.Input
}

func (m terminal) Init() tea.Cmd {
	return nil
}
