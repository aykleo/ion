package styles

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"github.com/aykleo/ion/data"
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
		if strings.Contains(line, "Name") && strings.Contains(line, "Value") && strings.Contains(line, "Salt") && strings.Contains(line, "Tags") && strings.Contains(line, "Updated") {
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

func isJsonOutput(output string) bool {
	trimmed := strings.TrimSpace(output)
	return strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]")
}

func formatJsonOutput(output string) string {
	lines := strings.Split(output, "\n")
	var formattedLines []string

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			formattedLines = append(formattedLines, "")
			continue
		}

		var result strings.Builder
		inString := false
		isProperty := false
		escaped := false

		for i, char := range line {
			switch char {
			case '"':
				if !escaped {
					if !inString {
						inString = true
						rest := line[i+1:]
						nextQuoteIndex := strings.Index(rest, "\"")
						if nextQuoteIndex != -1 {
							afterQuote := rest[nextQuoteIndex+1:]
							colonIndex := strings.Index(strings.TrimSpace(afterQuote), ":")
							if colonIndex == 0 || (colonIndex > 0 && strings.TrimSpace(afterQuote[:colonIndex]) == "") {
								isProperty = true
							}
						}
					} else {
						inString = false
						isProperty = false
					}
				}

				if isProperty {
					result.WriteString(MainTheme.Render(string(char)))
				} else {
					result.WriteString(NoStyle.Render(string(char)))
				}
				escaped = false
			case '[', ']':
				result.WriteString(MiscStyle.Render(string(char)))
				escaped = false
			case '{', '}':
				result.WriteString(MiscStyle.Render(string(char)))
				escaped = false
			case '\\':
				if inString && isProperty {
					result.WriteString(MainTheme.Render(string(char)))
				} else {
					result.WriteString(NoStyle.Render(string(char)))
				}
				escaped = !escaped
			default:
				if inString && isProperty {
					result.WriteString(MainTheme.Render(string(char)))
				} else {
					result.WriteString(NoStyle.Render(string(char)))
				}
				escaped = false
			}
		}

		formattedLines = append(formattedLines, result.String())
	}

	return strings.Join(formattedLines, "\n")
}

func formatTable(output string) string {
	lines := strings.Split(strings.TrimSpace(output), "\n")

	for _, line := range lines {
		if strings.Contains(line, "Name") && strings.Contains(line, "Value") && strings.Contains(line, "Salt") && strings.Contains(line, "Tags") && strings.Contains(line, "Updated") {
			return formatSecretListInTable(output)
		}
	}

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

func formatSecretRow(line string) string {
	if len(line) < 60 {
		return TableRowStyle.Render(line)
	}

	name := strings.TrimSpace(line[0:17])
	value := strings.TrimSpace(line[17:31])
	salt := strings.TrimSpace(line[31:44])
	tags := strings.TrimSpace(line[44:67])
	updated := strings.TrimSpace(line[67:])

	nameWithEmoji := "🔑 " + name
	styledName := MainTheme.Bold(true).Render(nameWithEmoji)

	nameCol := lipgloss.NewStyle().Width(17).Render(styledName)
	valueCol := lipgloss.NewStyle().Width(14).Render(value)
	saltCol := lipgloss.NewStyle().Width(15).Render(salt)
	tagsCol := lipgloss.NewStyle().Width(25).Render(tags)
	updatedCol := lipgloss.NewStyle().Render(updated)

	return TableRowStyle.Render(
		lipgloss.JoinHorizontal(lipgloss.Left,
			nameCol, " ",
			valueCol, " ",
			saltCol, " ",
			tagsCol, " ",
			updatedCol,
		),
	)
}

func formatSecretListInTable(output string) string {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var formattedLines []string
	var headerFound bool

	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			if len(formattedLines) > 0 {
				formattedLines = append(formattedLines, "")
			}
			continue
		}

		if strings.Contains(line, "Name") && strings.Contains(line, "Value") && strings.Contains(line, "Salt") && strings.Contains(line, "Tags") && strings.Contains(line, "Updated") {
			headerFound = true
			nameHeader := lipgloss.NewStyle().Width(17).Render("🔐 Name")
			valueHeader := lipgloss.NewStyle().Width(14).Render("Value")
			saltHeader := lipgloss.NewStyle().Width(15).Render("Salt")
			tagsHeader := lipgloss.NewStyle().Width(25).Render("Tags")
			updatedHeader := lipgloss.NewStyle().Render("Updated")

			headerLine := lipgloss.JoinHorizontal(lipgloss.Left,
				nameHeader, " ",
				valueHeader, " ",
				saltHeader, " ",
				tagsHeader, " ",
				updatedHeader,
			)

			styled := TableHeaderStyle.Render(headerLine)
			formattedLines = append(formattedLines, styled)
			continue
		}

		if regexp.MustCompile(`^-+\s+-+`).MatchString(trimmedLine) {
			separator := TableBorderStyle.Render(strings.Repeat("─", 90))
			formattedLines = append(formattedLines, separator)
			continue
		}

		if headerFound && i > 1 && trimmedLine != "" {
			formatted := formatSecretRow(line)
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

func FormatSecretsAsJSON(secrets []data.Secret) string {
	type JSONSecret struct {
		Name      string    `json:"name"`
		Value     string    `json:"value"`
		Salt      string    `json:"salt"`
		Tags      []string  `json:"tags"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	var jsonSecrets []JSONSecret
	for _, secret := range secrets {
		jsonSecrets = append(jsonSecrets, JSONSecret{
			Name:      secret.ID,
			Value:     secret.Value,
			Salt:      secret.Salt,
			Tags:      secret.Tags,
			CreatedAt: secret.CreatedAt,
			UpdatedAt: secret.UpdatedAt,
		})
	}

	jsonData, err := json.MarshalIndent(jsonSecrets, "", "  ")
	if err != nil {
		return "Error formatting JSON: " + err.Error()
	}

	return string(jsonData)
}
