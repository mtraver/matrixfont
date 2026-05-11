import terminalio

from adafruit_bitmap_font import bitmap_font


# Prints the given font's glyphs in the format used by github.com/mtraver/matrixfont.
# The default codepoint range is the ASCII printable range of 0x20 (space) to 0x7e (~).
def print_glyphs(font, first_codepoint=0x20, last_codepoint=0x7e):
  loaded_font = None
  if font is terminalio.FONT:
    loaded_font = font
  elif isinstance(font, str):
    loaded_font = bitmap_font.load_font(font)
  else:
    raise ValueError("Font must be terminalio.FONT or a path to a .bdf or .pcf file")

  # BDF fonts and terminalio.FONT have a shared bitmap / sprite sheet,
  # whereas PCF fonts have per-glyph bitmaps.
  shared_bitmap = hasattr(loaded_font, "bitmap")

  tile_w = None
  if shared_bitmap:
    tile_w = loaded_font.get_bounding_box()[0]

  for codepoint in range(first_codepoint, last_codepoint + 1):
    glyph = loaded_font.get_glyph(codepoint)
    if glyph is None:
      continue

    g_w = glyph.width
    g_h = glyph.height

    print("CHAR {}  // \"{}\"".format(hex(codepoint), chr(codepoint)))
    print("XOFF {}".format(glyph.dx))
    print("YOFF {}".format(glyph.dy))
    print("ADVANCE {}".format(glyph.shift_x))
    if g_w == 0 or g_h == 0:
      print("// no ink (empty glyph)")
      print("")
      continue

    for row in range(g_h):
      line = ""
      for col in range(g_w):
        if shared_bitmap:
          tile_x = glyph.tile_index * tile_w
          pixel = loaded_font.bitmap[tile_x + col, row]
        else:
          pixel = glyph.bitmap[col, row]

        line += "#" if pixel else "."

      print(line)

    print("")
