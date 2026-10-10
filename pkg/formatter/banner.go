package formatter

import (
	"fmt"
	"strings"

	"github.com/tfout/tfout/pkg/sanitize"
)

// DetailPair represents a key-value detail for the banner.
type DetailPair struct {
	Key   string
	Value string
}

// BannerData holds structured data to render a banner.
type BannerData struct {
	Title         string
	Status        Status
	Message       string
	Details       []DetailPair
	Items         []string
	SensitiveWarn bool
}

// RenderBanner formats BannerData into a string according to Config.
func RenderBanner(data BannerData, cfg Config) string {
	useColor := ShouldUseColor(cfg.ColorMode, cfg.Theme)
	sh := NewStyleHelper(useColor)
	border := GetBorderChars(cfg.BorderStyle, cfg.Theme)

	width := cfg.TermWidth
	if width < 50 {
		width = 50
	} else if width > 120 {
		width = 120
	}

	title := sanitize.SingleLine(data.Title)
	msg := sanitize.Multiline(data.Message)
	badge := sh.StatusBadge(data.Status)

	var sb strings.Builder

	// Border width calculation
	// Content inner width = width - 4 (2 border chars + 2 padding spaces)
	innerWidth := width - 4
	if cfg.BorderStyle == BorderNone {
		innerWidth = width
	}

	// Helper for top/bottom/separator border lines
	drawHorizontal := func(leftChar, midChar, rightChar string) string {
		if cfg.BorderStyle == BorderNone {
			return ""
		}
		line := leftChar + strings.Repeat(midChar, innerWidth+2) + rightChar
		return sh.Style(line, ANSIFgCyan)
	}

	// Helper for bordered lines
	wrapBorderLine := func(content string) string {
		if cfg.BorderStyle == BorderNone {
			return "  " + content
		}
		padded := PadRight(content, innerWidth)
		left := sh.Style(border.Vertical, ANSIFgCyan)
		right := sh.Style(border.Vertical, ANSIFgCyan)
		return left + " " + padded + " " + right
	}

	// Top Border
	if cfg.BorderStyle != BorderNone {
		sb.WriteString(drawHorizontal(border.TopLeft, border.Horizontal, border.TopRight))
		sb.WriteString("\n")
	}

	// Sensitive Warning Header if applicable
	if data.SensitiveWarn {
		warnText := sh.Style("WARNING: Sensitive values revealed in output", ANSIBoldRed)
		for _, line := range WrapText(warnText, innerWidth) {
			sb.WriteString(wrapBorderLine(line))
			sb.WriteString("\n")
		}
		if cfg.BorderStyle != BorderNone {
			sb.WriteString(drawHorizontal(border.TeeLeft, border.Horizontal, border.TeeRight))
			sb.WriteString("\n")
		}
	}

	// Header Line (Badge + Title)
	headerStr := badge + "  " + sh.Style(title, sh.HeaderColor(data.Status))
	for _, line := range WrapText(headerStr, innerWidth) {
		sb.WriteString(wrapBorderLine(line))
		sb.WriteString("\n")
	}

	// Message Section
	if msg != "" {
		if cfg.BorderStyle != BorderNone {
			sb.WriteString(drawHorizontal(border.TeeLeft, border.Horizontal, border.TeeRight))
			sb.WriteString("\n")
		} else {
			sb.WriteString("\n")
		}
		for _, line := range WrapText(msg, innerWidth) {
			sb.WriteString(wrapBorderLine(sh.Style(line, ANSIFgWhite)))
			sb.WriteString("\n")
		}
	}

	// Details Section (Key - Value pairs)
	if len(data.Details) > 0 {
		if cfg.BorderStyle != BorderNone {
			sb.WriteString(drawHorizontal(border.TeeLeft, border.Horizontal, border.TeeRight))
			sb.WriteString("\n")
		} else {
			sb.WriteString("\n")
		}

		maxKeyLen := 0
		for _, d := range data.Details {
			cleanKey := sanitize.SingleLine(d.Key)
			if len(cleanKey) > maxKeyLen {
				maxKeyLen = len(cleanKey)
			}
		}
		if maxKeyLen > 25 {
			maxKeyLen = 25
		}

		for _, d := range data.Details {
			cleanKey := sanitize.SingleLine(d.Key)
			cleanVal := sanitize.Multiline(d.Value)

			keyFormatted := sh.Style(PadRight(cleanKey+":", maxKeyLen+2), ANSIBold+ANSIFgCyan)
			valLines := strings.Split(cleanVal, "\n")

			for i, vLine := range valLines {
				var content string
				if i == 0 {
					content = fmt.Sprintf("%s %s", keyFormatted, sh.Style(vLine, ANSIFgWhite))
				} else {
					indent := strings.Repeat(" ", maxKeyLen+3)
					content = fmt.Sprintf("%s %s", indent, sh.Style(vLine, ANSIFgWhite))
				}

				for _, wrapped := range WrapText(content, innerWidth) {
					sb.WriteString(wrapBorderLine(wrapped))
					sb.WriteString("\n")
				}
			}
		}
	}

	// Items Section (Bullet Points)
	if len(data.Items) > 0 {
		if cfg.BorderStyle != BorderNone {
			sb.WriteString(drawHorizontal(border.TeeLeft, border.Horizontal, border.TeeRight))
			sb.WriteString("\n")
		} else {
			sb.WriteString("\n")
		}

		for _, item := range data.Items {
			cleanItem := sanitize.Multiline(item)
			itemStr := fmt.Sprintf("%s %s", sh.Style("•", ANSIBoldGreen), sh.Style(cleanItem, ANSIFgWhite))
			for _, wrapped := range WrapText(itemStr, innerWidth) {
				sb.WriteString(wrapBorderLine(wrapped))
				sb.WriteString("\n")
			}
		}
	}

	// Bottom Border
	if cfg.BorderStyle != BorderNone {
		sb.WriteString(drawHorizontal(border.BottomLeft, border.Horizontal, border.BottomRight))
		sb.WriteString("\n")
	}

	return sb.String()
}
