package scanner

import (
	"testing"

	"github.com/easonliao/gofp/token"
)

func TestScanner(t *testing.T) {
	var s Scanner
	s.Init([]byte("a=1.1 b=2()[]<><=>="))
	tok, lit, _ := s.Next()
	if tok != token.IDENT || lit != "a" {
		t.Error("error")
	}
	tok, lit, _ = s.Next()
	if tok != token.EQ || lit != "" {
		t.Error("error")
	}
	tok, lit, _ = s.Next()
	if tok != token.NUM || lit != "1.1" {
		t.Error("error")
	}
	tok, lit, _ = s.Next()
	if tok != token.IDENT || lit != "b" {
		t.Error("error")
	}
	tok, lit, _ = s.Next()
	if tok != token.EQ || lit != "" {
		t.Error("error")
	}
	tok, lit, _ = s.Next()
	if tok != token.NUM || lit != "2" {
		t.Error("error")
	}
	tok, lit, _ = s.Next()
	if tok != token.LPAREN || lit != "" {
		t.Error("error")
	}
	tok, lit, _ = s.Next()
	if tok != token.RPAREN || lit != "" {
		t.Error("error")
	}
	tok, lit, _ = s.Next()
	if tok != token.LBRACK || lit != "" {
		t.Error("error")
	}
	tok, lit, _ = s.Next()
	if tok != token.RBRACK || lit != "" {
		t.Error("error")
	}
	tok, lit, _ = s.Next()
	if tok != token.LT || lit != "" {
		t.Error("error")
	}
	tok, lit, _ = s.Next()
	if tok != token.GT || lit != "" {
		t.Error("error")
	}
	tok, lit, _ = s.Next()
	if tok != token.LE || lit != "" {
		t.Error("error")
	}
	tok, lit, _ = s.Next()
	if tok != token.GE || lit != "" {
		t.Error("error")
	}
	tok, lit, _ = s.Next()
	if tok != token.EOF || lit != "" {
		t.Error("error")
	}
}
