package matrixfont

import (
	"math"
	"strconv"
	"strings"
)

const (
	pointsPerInch = 72.27
)

// Slant codes as defined by the XLFD spec (section 3.1.2.4).
type Slant string

const (
	// Upright design.
	SlantRoman Slant = "R"

	// Italic design, slanted clockwise from the vertical.
	SlantItalic Slant = "I"

	// Obliqued upright design, slanted clockwise from the vertical.
	SlantOblique Slant = "O"

	// Italic design, slanted counterclockwise from the vertical.
	SlantReverseItalic Slant = "RI"

	// Obliqued upright design, slanted counterclockwise from the vertical.
	SlantReverseOblique Slant = "RO"

	// Other.
	SlantOther Slant = "OT"

	// TODO(mtraver) numeric values for polymorphic fonts.
)

// Spacing codes as defined by the XLFD spec (section 3.1.2.10).
type Spacing string

const (
	// A font whose logical character widths vary for each glyph. No other
	// restrictions are placed on the metrics of a proportional font.
	SpacingProportional Spacing = "P"

	// A font whose logical character widths are constant (i.e., every
	// glyph in the font has the same logical width). No other restrictions
	// are placed on the metrics of a monospaced font.
	SpacingMonospaced Spacing = "M"

	// A monospaced font that follows the standard typewriter character cell
	// model (the glyphs of the font can be modeled by X clients as "boxes" of
	// the same width and height that are imaged side-by-side to form text
	// strings or top-to-bottom to form text lines). By definition, all glyphs
	// have the same logical character width, and no glyphs have "ink" outside
	// of the character cell. There is no kerning (on a per-character basis with
	// positive metrics: 0 <= left-bearing <= right-bearing <= width; with negative
	// metrics: width <= left-bearing <= right-bearing <= zero). Also, the vertical
	// extents of the font do not exceed the vertical spacing (that is, on a
	// per-character basis: ascent <= font-ascent and descent <= font-descent).
	// The cell height = font-descent + font-ascent, and the width = average width.
	SpacingCharCell Spacing = "C"
)

// FontName represents an X Logical Font Description (XLFD) as specified in the X Consortium Standard, Version 1.5.
// It holds all 14 fields of an XLFD font name (section 3.1.2), which is formatted
// as a string of those fields separated by hyphens. For example:
//
//	-Adobe-Courier-Medium-R-Normal--10-100-75-75-M-60-ISO8859-1
//
// Reference: https://www.x.org/docs/XLFD/xlfd.pdf
type FontName struct {
	// Foundry is the name or identifier of the digital type foundry that
	// supplied the font data, or if different, the identifier of the
	// organization that last modified the font shape or metric information.
	Foundry string

	// FamilyName identifies the range or family of typeface designs that are all
	// variations of one basic typographic style. This must be spelled out in full,
	// with words separated spaces, as required. This name must be human-understandable
	// and suitable for presentation to a font user to identify the typeface family
	// (e.g., "Helvetica", "ITC Avant Garde Gothic").
	FamilyName string

	// WeightName identifies the font's typographic weight according to the
	// foundry's judgment. This name must be human-understandable and suitable
	// for presentation to a font user (e.g. "Medium", "Bold", "Light"). "0" is
	// used to indicate a polymorphic font.
	WeightName string

	// Slant indicates the overall posture of the typeface design used in the font.
	Slant Slant

	// SetwidthName gives the font's typographic proportionate width (i.e. the
	// nominal width per horizontal unit of the font), according to the foundry's
	// judgment (e.g. "Normal", "Condensed", "Narrow"). "0" is used to indicate a
	// polymorphic font.
	SetwidthName string

	// AddStyleName provides additional typographic style information that is not
	// captured by other fields but is needed to identify the particular font
	// (e.g. "Serif", "Sans Serif", "Informal", "Decorated"). The character "["
	// anywhere in the field is used to indicate a polymorphic font.
	AddStyleName string

	// PixelSize is the body size in pixels of the font at a particular point
	// size and y resolution. PixelSize is either an integer or a string beginning
	// with "[", which represents a matrix. PixelSize usually incorporates additional
	// vertical spacing that is considered part of the font design. (Note, however,
	// that this value is not necessarily equivalent to the height of the font
	// bounding box.) 0 is used to indicate a scalable font.
	PixelSize int

	// PointSize is the body size for which the font was designed. PointSize is
	// either an integer or a string beginning with "[", which represents a matrix.
	// This field usually incorporates additional vertical spacing that is considered
	// part of the font design. (Note, however, that PointSize is not necessarily
	// equivalent to the height of the font bounding box.) PointSize is expressed
	// in decipoints (where points are as defined in the X protocol or 72.27 points
	// equal 1 inch). 0 is used to indicate a scalable font.
	//
	// PointSize (decipoints) = PixelSize * (722.7 / ResolutionY)
	PointSize int

	// ResolutionX is the horizontal resolution, measured in pixels or dots per
	// inch (dpi), for which the font was designed. 0 is used to indicate a
	// scalable font.
	ResolutionX int

	// ResolutionY is the vertical resolution, measured in pixels or dots per
	// inch (dpi), for which the font was designed. 0 is used to indicate a
	// scalable font.
	ResolutionY int

	// Spacing indicates the escapement class of the font. This may be monospace
	// (fixed pitch), proportional (variable pitch), or charcell (a special
	// monospaced font that conforms to the traditional data-processing character
	// cell font model).
	Spacing Spacing

	// AverageWidth is the unweighted arithmetic mean of the absolute value
	// of the width of each glyph, in tenths of pixels. For monospaced
	// fonts this is the width of all glyphs * 10. Use negative values
	// (represented with a leading "~") for right-to-left fonts. Set to 0
	// for scalable fonts.

	// AverageWidth is the unweighted arithmetic mean of the absolute value of the
	// width of each glyph, in tenths of pixels. Use the negative of this value if
	// the dominant writing direction for the font is right-to-left. For monospaced
	// and character cell fonts, this is the width of all glyphs in the font. 0 is
	// used to indicate a scalable font.
	AverageWidth int

	// CharsetRegistry identifies the registration authority that owns the font's
	// character set encoding (e.g. "ISO8859").
	CharsetRegistry string

	// CharsetEncoding identifies the specific character set (e.g. "1" for
	// ISO 8859-1, which is ISO Latin-1). May include a subsetting hint in brackets.
	CharsetEncoding string
}

