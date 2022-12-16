package scanner

import (
	"fmt"
)

type TokenType int64
type LineNumber uint32
type ColumnNumber uint32

const (
	// Single-character tokens.
	LeftParen = iota
	RightParen
	LeftBrace
	RightBrace
	Comma
	Dot
	Minus
	Plus
	Semicolon
	Slash
	Star

	// One or two character tokens.
	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual

	Comment

	// Literals.
	Identifier
	String
	Number

	// Keywords.
	And
	Class
	Else
	Fun
	For
	If
	Nil
	Or
	Print
	Return
	Super
	This
	True
	False
	Var
	While

	Eof
)

var tokenTypeTostring = map[TokenType]string{
	// Single-character tokens.
	LeftParen:  "LeftParen",
	RightParen: "RightParen",
	LeftBrace:  "LeftBrace",
	RightBrace: "RightBrace",
	Comma:      "Comma",
	Dot:        "Dot",
	Minus:      "Minus",
	Plus:       "Plus",
	Semicolon:  "Semicolon",
	Slash:      "Slash",
	Star:       "Star",

	// One or two character tokens.
	Bang:         "Bang",
	BangEqual:    "BangEqual",
	Equal:        "Equal",
	EqualEqual:   "EqualEqual",
	Greater:      "Greater",
	GreaterEqual: "GreaterEqual",
	Less:         "Less",
	LessEqual:    "LessEqual",

	Comment: "Comment",

	// Literals.
	Identifier: "Identifier",
	String:     "Literal value",
	Number:     "Literal value",

	// Keywords.
	And:    "And",
	Class:  "Class",
	Else:   "Else",
	Fun:    "Fun",
	For:    "For",
	If:     "If",
	Nil:    "Nil",
	Or:     "Or",
	Print:  "Print",
	Return: "Return",
	Super:  "Super",
	This:   "This",
	True:   "True",
	Var:    "Var",
	While:  "While",

	Eof: "Eof",
}

type Token struct {
	Type   TokenType
	Lexeme string
	Line   LineNumber
	Value  interface{}
}

func (t Token) String() string {

	ts, prs := tokenTypeTostring[t.Type]

	if !prs {
		panic(fmt.Sprintf("unable to match TokenType to name: %d", t.Type))
	}

	var str string = fmt.Sprintf("token Type: %s Line: %d ", ts, t.Line)

	switch t.Type {
	case Comment:
		str += fmt.Sprintf("comment text: \"%s\"", t.Lexeme)
	case String:
		if t.Value == nil {
			panic("token type is a value literal but no value provided")
		}

		val, ok := t.Value.(string)
		if !ok {
			panic("token type is String but value is of wrong type")
		}
		str += fmt.Sprintf("string literal: \"%s\"", val)

	case Number:
		if t.Value == nil {
			panic("token type is a value literal but no value provided")
		}

		val, ok := t.Value.(float64)
		if !ok {
			panic("token type is String but value is of wrong type")
		}

		str += fmt.Sprintf("number literal: \"%f\"", val)

	case Identifier:
		str += fmt.Sprintf("identifier: \"%s\"", t.Lexeme)
	}

	return str
}
