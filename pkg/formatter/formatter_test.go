package formatter_test

import (
	"os"
	"strings"
	"testing"

	"github.com/tfout/tfout/pkg/formatter"
)

func TestShouldUseColor(t *testing.T) {
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
