package scanner

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/easonliao/gofp/token"
)

type Scanner struct {
	src      []byte
	offset   int
	rdoffset int
	err      error
	ch       rune
}

func (s *Scanner) Init(src []byte) {
	s.src = src
	s.next()
}

func (s *Scanner) Next() (tok token.Token, lit string, err error) {
	s.skipWhitespaces()

	switch ch := s.ch; {
	case isLetter(ch):
		lit = s.scanIdent()
		tok = token.Lookup(lit)

	case isDigit(ch):
		lit = s.scanNum()
		tok = token.NUM

	default:
		switch ch {
		case -1:
			tok = token.EOF
		case '+':
			tok = token.ADD
		case '-':
			tok = token.SUB
		case '/':
			tok = token.DIV
		case '*':
			tok = token.MULT
		case '=':
			tok = token.EQ
		case '[':
			tok = token.LBRACK
		case ']':
			tok = token.RBRACK
		case '(':
			tok = token.LPAREN
		case ')':
			tok = token.RPAREN
		case ',':
			tok = token.COMMA
		case '>':
			tok = token.GT
		case '<':
			tok = token.LT
		default:
			s.errorf("unregonized token %c", ch)
		}
		lit = ""
		s.next()
		if tok == token.LT || tok == token.GT {
			if s.ch == '=' {
				if tok == token.LT {
					tok = token.LE
				} else {
					tok = token.GE
				}
				s.next()
			}
		}
	}
	err = s.err
	return
}

func (s *Scanner) scanIdent() string {
	off := s.offset
	for isLetter(s.ch) || isDigit(s.ch) {
		s.next()
	}
	return string(s.src[off:s.offset])
}

func (s *Scanner) scanNum() string {
	off := s.offset
	for isDigit(s.ch) {
		s.next()
	}
	if s.ch == '.' {
		s.next()
		for isDigit(s.ch) {
			s.next()
		}
	}
	return string(s.src[off:s.offset])
}

func (s *Scanner) next() {
	if s.rdoffset == len(s.src) {
		s.offset = len(s.src)
		s.ch = -1
	} else {
		s.offset = s.rdoffset
		r, w := utf8.DecodeRune(s.src[s.rdoffset:])
		if r == utf8.RuneError {
			s.errorf("illegal utf8 encoding")
		}
		s.rdoffset += w
		s.ch = r
	}
}

func (s *Scanner) skipWhitespaces() {
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' || s.ch == '\r' {
		s.next()
	}
}

func (s *Scanner) errorf(format string, a ...interface{}) {
	if s.err != nil {
		s.err = fmt.Errorf(format, a)
	}
}

func isLetter(ch rune) bool {
	if 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch >= 0x80 && unicode.IsLetter(ch) {
		return true
	}
	return false
}

func isDigit(ch rune) bool {
	if '0' <= ch && ch <= '9' || ch >= 0x80 && unicode.IsDigit(ch) {
		return true
	}
	return false
}
