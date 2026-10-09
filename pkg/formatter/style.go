package formatter

import (
	"os"
	"strings"
)

// ColorMode represents color enablement options.
type ColorMode string

const (
	ColorAuto   ColorMode = "auto"
	ColorAlways ColorMode = "always"
	ColorNever  ColorMode = "never"
)

// Theme represents the visual theme.
type Theme string

const (
	ThemeDefault Theme = "default"
	ThemeMinimal Theme = "minimal"
	ThemeMono    Theme = "mono"
)

// BorderStyle represents the box border style.
type BorderStyle string

const (
	BorderSingle BorderStyle = "single"
	BorderDouble BorderStyle = "double"
	BorderNone   BorderStyle = "none"
)

// Status represents execution status.
type Status string

const (
	StatusSuccess Status = "success"
	StatusWarning Status = "warning"
	StatusError   Status = "error"
	StatusInfo    Status = "info"
)

// Config configures formatting options.
type Config struct {
	ColorMode     ColorMode
	Theme         Theme
	BorderStyle   BorderStyle
	ShowSensitive bool
	TermWidth     int
}

// DefaultConfig returns reasonable default configuration.
func DefaultConfig() Config {
	return Config{
		ColorMode:     ColorAuto,
		Theme:         ThemeDefault,
		BorderStyle:   BorderSingle,
		ShowSensitive: false,
		TermWidth:     80,
	}
}

// ShouldUseColor evaluates whether ANSI colors should be emitted based on ColorMode, NO_COLOR, and theme.
func ShouldUseColor(mode ColorMode, theme Theme) bool {
	if theme == ThemeMono {
		return false
	}
	if mode == ColorNever {
		return false
	}
	if mode == ColorAlways {
		return true
	}
	// ColorAuto: check NO_COLOR environment variable
	if noColor := os.Getenv("NO_COLOR"); noColor != "" {
		return false
	}
	// Check if terminal stdout is a TTY
	return isTTY(os.Stdout)
}

// Helper to check if file descriptor is a terminal / TTY.
func isTTY(f *os.File) bool {
	if f == nil {
		return false
	}
	stat, err := f.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}

// ANSI Escape Codes
const (
	ANSIReset      = "\x1b[0m"
	ANSIBold       = "\x1b[1m"
	ANSIDim        = "\x1b[2m"
	ANSIFgGreen    = "\x1b[32m"
	ANSIFgYellow   = "\x1b[33m"
	ANSIFgRed      = "\x1b[31m"
	ANSIFgCyan     = "\x1b[36m"
	ANSIFgBlue     = "\x1b[34m"
	ANSIFgWhite    = "\x1b[37m"
	ANSIBgGreen    = "\x1b[42m"
	ANSIBgYellow   = "\x1b[43m"
	ANSIBgRed      = "\x1b[41m"
	ANSIBgCyan     = "\x1b[46m"
	ANSIFgBlack    = "\x1b[30m"
	ANSIBoldGreen  = "\x1b[1;32m"
	ANSIBoldYellow = "\x1b[1;33m"
	ANSIBoldRed    = "\x1b[1;31m"
	ANSIBoldCyan   = "\x1b[1;36m"
)

// StyleHelper formats text with optional ANSI styles.
type StyleHelper struct {
	UseColor bool
}

func NewStyleHelper(useColor bool) *StyleHelper {
	return &StyleHelper{UseColor: useColor}
}

func (s *StyleHelper) Style(text string, styleCode string) string {
	if !s.UseColor || styleCode == "" {
		return text
	}
	return styleCode + text + ANSIReset
}

func (s *StyleHelper) StatusBadge(status Status) string {
	normalized := Status(strings.ToLower(string(status)))
	switch normalized {
	case StatusSuccess:
		if s.UseColor {
			return s.Style(" ✔ SUCCESS ", ANSIBold + ANSIFgBlack + ANSIBgGreen)
		}
		return "[SUCCESS]"
	case StatusWarning:
		if s.UseColor {
			return s.Style(" ⚠ WARNING ", ANSIBold + ANSIFgBlack + ANSIBgYellow)
		}
		return "[WARNING]"
	case StatusError:
		if s.UseColor {
			return s.Style(" ✖ ERROR ", ANSIBold + ANSIFgWhite + ANSIBgRed)
		}
		return "[ERROR]"
	case StatusInfo:
		fallthrough
	default:
		if s.UseColor {
			return s.Style(" ℹ INFO ", ANSIBold + ANSIFgBlack + ANSIBgCyan)
		}
		return "[INFO]"
	}
}

func (s *StyleHelper) HeaderColor(status Status) string {
	switch Status(strings.ToLower(string(status))) {
	case StatusSuccess:
		return ANSIBoldGreen
	case StatusWarning:
		return ANSIBoldYellow
	case StatusError:
		return ANSIBoldRed
	case StatusInfo:
		fallthrough
	default:
		return ANSIBoldCyan
	}
}
