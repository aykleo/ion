package textinput

type ITextInput interface {
	SetWidth(width int)
}

func (m *Input) SetWidth(width int) {
	m.input.Width = width / 2
}
