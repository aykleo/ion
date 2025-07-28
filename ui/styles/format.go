package styles

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func isTableOutput(output string) bool {
	lines := strings.Split(output, "\n")
	if len(lines) < 3 {
		return false
	}

	for _, line := range lines {
		if strings.Contains(line, "Mode") && strings.Contains(line, "LastWriteTime") {
			return true
		}
		if regexp.MustCompile(`^-+\s+-+`).MatchString(strings.TrimSpace(line)) {
			return true
		}
		if regexp.MustCompile(`^d-+\s+`).MatchString(strings.TrimSpace(line)) {
			return true
		}
	}

	return false
}

func formatTable(output string) string {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var formattedLines []string
	var headerFound bool

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			if len(formattedLines) > 0 {
				formattedLines = append(formattedLines, "")
			}
			continue
		}

		if strings.HasPrefix(trimmedLine, "Directory:") {
			styled := TableHeaderStyle.Render("📁 " + trimmedLine)
			formattedLines = append(formattedLines, styled)
			continue
		}

		if strings.Contains(line, "Mode") && strings.Contains(line, "LastWriteTime") {
			headerFound = true
			modeHeader := lipgloss.NewStyle().Width(8).Render("Mode")
			dateTimeHeader := lipgloss.NewStyle().Width(18).Render("LastWriteTime")
			sizeHeader := lipgloss.NewStyle().Width(12).Align(lipgloss.Right).Render("Length")
			nameHeader := lipgloss.NewStyle().Render("Name")

			headerLine := lipgloss.JoinHorizontal(lipgloss.Left,
				modeHeader, "  ",
				dateTimeHeader, " ",
				sizeHeader, "   ",
				nameHeader,
			)

			styled := TableHeaderStyle.Render("┌─ " + headerLine + " ─┐")
			formattedLines = append(formattedLines, styled)
			continue
		}

		if regexp.MustCompile(`^-+\s+-+`).MatchString(trimmedLine) {
			separator := TableBorderStyle.Render(strings.Repeat("─", 65))
			formattedLines = append(formattedLines, separator)
			continue
		}

		if headerFound && regexp.MustCompile(`^[da-z-]+\s+`).MatchString(trimmedLine) {
			formatted := formatDirectoryRow(line)
			formattedLines = append(formattedLines, formatted)
			continue
		}

		formattedLines = append(formattedLines, TableRowStyle.Render(line))
	}

	for len(formattedLines) > 0 && strings.TrimSpace(formattedLines[len(formattedLines)-1]) == "" {
		formattedLines = formattedLines[:len(formattedLines)-1]
	}

	return strings.Join(formattedLines, "\n")
}

func formatDirectoryRow(line string) string {
	fields := regexp.MustCompile(`\s+`).Split(strings.TrimSpace(line), -1)
	if len(fields) < 4 {
		return TableRowStyle.Render(line)
	}

	mode := fields[0]
	date := fields[1]
	time := fields[2]
	var size, name string

	if len(fields) == 4 {
		name = fields[3]
		size = ""
	} else {
		size = fields[3]
		name = strings.Join(fields[4:], " ")
	}

	var styledName string
	if strings.HasPrefix(mode, "d") {
		styledName = DirectoryStyle.Render("📁 " + name)
	} else {
		styledName = FileStyle.Render("📄 " + name)
	}

	modeCol := lipgloss.NewStyle().Width(8).Render(mode)
	dateTimeCol := lipgloss.NewStyle().Width(18).Render(date + " " + time)

	if size == "" {
		sizeCol := lipgloss.NewStyle().Width(12).Render("")
		return TableRowStyle.Render(
			lipgloss.JoinHorizontal(lipgloss.Left,
				modeCol,
				"  ",
				dateTimeCol,
				" ",
				sizeCol,
				"   ",
				styledName,
			),
		)
	} else {
		sizeCol := lipgloss.NewStyle().Width(12).Align(lipgloss.Right).Render(size)
		return TableRowStyle.Render(
			lipgloss.JoinHorizontal(lipgloss.Left,
				modeCol,
				"  ",
				dateTimeCol,
				" ",
				sizeCol,
				"   ",
				styledName,
			),
		)
	}
}
