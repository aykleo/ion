package pager

import (
	"regexp"
	"strings"
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func displayWidth(text string) int {
	return len(ansiRegex.ReplaceAllString(text, ""))
}

func splitPreservingANSI(text string, displayPos int) (string, string) {
	if displayPos <= 0 {
		return "", text
	}

	runes := []rune(text)
	currentDisplayPos := 0
	actualPos := 0

	for actualPos < len(runes) {
		remaining := string(runes[actualPos:])
		if strings.HasPrefix(remaining, "\x1b[") {
			for actualPos < len(runes) && runes[actualPos] != 'm' {
				actualPos++
			}
			if actualPos < len(runes) {
				actualPos++
			}
			continue
		}

		if currentDisplayPos >= displayPos {
			break
		}
		currentDisplayPos++
		actualPos++
	}

	if actualPos >= len(runes) {
		return text, ""
	}

	return string(runes[:actualPos]), string(runes[actualPos:])
}

func wrapText(text string, width int) string {
	if width <= 0 {
		return text
	}

	lines := strings.Split(text, "\n")
	var wrappedLines []string

	for _, line := range lines {
		if displayWidth(line) <= width {
			wrappedLines = append(wrappedLines, line)
			continue
		}

		if strings.TrimSpace(line) == "" {
			wrappedLines = append(wrappedLines, line)
			continue
		}

		currentLine := ""
		remaining := line

		for displayWidth(remaining) > 0 {
			if displayWidth(currentLine)+displayWidth(remaining) <= width {
				if currentLine == "" {
					currentLine = remaining
				} else {
					currentLine += remaining
				}
				break
			}

			availableWidth := width - displayWidth(currentLine)
			if availableWidth <= 0 {
				if currentLine != "" {
					wrappedLines = append(wrappedLines, currentLine)
					currentLine = ""
				}
				availableWidth = width
			}

			bestBreakPoint := availableWidth
			words := strings.Fields(remaining)
			if len(words) > 0 {
				wordStart := 0
				for _, word := range words {
					wordEnd := wordStart + displayWidth(word)
					if wordEnd <= availableWidth {
						wordStart = wordEnd + 1
						continue
					} else {
						if wordStart > 0 {
							bestBreakPoint = wordStart - 1
						} else if displayWidth(word) > availableWidth {
							bestBreakPoint = availableWidth
						}
						break
					}
				}
			}

			part, rest := splitPreservingANSI(remaining, bestBreakPoint)
			if part == "" {
				part, rest = splitPreservingANSI(remaining, 1)
			}

			if currentLine == "" {
				currentLine = part
			} else {
				currentLine += part
			}

			if displayWidth(currentLine) >= width || strings.HasSuffix(part, " ") {
				wrappedLines = append(wrappedLines, strings.TrimRight(currentLine, " "))
				currentLine = ""
			}

			remaining = strings.TrimLeft(rest, " ")
		}

		if currentLine != "" {
			wrappedLines = append(wrappedLines, currentLine)
		}
	}

	return strings.Join(wrappedLines, "\n")
}

func (m *Pager) wrapContent() string {
	if len(m.content) == 0 {
		return ""
	}

	effectiveWidth := m.viewport.Width - 4
	if effectiveWidth <= 0 {
		effectiveWidth = 80
	}

	var wrappedContent []string
	for _, item := range m.content {
		if strings.TrimSpace(item) == "" {
			wrappedContent = append(wrappedContent, "")
			continue
		}
		wrappedItem := wrapText(item, effectiveWidth)
		wrappedContent = append(wrappedContent, wrappedItem)
	}
	return strings.Join(wrappedContent, "\n\n")
}
