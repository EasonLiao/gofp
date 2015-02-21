package parser

import (
	"testing"
)

func TestParser(t *testing.T) {
	_, err := ParseExpr([]byte("a"))
	if err != nil {
		t.Error("error")
	}
	_, err = ParseExpr([]byte("(+ 1 2)"))
	if err != nil {
		t.Error("error")
	}
	_, err = ParseExpr([]byte("(def a 1)"))
	if err != nil {
		t.Error("error")
	}
	_, err = ParseExpr([]byte("(if (< 1 2) 1 2)"))
	if err != nil {
		t.Error("error")
	}
	_, err = ParseExpr([]byte("(do (add 1 2) (sub 1 2))"))
	if err != nil {
		t.Error("error")
	}
	_, err = ParseExpr([]byte("+ 1 2"))
	if err == nil {
		t.Error("error")
	}
	_, err = ParseExpr([]byte("(+ 1 2"))
	if err == nil {
		t.Error("error")
	}
	_, err = ParseExpr([]byte("(< 1)"))
	if err == nil {
		t.Error("error")
	}
}
