package matrixfont

import "math"

type Font struct {
	Header Header
	Glyphs []Glyph
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
