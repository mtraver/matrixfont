package matrixfont

import (
	"errors"
	"fmt"
	"strings"
)

type Glyph struct {
	Codepoint int

	Rows [][]bool

	// DX and DY are the X and Y offsets of the glyph. DX is the offset from the
	// cursor to the left side of the glyph and DY is the offset from the baseline
	// to the bottom of the glyph.
	DX, DY int

	// ShiftX is how far the cursor advances after this glyph is rendered.
	ShiftX int
}

func (g Glyph) Width() int {
	if len(g.Rows) == 0 {
		return 0
	}

	return len(g.Rows[0])
}

func (g Glyph) Height() int {
	return len(g.Rows)
}

func (g Glyph) Ascent() int {
	return g.DY + g.Height()
}

func (g Glyph) Descent() int {
	return -g.DY
}

func (g Glyph) Validate() error {
	// An empty glyph is valid as long it has ShiftX set.
	if len(g.Rows) == 0 {
		if g.ShiftX == 0 {
			return errors.New("no ink and no shift x (advance) set")
		}

		return nil
	}

	// Each glyph's bitmap must be a rectangle (i.e., it must have consistent row lengths).
	width := len(g.Rows[0])
	for i, r := range g.Rows {
		if len(r) != width {
			return fmt.Errorf(
				"row width inconsistent: row %d is %d px wide, expected %d", i, len(r), width)
		}
	}

	return nil
}

func (g Glyph) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "CHAR 0x%x  // %q\n", g.Codepoint, rune(g.Codepoint))
	for i, r := range g.Rows {
		fmt.Fprint(&sb, rowToString(r))
		if i < len(g.Rows)-1 {
			fmt.Fprint(&sb, "\n")
		}
	}

	return sb.String()
}

func (g Glyph) BDF() (string, error) {
	if err := g.Validate(); err != nil {
		return "", err
	}

	var sb strings.Builder

	fmt.Fprintf(&sb, "STARTCHAR U+%04X\n", g.Codepoint)
	fmt.Fprintf(&sb, "ENCODING %d\n", g.Codepoint)
	fmt.Fprintln(&sb, "SWIDTH 1000 0")
	fmt.Fprintf(&sb, "DWIDTH %d 0\n", g.ShiftX)
	fmt.Fprintf(&sb, "BBX %d %d %d %d\n", g.Width(), g.Height(), g.DX, g.DY)
	fmt.Fprintln(&sb, "BITMAP")
	for _, row := range g.Rows {
		fmt.Fprintln(&sb, rowToHex(row))
	}
	fmt.Fprintln(&sb, "ENDCHAR")

	return sb.String(), nil
}

// rowToHex converts a bitmap row to a hex string.
// The bits are left-aligned and padded to a full byte boundary.
func rowToHex(row []bool) string {
	width := len(row)
	numBytes := (width + 7) / 8
	bytes := make([]byte, numBytes)

	for i, bit := range row {
		if bit {
			byteIdx := i / 8
			bitIdx := 7 - (i % 8)
			bytes[byteIdx] |= 1 << bitIdx
		}
	}

	var sb strings.Builder
	for _, b := range bytes {
		fmt.Fprintf(&sb, "%02X", b)
	}
	return sb.String()
}

func rowToString(row []bool) string {
	var sb strings.Builder
	for _, bit := range row {
		if bit {
			sb.WriteByte('#')
		} else {
			sb.WriteByte('.')
		}
	}
	return sb.String()
}
