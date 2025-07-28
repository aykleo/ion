package textinput

import (
	"github.com/aykleo/ion/ui/styles"
	"github.com/charmbracelet/bubbles/textinput"
)

func InitInput() Input {
	ti := textinput.New()
	ti.Placeholder = "Type your command"
	ti.PlaceholderStyle = styles.Placeholder
	ti.Focus()
	ti.PromptStyle = styles.MainTheme
	ti.TextStyle = styles.MainTheme
	ti.Cursor.Style = styles.MainTheme
	ti.Width = 20

	firstWordStyle := styles.MainTheme

	return Input{
		input:          ti,
		err:            nil,
		FirstWordStyle: firstWordStyle,
	}
}
