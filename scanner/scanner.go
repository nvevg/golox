package scanner

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

var matchLexemeToKeyword = map[string]TokenType{
	"and":    And,
	"class":  Class,
	"else":   Else,
	"fun":    Fun,
	"for":    For,
	"if":     If,
	"nil":    Nil,
	"or":     Or,
	"print":  Print,
	"return": Return,
	"super":  Super,
	"this":   This,
	"true":   True,
	"false":  False,
	"var":    Var,
	"while":  While,
}

func isAllowedInId(r rune) bool {
	return unicode.IsDigit(r) || unicode.IsLetter(r) || r == '_'
}

type ScannerError struct {
	Line    LineNumber
	Column  ColumnNumber
	Message string
}

func (e ScannerError) Error() string {
	return fmt.Sprintf("Error: %s (%d:%d)", e.Message, e.Line, e.Column)
}

type Scanner struct {
	input  bufio.Reader
	Line   LineNumber
	Column ColumnNumber
}

func NewScanner(input bufio.Reader) *Scanner {
	return &Scanner{input: input, Line: 1, Column: 0}
}

func (sc *Scanner) Scan() (*Token, error) {

advance:

	c, e := sc.advance()

	if e != nil {
		return &Token{}, e
	}

	switch c {
	case '(':
		return sc.tok(LeftParen), nil
	case ')':
		return sc.tok(RightParen), nil
	case '{':
		return sc.tok(LeftBrace), nil
	case '}':
		return sc.tok(RightBrace), nil
	case ',':
		return sc.tok(Comma), nil
	case '.':
		return sc.tok(Dot), nil
	case '-':
		return sc.tok(Minus), nil
	case '+':
		return sc.tok(Plus), nil
	case ';':
		return sc.tok(Semicolon), nil
	case '*':
		return sc.tok(Star), nil
	case '/':
		return sc.matchSlash()
	case '!':
		return sc.matchBang()
	case '=':
		return sc.matchEqual()

	case '<':
		return sc.matchLess()

	case '>':
		return sc.matchGreater()

	case '"':
		return sc.matchStringLiteral()

	case '\n':
		sc.Line += 1
		sc.Column = 1
		goto advance
	case '\r':
		goto advance
	case '\t':
		goto advance
	case ' ':
		goto advance

	default:
		// IDs, Numbers, and keywords or...
		if unicode.IsDigit(c) {
			return sc.matchNumberLiteral(c)
		} else if unicode.IsLetter(c) {
			return sc.matchIdOrKeyword(c)
		}

		return nil, sc.unexpectedSymbol(c)
	}
}

func (sc *Scanner) advance() (rune, error) {
	r, sz, e := sc.input.ReadRune()

	// ??? EOF
	if e != nil || (sz == 1 && r == unicode.ReplacementChar) {
		return unicode.ReplacementChar, e
	}

	sc.Column += 1

	return r, nil
}

func (sc *Scanner) tok(tt TokenType) *Token {
	return &Token{Type: tt, Line: sc.Line}
}

func (sc *Scanner) identifier(id string) *Token {
	return &Token{Type: Identifier, Line: sc.Line, Lexeme: id}
}

func (sc *Scanner) comment(text string) *Token {
	return &Token{Type: Comment, Line: sc.Line, Lexeme: text}
}

func (sc *Scanner) literal(tt TokenType, value interface{}) *Token {
	return &Token{Type: tt, Line: sc.Line, Value: value}
}

func (sc *Scanner) makeScannerError(fmts string, args ...interface{}) error {
	return ScannerError{sc.Line, sc.Column, fmt.Sprintf(fmts, args...)}
}

func (sc *Scanner) unexpectedSymbol(r rune) error {
	return sc.makeScannerError("unexpected symbol: %s", string(r))
}

func (sc *Scanner) peek() (rune, error) {
	r, e := sc.advance()
	if e != nil {
		return unicode.ReplacementChar, e
	}

	e = sc.unread()
	if e != nil {
		return unicode.ReplacementChar, e
	}

	return r, nil
}

func (sc *Scanner) unread() error {
	e := sc.input.UnreadRune()
	if e != nil {
		return e
	}

	sc.Column -= 1

	return nil
}

