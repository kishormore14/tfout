package formatter

// BorderChars holds characters used to draw box borders.
type BorderChars struct {
	TopLeft     string
	TopRight    string
	BottomLeft  string
	BottomRight string
	Horizontal  string
	Vertical    string
	TeeLeft     string
	TeeRight    string
}

// GetBorderChars returns the border character set based on style and theme.
func GetBorderChars(style BorderStyle, theme Theme) BorderChars {
	if style == BorderNone {
		return BorderChars{
			TopLeft:     "",
			TopRight:    "",
			BottomLeft:  "",
			BottomRight: "",
			Horizontal:  "",
			Vertical:    "",
			TeeLeft:     "",
			TeeRight:    "",
		}
	}

	if theme == ThemeMono {
		// ASCII borders for mono theme
		if style == BorderDouble {
			return BorderChars{
				TopLeft:     "+",
				TopRight:    "+",
				BottomLeft:  "+",
				BottomRight: "+",
				Horizontal:  "=",
				Vertical:    "|",
				TeeLeft:     "+",
				TeeRight:    "+",
			}
		}
		return BorderChars{
			TopLeft:     "+",
			TopRight:    "+",
			BottomLeft:  "+",
			BottomRight: "+",
			Horizontal:  "-",
			Vertical:    "|",
			TeeLeft:     "+",
			TeeRight:    "+",
		}
	}

	switch style {
	case BorderDouble:
		return BorderChars{
			TopLeft:     "╔",
			TopRight:    "╗",
			BottomLeft:  "╚",
			BottomRight: "╝",
			Horizontal:  "═",
			Vertical:    "║",
			TeeLeft:     "╠",
			TeeRight:    "╣",
		}
	case BorderSingle:
		fallthrough
	default:
		return BorderChars{
			TopLeft:     "┌",
			TopRight:    "┐",
			BottomLeft:  "└",
			BottomRight: "┘",
			Horizontal:  "─",
			Vertical:    "│",
			TeeLeft:     "├",
			TeeRight:    "┤",
		}
	}
}
