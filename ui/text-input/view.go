package textinput

import (
	"strings"

	"github.com/aykleo/ion/ui/styles"
)

func (m Input) View() string {
	var b strings.Builder

	b.WriteString(m.input.View())

	b.WriteString(styles.NoStyle.Render("\n\nCursor mode is "))
	b.WriteString(styles.NoStyle.Render(m.cursorMode.String()))
	b.WriteString(styles.NoStyle.Render(" (ctrl+r to change style)\n\n"))

	return b.String()
}
