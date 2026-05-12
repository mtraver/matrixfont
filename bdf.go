package matrixfont

import (
	"fmt"
	"io"
)

// IsMonospace returns true if the font meets the requirements
// for a monospace font as defined in the XLFD spec, that is, that
// each glyph in the font has the same logical width.
func (f Font) IsMonospace() bool {
	if len(f.Glyphs) == 0 {
		return true
	}

	var shiftX *int
	for _, g := range f.Glyphs {
		if shiftX == nil {
			shiftX = &g.ShiftX
		}

		if g.ShiftX != *shiftX {
			return false
		}
	}

	return true
}

// IsCharCell returns true if the font meets the requirements
// for a char cell font as defined in the XLFD spec, that is:
// 1. All glyphs have the same logical width
// 2. No glyphs have ink outside of the character cell
// 3. At the font level, cell height = descent + ascent
func (f Font) IsCharCell() bool {
	if len(f.Glyphs) == 0 {
		return true
	}

	_, maxHeight, _, _ := f.BoundingBox()
	ascent, descent := f.AscentDescent()

	if maxHeight != ascent+descent {
		return false
	}

	var shiftX *int
	for _, g := range f.Glyphs {
		if shiftX == nil {
			shiftX = &g.ShiftX
		}

		if g.ShiftX != *shiftX {
			return false
		}

		// Verify that the horizontal extent doesn't exceed the cell.
		if !((0 <= g.DX) && (g.DX+g.Width() <= g.ShiftX)) {
			return false
		}

		// Verify that the vertical extent doesn't exceed the cell.
		if !(g.Ascent() <= ascent && g.Descent() <= descent) {
			return false
		}
	}

	return true
}

func (f Font) Validate() error {
	// Validate each glyph.
	for _, g := range f.Glyphs {
		if err := g.Validate(); err != nil {
			return fmt.Errorf("glyph U+%04X (%q): %w", g.Rune, g.Rune, err)
		}
	}

	// If a spacing value is specified, validate the font for compliance to spacing requirements.
	if f.Meta.Spacing == SpacingMonospaced && !f.IsMonospace() {
		return fmt.Errorf("spacing value %q given but font is not monospace", f.Meta.Spacing)
	}
	if f.Meta.Spacing == SpacingCharCell && !f.IsCharCell() {
		return fmt.Errorf("spacing value %q given but font is not char cell", f.Meta.Spacing)
	}

	return nil
}

func (f Font) XLFD() string {
	_, maxHeight, _, _ := f.BoundingBox()
	return f.Meta.XLFD(maxHeight, f.AvgWidth())
}

func (f Font) WriteBDF(w io.Writer) error {
	if len(f.Glyphs) == 0 {
		return fmt.Errorf("font has no glyphs")
	}

	if err := f.Validate(); err != nil {
		return err
	}

	ascent, descent := f.AscentDescent()
	maxWidth, maxHeight, minDX, minDY := f.BoundingBox()

	fmt.Fprintln(w, "STARTFONT 2.1")
	fmt.Fprintf(w, "FONT %s\n", f.XLFD())
	fmt.Fprintf(w, "SIZE %d %d %d\n", PointSize(maxHeight, f.Meta.DPI), f.Meta.DPI, f.Meta.DPI)
	fmt.Fprintf(w, "FONTBOUNDINGBOX %d %d %d %d\n", maxWidth, maxHeight, minDX, minDY)
	fmt.Fprintln(w, "STARTPROPERTIES 2")
	fmt.Fprintf(w, "FONT_ASCENT %d\n", ascent)
	fmt.Fprintf(w, "FONT_DESCENT %d\n", descent)
	fmt.Fprintln(w, "ENDPROPERTIES")
	fmt.Fprintf(w, "CHARS %d\n", len(f.Glyphs))

	for _, g := range f.Glyphs {
		glyphBDF, err := g.BDF()
		if err != nil {
			return err
		}

		fmt.Fprint(w, glyphBDF)
	}

	fmt.Fprintln(w, "ENDFONT")

	return nil
}