// String returns the full XLFD font name string.
// Numeric fields are rendered as their integer string equivalents.
// Empty string fields produce empty field values (consecutive hyphens).
func (f FontName) String() string {
	// A leading "~" (tilde) is used to indicate a negative average width.
	averageWidth := strconv.Itoa(abs(f.AverageWidth))
	if f.AverageWidth < 0 {
		averageWidth = "~" + averageWidth
	}

	fields := []string{
		f.Foundry,
		f.FamilyName,
		f.WeightName,
		string(f.Slant),
		f.SetwidthName,
		f.AddStyleName,
		strconv.Itoa(f.PixelSize),
		strconv.Itoa(f.PointSize),
		strconv.Itoa(f.ResolutionX),
		strconv.Itoa(f.ResolutionY),
		string(f.Spacing),
		averageWidth,
		f.CharsetRegistry,
		f.CharsetEncoding,
	}
	return "-" + strings.Join(fields, "-")
}

// PointSizeDecipoints computes the correct point size in decipoints given a
// pixel size and vertical resolution (dpi).
// The spec says, "Design POINT_SIZE cannot be calculated or approximated." This
// is because point size is a design intent value the relates to the dpi for which
// the font was designed. However, since there is a mathematical relationship between
// point size, pixel size, and dpi, this function is provided to calculate it.
func PointSizeDecipoints(pixelSize int, verticalDPI int) int {
	return int(math.Round(pointSize(pixelSize, verticalDPI) * 10))
}

// PointSize computes the correct point size in whole points given a pixel size
// and vertical resolution (dpi).
// The spec says, "Design POINT_SIZE cannot be calculated or approximated." This
// is because point size is a design intent value the relates to the dpi for which
// the font was designed. However, since there is a mathematical relationship between
// point size, pixel size, and dpi, this function is provided to calculate it.
func PointSize(pixelSize int, verticalDPI int) int {
	return int(math.Round(pointSize(pixelSize, verticalDPI)))
}

func pointSize(pixelSize int, verticalDPI int) float64 {
	if verticalDPI == 0 {
		return 0.0
	}

	return float64(pixelSize) * pointsPerInch / float64(verticalDPI)
}

// PixelSize computes the pixel size given a point size in decipoints and
// a vertical resolution (dpi).
//
// From the spec:
//
//	 DeciPointsPerInch = 722.7
//		PIXEL_SIZE = ROUND((RESOLUTION_Y * POINT_SIZE) / DeciPointsPerInch)
func PixelSize(pointSizeDecipoints int, resolutionY int) int {
	return int(math.Round(float64(resolutionY) * float64(pointSizeDecipoints) / 722.7))
}

// Scalable returns a copy of the FontName with PixelSize, PointSize,
// ResolutionX, ResolutionY, and AverageWidth set to 0, which is the
// XLFD convention for scalable fonts.
func (f FontName) Scalable() FontName {
	f.PixelSize = 0
	f.PointSize = 0
	f.ResolutionX = 0
	f.ResolutionY = 0
	f.AverageWidth = 0
	return f
}
