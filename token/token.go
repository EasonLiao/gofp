package token

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
	DEFN  // 'defn', declare function.
	LET   // 'let'
	IF    // 'if'
	FN    // 'fn'
	keyword_end
)

var tokens = [...]string{
	ILLEGAL: "[ILLEGAL]",
	EOF:     "[EOF]",
	NUM:     "[NUM]",
	LT:      "<",
	GT:      ">",
	LE:      "<=",
	GE:      ">=",
	EQ:      "=",
	LBRACK:  "[",
	RBRACK:  "]",
	LPAREN:  "(",
	RPAREN:  ")",
	COMMA:   ",",
	ADD:     "+",
	SUB:     "-",
	MULT:    "*",
	DIV:     "/",
	TRUE:    "true",
	FALSE:   "false",
	DO:      "do",
	DEF:     "def",
	DEFN:    "defn",
	LET:     "let",
	IF:      "if",
	FN:      "fn",
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for i := keyword_beg; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}

func Lookup(ident string) Token {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

func TokenName(tok Token) string {
	if int(tok) > len(tokens) {
		return "[INVALID TOKEN]"
	}
	return tokens[tok]
}
