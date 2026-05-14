package matrixfont

import (
	"math"
	"sort"
	"unicode"
)

type Font struct {
	Meta   Metadata
	Glyphs map[rune]Glyph
}

func (f Font) AscentDescent() (int, int) {
	var fontAscent, fontDescent int
	for _, g := range f.Glyphs {
		if len(g.Rows) == 0 {
			continue
		}

		ascent := g.Ascent()
		descent := g.Descent()

		if ascent > fontAscent {
			fontAscent = ascent
		}

		if descent > fontDescent {
			fontDescent = descent
		}
	}

	return fontAscent, fontDescent
}

func (f Font) AvgWidth() int {
	var total int
	for _, g := range f.Glyphs {
		total += g.ShiftX
	}

	return int(math.Round(float64(total) / float64(len(f.Glyphs))))
}

func (f Font) BoundingBox() (int, int, int, int) {
	var minDX, minDY int = math.MaxInt32, math.MaxInt32
	var maxRight, maxTop int = math.MinInt32, math.MinInt32

	for _, g := range f.Glyphs {
		if len(g.Rows) == 0 {
			continue
		}
		if g.DX < minDX {
			minDX = g.DX
		}
		if g.DY < minDY {
			minDY = g.DY
		}
		if g.DX+g.Width() > maxRight {
			maxRight = g.DX + g.Width()
		}
		if g.DY+g.Height() > maxTop {
			maxTop = g.DY + g.Height()
		}
	}

	maxWidth := maxRight - minDX
	maxHeight := maxTop - minDY

	return maxWidth, maxHeight, minDX, minDY
}

// OrderedPrintableRunes returns the font's printable runes in preview order: A-Z, a-z, 0-9, other printable ASCII, then any remaining printable runes in sorted order.
func (f Font) OrderedPrintableRunes() []rune {
	var result []rune
	seen := make(map[rune]bool)

	add := func(r rune) {
		if _, inFont := f.Glyphs[r]; inFont && !seen[r] {
			result = append(result, r)
			seen[r] = true
		}
	}

	for r := 'A'; r <= 'Z'; r++ {
		add(r)
	}
	for r := 'a'; r <= 'z'; r++ {
		add(r)
	}
	for r := '0'; r <= '9'; r++ {
		add(r)
	}
	for r := rune(0x20); r <= 0x7e; r++ {
		add(r)
	}

	var remaining []rune
	for r := range f.Glyphs {
		if !seen[r] && unicode.IsPrint(r) {
			remaining = append(remaining, r)
		}
	}
	sort.Slice(remaining, func(i, j int) bool { return remaining[i] < remaining[j] })
	result = append(result, remaining...)

	return result
}
