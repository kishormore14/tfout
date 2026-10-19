package formatter_test

import (
	"os"
	"strings"
	"testing"

	"github.com/tfout/tfout/pkg/formatter"
)

func TestShouldUseColor(t *testing.T) {
	// Isolate from any force-color variables in the developer/CI environment
	for _, key := range formatter.ForceColorEnvKeys {
		os.Unsetenv(key)
	}

	// Mono theme should never use color
	if formatter.ShouldUseColor(formatter.ColorAlways, formatter.ThemeMono) {
		t.Error("ShouldUseColor should return false for ThemeMono")
	}

	// ColorNever mode should never use color
	if formatter.ShouldUseColor(formatter.ColorNever, formatter.ThemeDefault) {
		t.Error("ShouldUseColor should return false for ColorNever")
	}

	// ColorAlways mode should return true
	if !formatter.ShouldUseColor(formatter.ColorAlways, formatter.ThemeDefault) {
		t.Error("ShouldUseColor should return true for ColorAlways")
	}

	// NO_COLOR test
	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	if formatter.ShouldUseColor(formatter.ColorAuto, formatter.ThemeDefault) {
		t.Error("ShouldUseColor should return false when NO_COLOR is set")
	}
}

func TestForceColorEnvOverrides(t *testing.T) {
	for _, key := range formatter.ForceColorEnvKeys {
		os.Unsetenv(key)
	}
	os.Unsetenv("NO_COLOR")

	// No force-color variable present
	if _, forced := formatter.ForceColorOverride(); forced {
		t.Error("ForceColorOverride should not force when no force-color variable is set")
	}

	for _, key := range formatter.ForceColorEnvKeys {
		t.Run(key+"_forces_on", func(t *testing.T) {
			os.Setenv(key, "1")
			defer os.Unsetenv(key)

			enabled, forced := formatter.ForceColorOverride()
			if !forced || !enabled {
				t.Errorf("%s=1 should force colors ON, got enabled=%v forced=%v", key, enabled, forced)
			}
			// Piped/CI stdout (not a TTY) must still get colors in auto mode
			if !formatter.ShouldUseColor(formatter.ColorAuto, formatter.ThemeDefault) {
				t.Errorf("ShouldUseColor(auto) should be true when %s=1", key)
			}
			// Explicit opt-in must win over NO_COLOR
			os.Setenv("NO_COLOR", "1")
			defer os.Unsetenv("NO_COLOR")
			if !formatter.ShouldUseColor(formatter.ColorAuto, formatter.ThemeDefault) {
				t.Errorf("%s=1 should override NO_COLOR", key)
			}
			// --color never and mono theme still win
			if formatter.ShouldUseColor(formatter.ColorNever, formatter.ThemeDefault) {
				t.Error("ColorNever should still disable color when force env is set")
			}
			if formatter.ShouldUseColor(formatter.ColorAuto, formatter.ThemeMono) {
				t.Error("Mono theme should still disable color when force env is set")
			}
		})

		t.Run(key+"_forces_off", func(t *testing.T) {
			os.Setenv(key, "0")
			defer os.Unsetenv(key)

			enabled, forced := formatter.ForceColorOverride()
			if !forced || enabled {
				t.Errorf("%s=0 should force colors OFF, got enabled=%v forced=%v", key, enabled, forced)
			}
			if formatter.ShouldUseColor(formatter.ColorAuto, formatter.ThemeDefault) {
				t.Errorf("ShouldUseColor(auto) should be false when %s=0", key)
			}
			// Explicit --color always still wins
			if !formatter.ShouldUseColor(formatter.ColorAlways, formatter.ThemeDefault) {
				t.Error("ColorAlways should still enable color when force env forces off")
			}
		})
	}
}

func TestVisibleWidthAndWrapText(t *testing.T) {
	styledStr := "\x1b[31mHello\x1b[0m World"
	if formatter.VisibleWidth(styledStr) != 11 {
		t.Errorf("Expected visible width 11, got %d", formatter.VisibleWidth(styledStr))
	}

	longText := "This is a long line of text that needs word wrapping to fit cleanly inside borders."
	wrapped := formatter.WrapText(longText, 20)
	for _, line := range wrapped {
		if formatter.VisibleWidth(line) > 20 {
			t.Errorf("Wrapped line %q exceeds max width 20", line)
		}
	}
}

func TestWrapTextOversizedSingleWord(t *testing.T) {
	oversizedWord := "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/k8s-default-ingress-tg/1234567890abcdef"
	wrapped := formatter.WrapText(oversizedWord, 25)
	for _, line := range wrapped {
		if formatter.VisibleWidth(line) > 25 {
			t.Errorf("Chunked line %q exceeds max width 25", line)
		}
	}
}

func TestRenderBanner(t *testing.T) {
	data := formatter.BannerData{
		Title:   "ROUTE53 HOSTED ZONE CREATED",
		Status:  formatter.StatusSuccess,
		Message: "Domain zone created successfully",
		Details: []formatter.DetailPair{
			{Key: "Domain", Value: "example.com"},
			{Key: "Name Servers", Value: "ns1.example.net\nns2.example.net"},
		},
		Items: []string{"Item 1", "Item 2"},
	}

	cfg := formatter.Config{
		ColorMode:     formatter.ColorNever,
		Theme:         formatter.ThemeMono,
		BorderStyle:   formatter.BorderSingle,
		ShowSensitive: false,
		TermWidth:     80,
	}

	out := formatter.RenderBanner(data, cfg)

	if !strings.Contains(out, "ROUTE53 HOSTED ZONE CREATED") {
		t.Error("Output missing title")
	}
	if !strings.Contains(out, "[SUCCESS]") {
		t.Error("Output missing success badge")
	}
	if !strings.Contains(out, "Domain:") || !strings.Contains(out, "example.com") {
		t.Error("Output missing details key/value")
	}
	if !strings.Contains(out, "ns1.example.net") || !strings.Contains(out, "ns2.example.net") {
		t.Error("Output missing multiline detail values")
	}
	if !strings.Contains(out, "Item 1") {
		t.Error("Output missing item")
	}
}

func TestBorderStyles(t *testing.T) {
	single := formatter.GetBorderChars(formatter.BorderSingle, formatter.ThemeDefault)
	if single.TopLeft != "┌" {
		t.Errorf("Expected ┌ for single border TopLeft, got %s", single.TopLeft)
	}

	double := formatter.GetBorderChars(formatter.BorderDouble, formatter.ThemeDefault)
	if double.TopLeft != "╔" {
		t.Errorf("Expected ╔ for double border TopLeft, got %s", double.TopLeft)
	}

	monoSingle := formatter.GetBorderChars(formatter.BorderSingle, formatter.ThemeMono)
	if monoSingle.TopLeft != "+" {
		t.Errorf("Expected + for mono border TopLeft, got %s", monoSingle.TopLeft)
	}
}
