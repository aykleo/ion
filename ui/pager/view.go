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
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m *Pager) footerView() string {
	info := styles.PagerInfoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}
