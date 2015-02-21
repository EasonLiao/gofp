package ast

import (
	"testing"
)

func TestScope(t *testing.T) {
	var sc *Scope = NewScope(nil)
	sc.Insert("foo", createBoolean(true))
	var scNew = NewScope(sc)
	if scNew.Lookup("foo") == nil {
		t.Errorf("error")
	}
	scNew.Insert("foo", NilObj)
	if scNew.Lookup("foo") != NilObj {
		t.Errorf("error")
	}
	if sc.Lookup("foo") == NilObj {
		t.Errorf("error")
	}
	if sc.Lookup("bar") != nil {
		t.Errorf("error")
	}
}
