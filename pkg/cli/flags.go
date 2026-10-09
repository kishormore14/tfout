package cli

import (
	"flag"
	"fmt"
	"strings"

	"github.com/tfout/tfout/pkg/formatter"
)

// GlobalFlags holds common flags across tfout commands.
type GlobalFlags struct {
	Color         string
	Theme         string
	Border        string
	ShowSensitive bool
}

// AddGlobalFlags registers standard style flags on a FlagSet.
func AddGlobalFlags(fs *flag.FlagSet, gf *GlobalFlags) {
	fs.StringVar(&gf.Color, "color", "auto", "Color mode: auto, always, never")
	fs.StringVar(&gf.Theme, "theme", "default", "Theme: default, minimal, mono")
	fs.StringVar(&gf.Border, "border", "single", "Border style: single, double, none")
	fs.BoolVar(&gf.ShowSensitive, "show-sensitive", false, "Reveal sensitive values in output (displays warning)")
}

// ToFormatterConfig converts GlobalFlags into formatter.Config.
func (gf *GlobalFlags) ToFormatterConfig() (formatter.Config, error) {
	colorMode := formatter.ColorMode(strings.ToLower(strings.TrimSpace(gf.Color)))
	if colorMode != formatter.ColorAuto && colorMode != formatter.ColorAlways && colorMode != formatter.ColorNever {
		return formatter.Config{}, fmt.Errorf("invalid color mode %q (must be auto, always, or never)", gf.Color)
	}

	theme := formatter.Theme(strings.ToLower(strings.TrimSpace(gf.Theme)))
	if theme != formatter.ThemeDefault && theme != formatter.ThemeMinimal && theme != formatter.ThemeMono {
		return formatter.Config{}, fmt.Errorf("invalid theme %q (must be default, minimal, or mono)", gf.Theme)
	}

	borderStyle := formatter.BorderStyle(strings.ToLower(strings.TrimSpace(gf.Border)))
	if borderStyle != formatter.BorderSingle && borderStyle != formatter.BorderDouble && borderStyle != formatter.BorderNone {
		return formatter.Config{}, fmt.Errorf("invalid border style %q (must be single, double, or none)", gf.Border)
	}

	return formatter.Config{
		ColorMode:     colorMode,
		Theme:         theme,
		BorderStyle:   borderStyle,
		ShowSensitive: gf.ShowSensitive,
		TermWidth:     80,
	}, nil
}

// DetailFlag supports repeating --detail "Key" "Value" or --detail "Key=Value" flags.
type DetailFlag []formatter.DetailPair

func (d *DetailFlag) String() string {
	return fmt.Sprintf("%v", []formatter.DetailPair(*d))
}

func (d *DetailFlag) Set(value string) error {
	parts := strings.SplitN(value, "=", 2)
	if len(parts) == 2 {
		*d = append(*d, formatter.DetailPair{
			Key:   strings.TrimSpace(parts[0]),
			Value: strings.TrimSpace(parts[1]),
		})
		return nil
	}

	// Single string - treat as key with empty value or key waiting for value
	*d = append(*d, formatter.DetailPair{
		Key:   strings.TrimSpace(value),
		Value: "",
	})
	return nil
}

// StringSliceFlag supports repeating --item "val" flags.
type StringSliceFlag []string

func (s *StringSliceFlag) String() string {
	return strings.Join(*s, ", ")
}

func (s *StringSliceFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}
