package sanitize_test

import (
	"testing"

	"github.com/tfout/tfout/pkg/sanitize"
)

func TestSingleLineSanitize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "plain text",
			input:    "Route 53 Zone",
			expected: "Route 53 Zone",
		},
		{
			name:     "Green SGR color code preservation",
			input:    "\x1b[32malb-123456.us-east-1.elb.amazonaws.com\x1b[0m",
			expected: "\x1b[32malb-123456.us-east-1.elb.amazonaws.com\x1b[0m",
		},
		{
			name:     "OSC title injection stripping",
			input:    "\x1b]0;Title Hijack\x07Malicious Header",
			expected: "Malicious Header",
		},
		{
			name:     "newline and control characters in single line",
			input:    "Header\nSubHeader\r\x07Bell",
			expected: "Header SubHeader Bell",
		},
		{
			name:     "tab expansion",
			input:    "Key:\tValue",
			expected: "Key: Value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitize.SingleLine(tt.input)
			if got != tt.expected {
				t.Errorf("SingleLine(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestMultilineSanitize(t *testing.T) {
	input := "\x1b[32mLine 1\x1b[0m\nLine 2\r\x07With Bell\nLine 3"
	expected := "\x1b[32mLine 1\x1b[0m\nLine 2 With Bell\nLine 3"

	got := sanitize.Multiline(input)
	if got != expected {
		t.Errorf("Multiline(%q) = %q; want %q", input, got, expected)
	}
}
