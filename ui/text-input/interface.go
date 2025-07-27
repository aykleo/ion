package textinput

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ITextInput interface {
	SetWidth(width int)
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	View() string
	Init() tea.Cmd
}

func (m *Input) SetWidth(width int) {
	m.input.Width = width / 2
}

func (m *Input) Init() tea.Cmd {
	return textinput.Blink
}

func NewTextInput() ITextInput {
	input := InitInput()
	return &input
}
