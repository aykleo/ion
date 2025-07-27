package textinput

import (
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
