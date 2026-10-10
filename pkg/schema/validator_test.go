package schema_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/tfout/tfout/pkg/schema"
)

func TestParseAndValidate(t *testing.T) {
	validJSON := `{
		"schema_version": "1.0",
		"title": "Database Instance Created",
		"status": "success",
		"message": "RDS instance deployed successfully.",
		"details": {
			"Endpoint": "db.example.internal",
			"Port": "5432"
		},
		"items": [
			"Automatic backups enabled",
			"Multi-AZ standby ready"
		]
	}`

	doc, err := schema.ParseAndValidate(strings.NewReader(validJSON))
	if err != nil {
		t.Fatalf("ParseAndValidate failed unexpectedly: %v", err)
	}

	if doc.Title != "Database Instance Created" {
		t.Errorf("Expected title 'Database Instance Created', got %q", doc.Title)
	}
	if doc.Status != "success" {
		t.Errorf("Expected status 'success', got %q", doc.Status)
	}
	if len(doc.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(doc.Items))
	}
}

func TestParseAndValidateErrors(t *testing.T) {
	tests := []struct {
		name        string
		jsonInput   string
		errContains string
	}{
		{
			name:        "empty input",
			jsonInput:   "",
			errContains: "JSON input is empty",
		},
		{
			name:        "invalid JSON syntax",
			jsonInput:   `{invalid}`,
			errContains: "invalid JSON syntax",
		},
		{
			name:        "missing schema_version",
			jsonInput:   `{"title": "Test"}`,
			errContains: "missing required field 'schema_version'",
		},
		{
			name:        "unsupported schema_version",
			jsonInput:   `{"schema_version": "99.0", "title": "Test"}`,
			errContains: "unsupported schema_version",
		},
		{
			name:        "missing title",
			jsonInput:   `{"schema_version": "1.0"}`,
			errContains: "missing or empty required field 'title'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := schema.ParseAndValidate(strings.NewReader(tt.jsonInput))
			if err == nil {
				t.Fatalf("Expected error containing %q, got nil", tt.errContains)
			}
			if !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("Error %q does not contain %q", err.Error(), tt.errContains)
			}
		})
	}
}

func TestParseAndValidateOversizedInput(t *testing.T) {
	oversized := make([]byte, schema.MaxInputSize+100)
	for i := range oversized {
		oversized[i] = ' '
	}
	_, err := schema.ParseAndValidate(bytes.NewReader(oversized))
	if err == nil {
		t.Fatal("Expected error for oversized input, got nil")
	}
	if !strings.Contains(err.Error(), "exceeds maximum allowed size") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestDeterministicDetailsOrdering(t *testing.T) {
	jsonInput := `{
		"schema_version": "1.0",
		"title": "Test",
		"details": {
			"Zebra": "1",
			"Alpha": "2",
			"Beta": "3"
		}
	}`
	doc, err := schema.ParseAndValidate(strings.NewReader(jsonInput))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	details := doc.OrderedDetails()
	if len(details) != 3 {
		t.Fatalf("Expected 3 details, got %d", len(details))
	}
	if details[0].Key != "Alpha" || details[1].Key != "Beta" || details[2].Key != "Zebra" {
		t.Errorf("Details not sorted deterministically: %+v", details)
	}
}