func (sc *Scanner) matchBang() (*Token, error) {
	r, e := sc.peek()
	if e != nil {
		return nil, e
	}

	if r == '=' {
		_, e = sc.advance()
		if e != nil {
			return nil, e
		}

		return sc.tok(BangEqual), nil
	}

	return sc.tok(Bang), nil
}

func (sc *Scanner) matchEqual() (*Token, error) {
	r, e := sc.peek()
	if e != nil {
		return nil, e
	}

	var token *Token = nil
	var consumed bool = false

	switch r {
	case '=':
		token = sc.tok(EqualEqual)
		consumed = true
	case '<':
		token = sc.tok(LessEqual)
		consumed = true
	default:
		token = sc.tok(Equal)
	}

	if consumed {
		_, e = sc.advance()

		if e != nil {
			return nil, e
		}
	}

	return token, nil
}

func (sc *Scanner) matchLess() (*Token, error) {
	c, e := sc.peek()
	if e != nil {
		return nil, e
	}

	if c == '=' {
		_, e = sc.advance()

		if e != nil {
			return nil, e
		}

		return sc.tok(LessEqual), nil
	}

	return sc.tok(Less), nil
}

func (sc *Scanner) matchGreater() (*Token, error) {
	c, e := sc.peek()
	if e != nil {
		return nil, e
	}

	if c == '=' {
		_, e = sc.advance()

		if e != nil {
			return nil, e
		}

		return sc.tok(GreaterEqual), nil
	}

	return sc.tok(Greater), nil
}

func (sc *Scanner) matchStringLiteral() (*Token, error) {
	var sb strings.Builder
	c, e := sc.advance()
	for ; e == nil && c != '"'; c, e = sc.advance() {
		if c == '\n' {
			sc.Line += 1
		}
		sb.WriteRune(c)
	}

	if e != nil && e != io.EOF {
		return nil, e
	} else if e == io.EOF {
		return nil, sc.makeScannerError("unterminated string")
	}

	return sc.literal(String, sb.String()), nil
}

func (sc *Scanner) matchSlash() (*Token, error) {
	c, e := sc.peek()
	if e != nil {
		return nil, e
	}

	var sb strings.Builder

	if c == '/' {
		// munch the second slash
		sc.advance()

		// a comment goes until the end of the line
		for {
			c, e = sc.advance()

			if e != nil && e != io.EOF {
				return nil, e
			} else if e == io.EOF {
				break
			}

			if c == '\n' {
				// let the sc.Scan() handle newlines
				sc.unread()
				break
			}

			sb.WriteRune(c)
		}

		return sc.comment(sb.String()), nil
	}

	return sc.tok(Slash), nil
}

func (sc *Scanner) matchNumberLiteral(first rune) (*Token, error) {
	var sb strings.Builder
	// not missing a first advance()-d digit
	sb.WriteRune(first)

	had_decimal_point := false

munch_numbers:

	c, e := sc.advance()
	for ; e == nil && unicode.IsDigit(c); c, e = sc.advance() {
		sb.WriteRune(c)
	}

	if e != nil && e != io.EOF {
		return nil, e
	}

	if c == '.' && !had_decimal_point {
		sb.WriteRune('.')

		had_decimal_point = true

		goto munch_numbers
	} else {
		// the last advance()-d rune is not a dot - let Scan() handle it
		sc.unread()
	}

	f, e := strconv.ParseFloat(sb.String(), 64)
	if e != nil {
		panic(e.Error())
	}

	return sc.literal(Number, f), nil
}

func (sc *Scanner) matchIdOrKeyword(first rune) (*Token, error) {
	var sb strings.Builder
	// not missing a first advance()-d rune
	sb.WriteRune(first)

	c, e := sc.peek()
	for ; e == nil && isAllowedInId(c); c, e = sc.peek() {
		sb.WriteRune(c)
		sc.advance()
	}

	if e != nil && e != io.EOF {
		return nil, e
	}

	lexeme := sb.String()
	kw, is_kw := matchLexemeToKeyword[lexeme]
	if is_kw {
		return sc.tok(kw), nil
	}

	return sc.identifier(sb.String()), nil
}
