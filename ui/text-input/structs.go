package textinput

import (
	"github.com/aykleo/ion/ui/styles"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
)

type (
	errMsg error
)

type Input struct {
	input      textinput.Model
	cursorMode cursor.Mode
	err        error
}

func InitInput() Input {
	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.Focus()
	ti.PromptStyle = styles.MainTheme
	ti.TextStyle = styles.MainTheme
	ti.Cursor.Style = styles.MainTheme
	ti.Width = 20

	return Input{
		input: ti,
		err:   nil,
	}
}
