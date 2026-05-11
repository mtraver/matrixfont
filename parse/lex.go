//go:generate stringer -linecomment -type=tokenType
package parse

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type tokenType int

const (
	tokenEOF   tokenType = iota // EOF
	tokenError                  // error

	// Header tokens.
	tokenFOUNDRY          // FOUNDRY
	tokenFAMILY           // FAMILY
	tokenWEIGHT           // WEIGHT
	tokenSLANT            // SLANT
	tokenWIDTH            // WIDTH
	tokenSTYLE            // STYLE
	tokenDPI              // DPI
	tokenSPACING          // SPACING
	tokenCHARSET_REGISTRY // CHARSET_REGISTRY
	tokenCHARSET_ENCODING // CHARSET_ENCODING

	// Glyph tokens.
	tokenCHAR      // CHAR
	tokenXOFF      // XOFF
	tokenYOFF      // YOFF
	tokenADVANCE   // ADVANCE
	tokenBitmapRow // bitmap row
)

var linePatterns = []struct {
	re      *regexp.Regexp
	typ     tokenType
	process func(string) token
}{
	// Header section.
	{regexp.MustCompile(`^FOUNDRY\s+(.+)$`), tokenFOUNDRY, strToken(tokenFOUNDRY)},
	{regexp.MustCompile(`^FAMILY\s+(.+)$`), tokenFAMILY, strToken(tokenFAMILY)},
	{regexp.MustCompile(`^WEIGHT\s+(.+)$`), tokenWEIGHT, strToken(tokenWEIGHT)},
	{regexp.MustCompile(`^SLANT\s+(.+)$`), tokenSLANT, strToken(tokenSLANT)},
	{regexp.MustCompile(`^WIDTH\s+(.+)$`), tokenWIDTH, strToken(tokenWIDTH)},
	{regexp.MustCompile(`^STYLE\s+(.+)$`), tokenSTYLE, strToken(tokenSTYLE)},
	{regexp.MustCompile(`^DPI\s+(.+)$`), tokenDPI, intToken(tokenDPI)},
	{regexp.MustCompile(`^SPACING\s+(.+)$`), tokenSPACING, strToken(tokenSPACING)},
	{regexp.MustCompile(`^CHARSET_REGISTRY\s+(.+)$`), tokenCHARSET_REGISTRY, strToken(tokenCHARSET_REGISTRY)},
	{regexp.MustCompile(`^CHARSET_ENCODING\s+(.+)$`), tokenCHARSET_ENCODING, strToken(tokenCHARSET_ENCODING)},

	// Glyph section.
	{regexp.MustCompile(`^CHAR\s+(.+)$`), tokenCHAR, codepointToken(tokenCHAR)},
	{regexp.MustCompile(`^XOFF\s+(.+)$`), tokenXOFF, intToken(tokenXOFF)},
	{regexp.MustCompile(`^YOFF\s+(.+)$`), tokenYOFF, intToken(tokenYOFF)},
	{regexp.MustCompile(`^ADVANCE\s+(.+)$`), tokenADVANCE, intToken(tokenADVANCE)},
	{regexp.MustCompile(`^([#.]+)$`), tokenBitmapRow, strToken(tokenBitmapRow)},
}

type token struct {
	typ tokenType

	intValue int
	strValue string
	err      error
}

type lexer struct {
	scanner *bufio.Scanner
	tokens  chan token
}

func lex(r io.Reader) *lexer {
	l := &lexer{
		scanner: bufio.NewScanner(r),
		tokens:  make(chan token),
	}

	go l.run()
	return l
}

func (l *lexer) run() {
	defer close(l.tokens)

	for l.scanner.Scan() {
		line := preprocess(l.scanner.Text())

		// Skip blank lines. Comment lines will have been transformed into
		// blank lines by preprocessing so this also skips comments.
		if line == "" {
			continue
		}

		l.tokens <- l.classifyLine(line)
	}

	if err := l.scanner.Err(); err != nil {
		l.tokens <- token{typ: tokenError, err: err}
	}

	l.tokens <- token{typ: tokenEOF}
}

func (l *lexer) classifyLine(line string) token {
	for _, p := range linePatterns {
		if matches := p.re.FindStringSubmatch(line); matches != nil {
			return p.process(matches[1])
		}
	}

	return token{typ: tokenError, err: fmt.Errorf("unrecognized line: %q", line)}
}

func preprocess(line string) string {
	line = strings.TrimSpace(line)

	// Strip // comments. This works for both leading and inline comments.
	if idx := strings.Index(line, "//"); idx != -1 {
		line = strings.TrimSpace(line[:idx])
	}

	return line
}

func strToken(typ tokenType) func(string) token {
	return func(s string) token {
		return token{typ: typ, strValue: s}
	}
}

func intToken(typ tokenType) func(string) token {
	return func(s string) token {
		v, err := strconv.Atoi(s)
		if err != nil {
			return token{typ: tokenError, err: fmt.Errorf("invalid %s value %q: must be an integer", typ, s)}
		}
		return token{typ: typ, intValue: v}
	}
}

func codepointToken(typ tokenType) func(string) token {
	return func(s string) token {
		cp, err := parseCodepoint(s)
		if err != nil {
			return token{typ: tokenError, err: fmt.Errorf("invalid codepoint %q: %w", s, err)}
		}
		return token{typ: typ, intValue: cp}
	}
}

func parseCodepoint(s string) (int, error) {
	lower := strings.ToLower(s)
	if strings.HasPrefix(lower, "u+") || strings.HasPrefix(lower, "0x") {
		v, err := strconv.ParseInt(lower[2:], 16, 32)
		return int(v), err
	}
	v, err := strconv.ParseInt(s, 10, 32)
	return int(v), err
}
