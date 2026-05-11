package parse

import (
	"fmt"
	"io"
	"strings"

	"github.com/mtraver/matrixfont"
)

type parseContext struct {
	meta    matrixfont.Metadata
	glyphs  []matrixfont.Glyph
	current *matrixfont.Glyph
}

func (ctx *parseContext) startGlyph(codepoint int) {
	ctx.current = &matrixfont.Glyph{Codepoint: codepoint}
}

func (ctx *parseContext) flush() {
	if ctx.current != nil {
		ctx.glyphs = append(ctx.glyphs, *ctx.current)
		ctx.current = nil
	}
}

func (ctx *parseContext) font() matrixfont.Font {
	return matrixfont.Font{Meta: ctx.meta, Glyphs: ctx.glyphs}
}

type stateFn func(*parseContext, token) (stateFn, error)

type Parser struct {
	lexer *lexer
	ctx   *parseContext
}

func Parse(r io.Reader) (matrixfont.Font, error) {
	p := &Parser{
		lexer: lex(r),
		ctx:   &parseContext{},
	}
	return p.run()
}

func (p *Parser) run() (matrixfont.Font, error) {
	var err error
	var state stateFn = p.parseMetadata
	for state != nil {
		tok := <-p.lexer.tokens
		if tok.typ == tokenEOF {
			break
		}

		if tok.typ == tokenError {
			return matrixfont.Font{}, tok.err
		}

		state, err = state(p.ctx, tok)
		if err != nil {
			return matrixfont.Font{}, err
		}
	}

	p.ctx.flush()
	font, err := postprocess(p.ctx.font())
	if err != nil {
		return matrixfont.Font{}, err
	}

	if err := font.Validate(); err != nil {
		return matrixfont.Font{}, err
	}

	return font, nil
}

func (p *Parser) parseMetadata(ctx *parseContext, tok token) (stateFn, error) {
	switch tok.typ {
	case tokenCHAR:
		ctx.startGlyph(tok.intValue)
		return p.parseGlyphMeta, nil

	case tokenFOUNDRY:
		ctx.meta.Foundry = tok.strValue
		return p.parseMetadata, nil

	case tokenFAMILY:
		ctx.meta.Family = tok.strValue
		return p.parseMetadata, nil

	case tokenWEIGHT:
		ctx.meta.Weight = tok.strValue
		return p.parseMetadata, nil

	case tokenSLANT:
		ctx.meta.Slant = matrixfont.Slant(strings.ToUpper(tok.strValue))
		return p.parseMetadata, nil

	case tokenWIDTH:
		ctx.meta.Width = tok.strValue
		return p.parseMetadata, nil

	case tokenSTYLE:
		ctx.meta.Style = tok.strValue
		return p.parseMetadata, nil

	case tokenDPI:
		ctx.meta.DPI = tok.intValue
		return p.parseMetadata, nil

	case tokenSPACING:
		ctx.meta.Spacing = matrixfont.Spacing(strings.ToUpper(tok.strValue))
		return p.parseMetadata, nil

	case tokenCHARSET_REGISTRY:
		ctx.meta.CharsetRegistry = tok.strValue
		return p.parseMetadata, nil

	case tokenCHARSET_ENCODING:
		ctx.meta.CharsetEncoding = tok.strValue
		return p.parseMetadata, nil

	default:
		return nil, fmt.Errorf("unexpected token in metadata section: %v", tok)
	}
}

func (p *Parser) parseGlyphMeta(ctx *parseContext, tok token) (stateFn, error) {
	switch tok.typ {
	case tokenCHAR:
		ctx.flush()
		ctx.startGlyph(tok.intValue)
		return p.parseGlyphMeta, nil

	case tokenXOFF:
		ctx.current.DX = tok.intValue
		return p.parseGlyphMeta, nil

	case tokenYOFF:
		ctx.current.DY = tok.intValue
		return p.parseGlyphMeta, nil

	case tokenADVANCE:
		ctx.current.ShiftX = tok.intValue
		return p.parseGlyphMeta, nil

	case tokenBitmapRow:
		ctx.current.Rows = append(ctx.current.Rows, parseBitmapRow(tok.strValue))
		return p.parseGlyphBitmap, nil

	default:
		return nil, fmt.Errorf("unexpected token in glyph metadata: %v", tok)
	}
}

func (p *Parser) parseGlyphBitmap(ctx *parseContext, tok token) (stateFn, error) {
	switch tok.typ {
	case tokenCHAR:
		ctx.flush()
		ctx.startGlyph(tok.intValue)
		return p.parseGlyphMeta, nil

	case tokenBitmapRow:
		ctx.current.Rows = append(ctx.current.Rows, parseBitmapRow(tok.strValue))
		return p.parseGlyphBitmap, nil

	default:
		return nil, fmt.Errorf("unexpected token in glyph bitmap: %v", tok)
	}
}

func parseBitmapRow(row string) []bool {
	result := make([]bool, len(row))
	for i, c := range row {
		result[i] = c == '#'
	}

	return result
}

func postprocess(font matrixfont.Font) (matrixfont.Font, error) {
	// Empty glyphs must have ADVANCE set.
	for i, g := range font.Glyphs {
		if g.Width() == 0 && g.ShiftX == 0 {
			return matrixfont.Font{}, fmt.Errorf("glyph U+%04X has no ink and no ADVANCE set", g.Codepoint)
		}
		if g.ShiftX == 0 {
			font.Glyphs[i].ShiftX = g.DX + g.Width()
		}
	}

	return font, nil
}
