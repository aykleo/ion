package textinput

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type (
	errMsg error
)

type Input struct {
	input          textinput.Model
	cursorMode     cursor.Mode
	err            error
	FirstWordStyle lipgloss.Style
	zenMode        bool
}

type CommandMsg struct {
	Alias        string
	Command      string
	Args         []string
	IsIonCommand bool
}
