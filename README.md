# Easily create and edit fonts for LED matrices

Create a font via an easy-to-edit, visual, text-based format!

Define your font like this:

```text
// These optional header key-value pairs define font metadata used
// to generate an XLFD font name string for the BDF file.
FOUNDRY mtraver
FAMILY ledmatrix
WEIGHT light
SLANT r
WIDTH normal
STYLE sans serif
DPI 6  // The DPI of a 4mm pitch LED matrix is 6.4
SPACING c
CHARSET_REGISTRY ISO8859
CHARSET_ENCODING 1

// Glyphs are defined by a "CHAR" line containing the codepoint
// value as either hex or decimal, followed by a grid of '#' and '.'
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

## Modify an existing font

Want to start from `terminalio.FONT` or a font that you've converted from .ttf to .bdf/.pcf and modify it? You can do that! First you'll need to export it in matrixfont format from your CircuitPython device. Copy `font_utils.py` onto your device and use it like this:

```python
import terminalio
import font_utils

# Export the terminalio default font.
print("======= START FONT =======")
font_utils.print_glyphs(terminalio.FONT)
print("======= END FONT =======")

# Export any font file.
print("======= START FONT =======")
font_utils.print_glyphs("lib/font_free_mono_12/font.pcf")
print("======= END FONT =======")
```

Save the output between the START FONT and END FONT lines to a file, edit the glyphs as you please, and convert it to .bdf using `txt2bdf` as described above.

## Background

When building my first LED matrix project (using [a matrix from Adafruit](https://www.adafruit.com/category/327)) I started off using CircuitPython's built-in `terminalio.FONT`. It looks good on an LED matrix but it's taller than it needs to be. Given the very low DPI of LED matrices (a 4mm pitch LED matrix has a DPI of 6.4), I wanted a more compact but still readable font. Space is at a premium!

Aside from `terminalio.FONT`, the easy and well-documented font options are:
1. The pre-made fonts from Adafruit (download one of the releases, which contain .pcf files): https://github.com/adafruit/circuitpython-fonts
2. Using `otf2bdf` to convert any .ttf file to [BDF](https://en.wikipedia.org/wiki/Glyph_Bitmap_Distribution_Format) format, which CircuitPython supports. The .bdf can then then be converted to .pcf if you like.

I was unsatisfied with these options. The core issue is that taking a modern, detailed font and displaying it on an exceptionally low DPI LED matrix will either 1) require a huge number of pixels (perhaps 15 pixels tall, which on a 64x32 matrix means _at most_ two lines of text in landscape orientation) or 2) downsampling to the point of illegibility. For example, FreeMono 8, 10, and 12, all included in the Adafruit font bundle, look quite bad. Capitals in FreeMono 12 are 7 pixels tall and with 7 pixels it's possible to achieve attractive, legible text on an LED matrix. I wanted to make that happen!

Given all this I set out to create my own fonts specifically designed for LED matrices and thus this tool was born. I hope it proves useful to someone else!
