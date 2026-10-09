package terraform

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/tfout/tfout/pkg/formatter"
	"github.com/tfout/tfout/pkg/sanitize"
)

// OutputValue represents a single output from `terraform output -json`.
type OutputValue struct {
	Type      interface{} `json:"type"`
	Value     interface{} `json:"value"`
	Sensitive bool        `json:"sensitive"`
}

// OutputsMap is a map of output name to OutputValue.
type OutputsMap map[string]OutputValue

// ParseOutputs parses raw JSON from `terraform output -json`.
func ParseOutputs(jsonData []byte) (OutputsMap, error) {
	if len(strings.TrimSpace(string(jsonData))) == 0 {
		return nil, fmt.Errorf("terraform outputs JSON is empty")
	}

	var outputs OutputsMap
	if err := json.Unmarshal(jsonData, &outputs); err != nil {
		return nil, fmt.Errorf("failed to parse terraform output JSON: %w", err)
	}

	return outputs, nil
}

// FormatValue formats a raw interface{} value into a clean, human-readable string.
func FormatValue(val interface{}, sensitive bool, showSensitive bool) (string, bool) {
	if sensitive && !showSensitive {
		return "[SENSITIVE REDACTED]", true
	}

	if val == nil {
		return "null", sensitive
	}

	switch v := val.(type) {
	case string:
		return sanitize.Multiline(v), sensitive
	case bool:
		return fmt.Sprintf("%t", v), sensitive
	case float64:
		// Check if it's an integer within int64 bounds
		if v >= float64(math.MinInt64) && v <= float64(math.MaxInt64) && v == float64(int64(v)) {
			return fmt.Sprintf("%d", int64(v)), sensitive
		}
		return fmt.Sprintf("%g", v), sensitive
	default:
		// For maps, slices, complex objects: pretty-print JSON
		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return fmt.Sprintf("%v", v), sensitive
		}
		return sanitize.Multiline(string(b)), sensitive
	}
}

// ConvertToBanner converts an OutputsMap to formatter.BannerData.
func (om OutputsMap) ConvertToBanner(dir string, cfg formatter.Config) formatter.BannerData {
	var details []formatter.DetailPair
	hasRevealedSensitive := false
	sensitiveCount := 0

	// Sort keys alphabetically for predictable output
	keys := make([]string, 0, len(om))
	for k := range om {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		output := om[k]
		valStr, _ := FormatValue(output.Value, output.Sensitive, cfg.ShowSensitive)

		if output.Sensitive {
			sensitiveCount++
			if cfg.ShowSensitive {
				hasRevealedSensitive = true
			}
		}

		cleanKey := sanitize.SingleLine(k)
		if output.Sensitive {
			if cfg.ShowSensitive {
				cleanKey = cleanKey + " (sensitive)"
			} else {
				cleanKey = cleanKey + " (redacted)"
			}
		}

		details = append(details, formatter.DetailPair{
			Key:   cleanKey,
			Value: valStr,
		})
	}

	msg := fmt.Sprintf("Terraform Directory: %s | Total Outputs: %d", dir, len(om))
	if sensitiveCount > 0 && !cfg.ShowSensitive {
		msg += fmt.Sprintf(" (%d sensitive redacted)", sensitiveCount)
	}

	title := "TERRAFORM OUTPUTS"
	if len(om) == 0 {
		title = "NO TERRAFORM OUTPUTS FOUND"
		msg = fmt.Sprintf("Terraform Directory: %s | No output values declared or populated.", dir)
	}

	return formatter.BannerData{
		Title:         title,
		Status:        formatter.StatusSuccess,
		Message:       msg,
		Details:       details,
		SensitiveWarn: hasRevealedSensitive,
	}
}
