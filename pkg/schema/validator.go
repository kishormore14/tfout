package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/tfout/tfout/pkg/sanitize"
)

// SupportedSchemaVersion is the current version supported by tfout.
const SupportedSchemaVersion = "1.0"

// MaxInputSize is 10 MB limit for JSON inputs to prevent memory exhaustion.
const MaxInputSize = 10 * 1024 * 1024

// Detail represents a key-value detail entry.
type Detail struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// RenderDocument represents the JSON structure for `tfout render`.
type RenderDocument struct {
	SchemaVersion string                 `json:"schema_version"`
	Title         string                 `json:"title"`
	Status        string                 `json:"status,omitempty"`
	Message       string                 `json:"message,omitempty"`
	DetailsMap    map[string]interface{} `json:"details,omitempty"`
	DetailsList   []Detail               `json:"details_list,omitempty"`
	Items         []string               `json:"items,omitempty"`
	Theme         string                 `json:"theme,omitempty"`
	BorderStyle   string                 `json:"border_style,omitempty"`
}

// OrderedDetails returns key-value pairs extracted from either DetailsList or DetailsMap.
func (doc *RenderDocument) OrderedDetails() []Detail {
	if len(doc.DetailsList) > 0 {
		var result []Detail
		for _, d := range doc.DetailsList {
			result = append(result, Detail{
				Key:   sanitize.SingleLine(d.Key),
				Value: sanitize.Multiline(d.Value),
			})
		}
		return result
	}

	if len(doc.DetailsMap) > 0 {
		var result []Detail
		for k, v := range doc.DetailsMap {
			valStr := fmt.Sprintf("%v", v)
			result = append(result, Detail{
				Key:   sanitize.SingleLine(k),
				Value: sanitize.Multiline(valStr),
			})
		}
		return result
	}

	return nil
}

// ParseAndValidate parses JSON data from reader with size limit and validates schema.
func ParseAndValidate(r io.Reader) (*RenderDocument, error) {
	lr := io.LimitReader(r, MaxInputSize+1)
	data, err := io.ReadAll(lr)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON input: %w", err)
	}

	if len(data) > MaxInputSize {
		return nil, fmt.Errorf("JSON input exceeds maximum allowed size of %d bytes", MaxInputSize)
	}

	if len(strings.TrimSpace(string(data))) == 0 {
		return nil, errors.New("JSON input is empty")
	}

	var doc RenderDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("invalid JSON syntax: %w", err)
	}

	if doc.SchemaVersion == "" {
		return nil, errors.New("missing required field 'schema_version'")
	}

	if doc.SchemaVersion != SupportedSchemaVersion && doc.SchemaVersion != "1" {
		return nil, fmt.Errorf("unsupported schema_version %q (supported version: %s)", doc.SchemaVersion, SupportedSchemaVersion)
	}

	if strings.TrimSpace(doc.Title) == "" {
		return nil, errors.New("missing or empty required field 'title'")
	}

	// Sanitize string fields
	doc.Title = sanitize.SingleLine(doc.Title)
	if doc.Message != "" {
		doc.Message = sanitize.Multiline(doc.Message)
	}
	if doc.Status != "" {
		doc.Status = strings.ToLower(strings.TrimSpace(doc.Status))
	} else {
		doc.Status = "info"
	}

	for i, item := range doc.Items {
		doc.Items[i] = sanitize.Multiline(item)
	}

	return &doc, nil
}
