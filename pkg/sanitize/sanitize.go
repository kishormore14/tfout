package sanitize

import (
	"regexp"
	"strings"
	"unicode"
)

// ansiDangerousRegex matches non-SGR ANSI escape sequences (OSC title hijacks, cursor movements, screen clear, etc.)
// Preserves safe SGR text styling/color codes (ending in 'm', e.g. \x1b[32m for green).
var ansiDangerousRegex = regexp.MustCompile(`\x1b\][^\x07\x1b]*[\x07\x1b\\]|\x1b\[[0-9;]*[a-zA-LN-Z]`)

// Text strips dangerous terminal control characters from a string while preserving safe ANSI color codes.
// If allowNewline is false, newlines (\r, \n) are replaced with spaces.
func Text(s string, allowNewline bool) string {
	// Strip dangerous ANSI control sequences (keep SGR color codes ending in 'm')
	cleaned := ansiDangerousRegex.ReplaceAllString(s, "")

	var sb strings.Builder
	sb.Grow(len(cleaned))

	for _, r := range cleaned {
		if r == '\n' || r == '\r' {
			if allowNewline {
				sb.WriteRune('\n')
			} else {
				sb.WriteRune(' ')
			}
			continue
		}
		if r == '\t' {
			sb.WriteString("    ") // Expand tab to 4 spaces for consistent width calculation
			continue
		}
		// Strip ASCII control chars (0x00-0x1F, 0x7F except ESC 0x1B) and unicode control runes
		if unicode.IsControl(r) && r != '\x1b' {
			continue
		}
		sb.WriteRune(r)
	}

	res := sb.String()
	if !allowNewline {
		// Collapse multiple consecutive spaces created by newline replacement
		res = collapseSpaces(res)
	}
	return strings.TrimSpace(res)
}

// SingleLine cleans a string for single-line usage like headers, titles, or keys.
func SingleLine(s string) string {
	return Text(s, false)
}

// Multiline cleans a multiline string value, preserving newlines while stripping control sequences.
func Multiline(s string) string {
	lines := strings.Split(s, "\n")
	var cleanedLines []string
	for _, line := range lines {
		cleaned := Text(line, false)
		cleanedLines = append(cleanedLines, cleaned)
	}
	return strings.Join(cleanedLines, "\n")
}

func collapseSpaces(s string) string {
	var sb strings.Builder
	inSpace := false
	for _, r := range s {
		if r == ' ' {
			if !inSpace {
				sb.WriteRune(' ')
				inSpace = true
			}
		} else {
			sb.WriteRune(r)
			inSpace = false
		}
	}
	return sb.String()
}
