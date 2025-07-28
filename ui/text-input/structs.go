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
}
