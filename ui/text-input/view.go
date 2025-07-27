package textinput

import (
	"strings"
)

func (m Input) View() string {
	var b strings.Builder

	b.WriteString(m.input.View())

	return b.String()
}
