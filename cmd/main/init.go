package main

import tea "github.com/charmbracelet/bubbletea"

func (m terminal) Init() tea.Cmd {
	return m.input.Init()
}
