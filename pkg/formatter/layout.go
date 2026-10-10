package formatter

import (
	"regexp"
	"strings"
)

var stripAnsiRegex = regexp.MustCompile(`\x1b\][^\x07\x1b]*[\x07\x1b\\]|\x1b\[[0-9;]*[a-zA-Z]`)

// VisibleWidth returns the visible terminal column count of a string, ignoring ANSI escape codes.
func VisibleWidth(s string) int {
	clean := stripAnsiRegex.ReplaceAllString(s, "")
	width := 0
	for _, r := range clean {
		width += runeWidth(r)
	}
	return width
}

// runeWidth estimates terminal display column width for a rune.
func runeWidth(r rune) int {
	// Control characters or zero-width
	if r < 32 || (r >= 0x7F && r < 0xA0) {
		return 0
	}
	// Wide CJK characters and wide surrogate emojis occupy 2 cells
	if (r >= 0x1F300 && r <= 0x1F9FF) || (r >= 0x2E80 && r <= 0x9FFF) {
		return 2
	}
	// Standard symbols (✔, ⚠, ✖, ℹ, •) and alphanumeric runes occupy 1 cell
	return 1
}

// WrapText wraps text into lines of at most maxLen visible characters.
// Handles long unbroken strings (e.g. URLs, ARNs) without breaking border alignment.
func WrapText(text string, maxLen int) []string {
	if maxLen <= 10 {
		maxLen = 10
	}

	lines := strings.Split(text, "\n")
	var result []string

	for _, line := range lines {
		if VisibleWidth(line) <= maxLen {
			result = append(result, line)
			continue
		}

		words := strings.Fields(line)
		if len(words) == 0 {
			result = append(result, "")
			continue
		}

		// Expand words that individually exceed maxLen into chunked words
		var processedWords []string
		for _, w := range words {
			if VisibleWidth(w) > maxLen {
				processedWords = append(processedWords, chunkWord(w, maxLen)...)
			} else {
				processedWords = append(processedWords, w)
			}
		}

		var currentLine strings.Builder
		currentWidth := 0

		for _, word := range processedWords {
			wordWidth := VisibleWidth(word)

			if currentWidth == 0 {
				currentLine.WriteString(word)
				currentWidth = wordWidth
			} else if currentWidth+1+wordWidth <= maxLen {
				currentLine.WriteString(" ")
				currentLine.WriteString(word)
				currentWidth += 1 + wordWidth
			} else {
				result = append(result, currentLine.String())
				currentLine.Reset()
				currentLine.WriteString(word)
				currentWidth = wordWidth
			}
		}

		if currentLine.Len() > 0 {
			result = append(result, currentLine.String())
		}
	}

	return result
}

// chunkWord splits an oversized word into sub-chunks of at most maxLen visible characters.
func chunkWord(word string, maxLen int) []string {
	var chunks []string
	runes := []rune(word)

	for len(runes) > 0 {
		end := maxLen
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[:end]))
		runes = runes[end:]
	}

	return chunks
}

// PadRight pads a string with spaces up to targetWidth visible characters.
func PadRight(s string, targetWidth int) string {
	vw := VisibleWidth(s)
	if vw >= targetWidth {
		return s
	}
	return s + strings.Repeat(" ", targetWidth-vw)
}
