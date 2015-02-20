package token

import (
	"fmt"
)

type Token int

const (
	ILLEGAL Token = iota
	EOF
	COMMENT

	IDENT // identifier.

	literal_beg
	NUM    // '1.2'
	LT     // '<'
	GT     // '>'
	LE     // '<='
	GE     // '>='
	EQ     // '='
	LBRACK // '['
	RBRACK // ']'
	LPAREN // '('
	RPAREN // ')'
	COMMA  // ','
	ADD    // '+'
	SUB    // '-'
	MULT   // '*'
	DIV    // '/'
	literal_end

	keyword_beg
	TRUE  // 'true'
	FALSE // 'false'
	DO    // 'do'
	DEF   // 'def', declare variable.
	LET   // 'let'
	IF    // 'if'
	FN    // 'fn'
	keyword_end
)

var tokens = [...]string{
	LT:     "<",
	GT:     ">",
	LE:     "<=",
	GE:     ">=",
	EQ:     "=",
	LBRACK: "[",
	RBRACK: "]",
	LPAREN: "(",
	RPAREN: ")",
	COMMA:  ",",
	ADD:    "+",
	SUB:    "-",
	MULT:   "*",
	DIV:    "/",
	TRUE:   "true",
	FALSE:  "false",
	DO:     "do",
	DEF:    "def",
	LET:    "let",
	IF:     "if",
	FN:     "fn",
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for i := keyword_beg; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
	fmt.Println("len", len(keywords))
}

func Lookup(ident string) Token {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
