package formatter

import (
	"github.com/tfout/tfout/pkg/schema"
)

// RenderDocumentToBanner converts a schema.RenderDocument into BannerData and renders it.
func RenderDocumentToBanner(doc *schema.RenderDocument, cfg Config) string {
	// If document specifies theme/border_style and CLI did not explicitly override, apply document settings
	if doc.Theme != "" && cfg.Theme == ThemeDefault {
		cfg.Theme = Theme(doc.Theme)
	}
	if doc.BorderStyle != "" && cfg.BorderStyle == BorderSingle {
		cfg.BorderStyle = BorderStyle(doc.BorderStyle)
	}

	details := doc.OrderedDetails()
	var detailPairs []DetailPair
	for _, d := range details {
		detailPairs = append(detailPairs, DetailPair{
			Key:   d.Key,
			Value: d.Value,
		})
	}

	bannerData := BannerData{
		Title:         doc.Title,
		Status:        Status(doc.Status),
		Message:       doc.Message,
		Details:       detailPairs,
		Items:         doc.Items,
		SensitiveWarn: false,
	}

	return RenderBanner(bannerData, cfg)
}
