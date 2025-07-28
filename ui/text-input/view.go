package textinput

import (
	"strings"
	"unicode"

	"github.com/aykleo/ion/ui/styles"
)

func (m Input) View() string {
	value := m.input.Value()

	if len(value) == 0 {
		return m.input.View()
	}

	originalView := m.input.View()

	firstWordEnd := findFirstWordEnd(value)

	if firstWordEnd <= 0 {
		return originalView
	}

	prompt := m.input.PromptStyle.Render(m.input.Prompt)
	firstWord := value[:firstWordEnd]

	styledFirstWord := m.FirstWordStyle.Render(firstWord)

	viewWithoutPrompt := strings.TrimPrefix(originalView, prompt)

	if firstWordEnd >= len(value) {
		modifiedView := strings.Replace(viewWithoutPrompt, firstWord, styledFirstWord, 1)
		return prompt + modifiedView
	}

	remaining := value[firstWordEnd:]
	styledRemaining := styles.NoStyle.Render(remaining)

	textPart := firstWord + remaining
	styledTextPart := styledFirstWord + styledRemaining

	modifiedView := strings.Replace(viewWithoutPrompt, textPart, styledTextPart, 1)
	return prompt + modifiedView
}

func findFirstWordEnd(text string) int {
	if len(text) == 0 {
		return 0
	}

	for i, r := range text {
		if unicode.IsSpace(r) {
			return i
		}
	}

	return len(text)
}
