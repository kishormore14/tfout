package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/tfout/tfout/pkg/cli"
)

func TestExecuteHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := cli.Execute([]string{"--help"}, nil, &stdout, &stderr)
	if code != 0 {
		t.Errorf("Expected exit code 0 for --help, got %d", code)
	}
	if !strings.Contains(stdout.String(), "tfout — A Reusable Open-Source Terminal Formatter") {
		t.Error("Expected help message in stdout")
	}
}

func TestExecuteVersion(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := cli.Execute([]string{"version"}, nil, &stdout, &stderr)
	if code != 0 {
		t.Errorf("Expected exit code 0 for version, got %d", code)
	}
	if !strings.Contains(stdout.String(), "tfout version") {
		t.Error("Expected version message in stdout")
	}
}

func TestExecuteDoctor(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := cli.Execute([]string{"doctor"}, nil, &stdout, &stderr)
	if code != 0 {
		t.Errorf("Expected exit code 0 for doctor, got %d", code)
	}
	if !strings.Contains(stdout.String(), "TFOUT SYSTEM DOCTOR") {
		t.Error("Expected doctor banner in stdout")
	}
}

func TestExecuteBanner(t *testing.T) {
	var stdout, stderr bytes.Buffer
	args := []string{
		"banner",
		"--title", "TEST BANNER",
		"--status", "success",
		"--detail", "Key=Value",
		"--color", "never",
		"--theme", "mono",
	}

	code := cli.Execute(args, nil, &stdout, &stderr)
	if code != 0 {
		t.Errorf("Expected exit code 0 for banner, got %d. Stderr: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "TEST BANNER") {
		t.Error("Expected title in banner output")
	}
}

func TestExecuteBannerMissingTitle(t *testing.T) {
	var stdout, stderr bytes.Buffer
	args := []string{"banner"}

	code := cli.Execute(args, nil, &stdout, &stderr)
	if code == 0 {
		t.Error("Expected non-zero exit code for banner without title")
	}
	if !strings.Contains(stderr.String(), "missing required option --title") {
		t.Error("Expected error message about missing title")
	}
}

func TestExecuteRenderStdin(t *testing.T) {
	jsonInput := `{
		"schema_version": "1.0",
		"title": "Render Test",
		"status": "info"
	}`

	stdin := strings.NewReader(jsonInput)
	var stdout, stderr bytes.Buffer

	code := cli.Execute([]string{"render", "--color", "never"}, stdin, &stdout, &stderr)
	if code != 0 {
		t.Errorf("Expected exit code 0 for render stdin, got %d. Stderr: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Render Test") {
		t.Error("Expected title in rendered output")
	}
}

func TestExecuteUnknownCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := cli.Execute([]string{"unknown_cmd"}, nil, &stdout, &stderr)
	if code == 0 {
		t.Error("Expected non-zero exit code for unknown command")
	}
	if !strings.Contains(stderr.String(), "Unknown command \"unknown_cmd\"") {
		t.Error("Expected unknown command error message")
	}
}
