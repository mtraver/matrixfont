package matrixfont

// Header holds the XLFD metadata parsed from the top of a matrix font file.
type Header struct {
	Foundry         string
	Family          string
	Weight          string
	Slant           Slant
	Width           string
	Style           string
	DPI             int
	Spacing         Spacing
	CharsetRegistry string
	CharsetEncoding string
}

func (h Header) XLFD(glyphHeightPx, avgGlyphWidhtPx int) string {
	foundry := h.Foundry
	if foundry == "" {
		foundry = "unknown"
	}

	family := h.Family
	if family == "" {
		family = "unknown"
	}

	fontName := FontName{
		Foundry:         foundry,
		FamilyName:      family,
		WeightName:      h.Weight,
		Slant:           Slant(h.Slant),
		SetwidthName:    h.Width,
		PixelSize:       glyphHeightPx,
		PointSize:       PointSizeDecipoints(glyphHeightPx, h.DPI),
		ResolutionX:     h.DPI,
		ResolutionY:     h.DPI,
		Spacing:         Spacing(h.Spacing),
		AverageWidth:    avgGlyphWidhtPx * 10, // decipixels
		CharsetRegistry: h.CharsetRegistry,
		CharsetEncoding: h.CharsetEncoding,
	}

	return fontName.String()
}
