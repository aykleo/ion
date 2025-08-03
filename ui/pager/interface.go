package pager

import tea "github.com/charmbracelet/bubbletea"

type IPager interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	View() string

	SetWidth(width int) int
	SetHeight(height int) int
	SetCurrentFolder(path string)

	AppendCommand(msg string) (UpdateContentMsg, tea.Cmd)
	ResetMsgs()
}

func NewPager() IPager {
	pager := InitPager()
	return &pager
}

func (m *Pager) Init() tea.Cmd {
	return nil
}

func (m *Pager) SetWidth(width int) int {
	return int(float64(width) * 0.98)
}
func (m *Pager) SetHeight(height int) int {

	if height == 0 {
		return 0
	}
	if height < 10 {
		if m.zenMode {
			return int(float64(height) * 0.45)
		}
		return int(float64(height) * 0.4)
	}
	if height < 20 {
		if m.zenMode {
			return int(float64(height) * 0.75)
		}
		return int(float64(height) * 0.7)
	}
	if height < 30 {
		if m.zenMode {
			return int(float64(height) * 0.85)
		}
		return int(float64(height) * 0.8)
	}
	if height < 40 {
		if m.zenMode {
			return int(float64(height) * 0.9)
		}
		return int(float64(height) * 0.85)
	}
	if m.zenMode {
		return int(float64(height) * 0.95)
	}
	return int(float64(height) * 0.9)

}

func (m *Pager) SetCurrentFolder(path string) {
	m.currentPath = path
}

func (m *Pager) AppendCommand(msg string) (UpdateContentMsg, tea.Cmd) {
	contentSize := len(m.content)
	if contentSize < 200 {
		m.content = append(m.content, msg)
	} else {
		m.content = append(m.content[1:], msg)
	}

	return UpdateContentMsg{}, func() tea.Msg {
		return UpdateContentMsg{}
	}
}

func (m *Pager) ResetMsgs() {
	m.content = []string{}
}
