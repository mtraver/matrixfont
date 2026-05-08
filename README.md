# Easily create and edit fonts for LED matrices

What are your font options for use with LED matrices such as [those from Adafruit](https://www.adafruit.com/category/327)? Well, there's the built-in `terminalio.FONT` which looks good on an LED matrix but is less compact than it could be. Tools such as `otf2bdf` can convert .ttf files to [BDF](https://en.wikipedia.org/wiki/Glyph_Bitmap_Distribution_Format) format, which CircuitPython supports, so theorethically any font can be used on an LED matrix. Whether or not such converted fonts look good is another story... When displayed on a very low DPI LED matrix, and especially if scaled down, they often look terrible.

Unsatisfied with existing options, I embarked on the task of creating my own font. To do so I wrote this tool that allows a font to be created via an easy-to-edit, visual, text-based format.

Define your font like this:

```txt
// These optional header key-value pairs define font metadata used
// to generate an XLFD font name string for the BDF file.
FOUNDRY mtraver
FAMILY ledmatrix
WEIGHT light
SLANT r
WIDTH normal
STYLE sans serif
DPI 6
SPACING c
CHARSET_REGISTRY ISO8859
CHARSET_ENCODING 1

// Glyph are defined by a "CHAR" line containing the codepoint
// value in either hex or decimal, followed by a grid of '#' and '.'
// characters representing the glyph's bitmap.
CHAR 0x41  // "A"
.###..
#...#.
#...#.
#####.
#...#.
#...#.
#...#.

CHAR 0x42  // "B"
####..
#...#.
#...#.
####..
#...#.
#...#.
####..

CHAR 0x43  // "C"
.###..
#...#.
#.....
#.....
#.....
#...#.
.###..
```

Then create a .bdf file from it like this:

```sh
go run ./cmd/txt2bdf my_font.txt my_font.bdf
```

Load it onto your Adafruit Matrix Portal or whatever and off you go!
