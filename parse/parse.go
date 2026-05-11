package parse

import (
	"fmt"
	"io"
	"strings"

	"github.com/mtraver/matrixfont"
)

type parseContext struct {
	header  matrixfont.Header
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
	return matrixfont.Font{Header: ctx.header, Glyphs: ctx.glyphs}
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
	var state stateFn = p.parseHeader
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
	return postprocess(p.ctx.font())
}

func (p *Parser) parseHeader(ctx *parseContext, tok token) (stateFn, error) {
	switch tok.typ {
	case tokenCHAR:
		ctx.startGlyph(tok.intValue)
		return p.parseGlyphMeta, nil

	case tokenFOUNDRY:
		ctx.header.Foundry = tok.strValue
		return p.parseHeader, nil

	case tokenFAMILY:
		ctx.header.Family = tok.strValue
		return p.parseHeader, nil

	case tokenWEIGHT:
		ctx.header.Weight = tok.strValue
		return p.parseHeader, nil

	case tokenSLANT:
		ctx.header.Slant = matrixfont.Slant(strings.ToUpper(tok.strValue))
		return p.parseHeader, nil

	case tokenWIDTH:
		ctx.header.Width = tok.strValue
		return p.parseHeader, nil

	case tokenSTYLE:
		ctx.header.Style = tok.strValue
		return p.parseHeader, nil

	case tokenDPI:
		ctx.header.DPI = tok.intValue
		return p.parseHeader, nil

	case tokenSPACING:
		ctx.header.Spacing = matrixfont.Spacing(strings.ToUpper(tok.strValue))
		return p.parseHeader, nil

	case tokenCHARSET_REGISTRY:
		ctx.header.CharsetRegistry = tok.strValue
		return p.parseHeader, nil

	case tokenCHARSET_ENCODING:
		ctx.header.CharsetEncoding = tok.strValue
		return p.parseHeader, nil

	default:
		return nil, fmt.Errorf("unexpected token in header: %v", tok)
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
