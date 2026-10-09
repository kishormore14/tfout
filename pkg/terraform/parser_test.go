package terraform_test

import (
	"strings"
	"testing"

	"github.com/tfout/tfout/pkg/formatter"
	"github.com/tfout/tfout/pkg/terraform"
)

func TestParseOutputs(t *testing.T) {
	sampleJSON := `{
		"domain_name": {
			"sensitive": false,
			"type": "string",
			"value": "example.com"
		},
		"db_password": {
			"sensitive": true,
			"type": "string",
			"value": "SuperSecret123!"
		},
		"name_servers": {
			"sensitive": false,
			"type": ["list", "string"],
			"value": ["ns1.example.net", "ns2.example.net"]
		},
		"instance_count": {
			"sensitive": false,
			"type": "number",
			"value": 3
		},
		"enabled": {
			"sensitive": false,
			"type": "bool",
			"value": true
		}
	}`

	outputs, err := terraform.ParseOutputs([]byte(sampleJSON))
	if err != nil {
		t.Fatalf("ParseOutputs failed: %v", err)
	}

	if len(outputs) != 5 {
		t.Errorf("Expected 5 outputs, got %d", len(outputs))
	}

	// Test sensitive redaction default
	val, isSens := terraform.FormatValue(outputs["db_password"].Value, outputs["db_password"].Sensitive, false)
	if val != "[SENSITIVE REDACTED]" {
		t.Errorf("Expected redacted sensitive value, got %q", val)
	}
	if !isSens {
		t.Error("Expected sensitive flag to be true")
	}

	// Test sensitive revealed option
	valRevealed, _ := terraform.FormatValue(outputs["db_password"].Value, outputs["db_password"].Sensitive, true)
	if valRevealed != "SuperSecret123!" {
		t.Errorf("Expected revealed sensitive value 'SuperSecret123!', got %q", valRevealed)
	}

	// Test banner conversion
	cfg := formatter.Config{
		ColorMode:     formatter.ColorNever,
		Theme:         formatter.ThemeMono,
		BorderStyle:   formatter.BorderSingle,
		ShowSensitive: false,
		TermWidth:     80,
	}

	bannerData := outputs.ConvertToBanner(".", cfg)
	if bannerData.Title != "TERRAFORM OUTPUTS" {
		t.Errorf("Expected title 'TERRAFORM OUTPUTS', got %q", bannerData.Title)
	}

	rendered := formatter.RenderBanner(bannerData, cfg)
	if !strings.Contains(rendered, "db_password (redacted)") {
		t.Error("Rendered output missing redacted field label")
	}
	if !strings.Contains(rendered, "ns1.example.net") {
		t.Error("Rendered output missing list item")
	}
}

func TestParseEmptyOutputs(t *testing.T) {
	_, err := terraform.ParseOutputs([]byte(""))
	if err == nil {
		t.Error("Expected error on empty input, got nil")
	}
}
