package textinput

import (
	"github.com/aykleo/ion/ui/styles"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg error
)

type Input struct {
	input      textinput.Model
	cursorMode cursor.Mode
	err        error
}

func InitialInput() Input {
	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.Focus()
	ti.PromptStyle = styles.CursorStyle
	ti.TextStyle = styles.CursorStyle
	ti.Width = 20

	return Input{
		input: ti,
		err:   nil,
	}
}

func (m Input) Init() tea.Cmd {
	return textinput.Blink
}
