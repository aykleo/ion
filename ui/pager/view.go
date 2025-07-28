package pager

import (
	"fmt"
	"strings"

	"github.com/aykleo/ion/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

func (m *Pager) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(m.currentPath), m.viewport.View(), m.footerView())
}

func (m *Pager) headerView(t string) string {
	title := styles.PagerTitleStyle.Render(t)
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	line = styles.MainTheme.Render(line)
	return styles.JoinHorizontal(title, line)
}

func (m *Pager) footerView() string {
	info := styles.PagerInfoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	line = styles.MainTheme.Render(line)
	return styles.JoinHorizontal(line, info)
}
