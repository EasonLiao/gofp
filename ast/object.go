package ast

type Object struct {
	Kind  ObjKind
	Value interface{}
}

var NilObj = &Object{Kind: Nil, Value: nil}

type ObjKind int

const (
	Bad ObjKind = iota
	Double
	Func
	Boolean
	List
	Nil
)

func (o ObjKind) String() string {
	switch o {
	case Bad:
		return "Bad"
	case Double:
		return "Double"
	case Func:
		return "Function"
	case Boolean:
		return "Boolean"
	case List:
		return "List"
	case Nil:
		return "Nil"
	}
	return "UNKNOWN"
}

func createDouble(v float64) *Object {
	return &Object{Kind: Double, Value: v}
}

func createBoolean(b bool) *Object {
	return &Object{Kind: Boolean, Value: b}
}

func createList(list []*Object) *Object {
	return &Object{Kind: List, Value: list}
}

func createFunc(sc *Scope, params []string, body Expr) *Object {
	return &Object{Kind: Func, Value: &FuncValue{Closure: sc, Params: params, Body: body}}
}

type FuncValue struct {
	Closure *Scope
	Params  []string
	Body    Expr
}
