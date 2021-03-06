package ast

import (
	"fmt"

	"github.com/easonliao/gofp/token"
)

type Expr interface {
	Eval(sc *Scope) (*Object, error)
	// Given the scope, collects all the unresolved identifiers in the expression, used for closure capture.
	collectUnresolvedNames(*Scope, map[string]bool)
}

type (
	NilExpr struct {
	}

	IdentExpr struct {
		Name string
	}

	NumExpr struct {
		Value float64
	}

	BooleanExpr struct {
		Bool bool
	}

	DefExpr struct {
		Ident *IdentExpr
		Expr  Expr
	}

	DefnExpr struct {
		Ident *IdentExpr
		Expr  Expr
	}

	FuncExpr struct {
		Params []*IdentExpr
		Expr   Expr
	}

	ExprList struct {
		Exprs []Expr
	}

	CallExpr struct {
		// Fun is an expression returns a function object.
		Fun  Expr
		Args *ExprList
	}

	DoExpr struct {
		Exprs *ExprList
	}

	IfExpr struct {
		Cond Expr
		Then Expr
		Else Expr
	}

	BinaryOp struct {
		Op    token.Token
		Left  Expr
		Right Expr
	}

	MultiOp struct {
		Op    token.Token
		Exprs *ExprList
	}

	BindExpr struct {
		Ident *IdentExpr
		Value Expr
	}

	LetExpr struct {
		Bindings []*BindExpr
		Body     Expr
	}
)

func (*NilExpr) Eval(sc *Scope) (*Object, error) {
	return NilObj, nil
}

// Eval implementation.
func (expr *IdentExpr) Eval(sc *Scope) (*Object, error) {
	obj := sc.Lookup(expr.Name)
	if obj == nil {
		return nil, fmt.Errorf("%q is not defined.", expr.Name)
	}
	return obj, nil
}

func (expr *NumExpr) Eval(sc *Scope) (*Object, error) {
	return createDouble(expr.Value), nil
}

func (expr *BooleanExpr) Eval(sc *Scope) (*Object, error) {
	return createBoolean(expr.Bool), nil
}

func (expr *DefExpr) Eval(sc *Scope) (*Object, error) {
	obj, err := expr.Expr.Eval(sc)
	if err != nil {
		return nil, err
	}
	if obj == NilObj {
		return nil, fmt.Errorf("Can't bind nil object to symbol.")
	}
	// Put it into symbol table.
	sc.Insert(expr.Ident.Name, obj)
	return NilObj, nil
}

func (expr *DefnExpr) Eval(sc *Scope) (*Object, error) {
	newScope := NewScope(sc)
	// Since defn binds a name to a function which can call itself, binds the fun name with Self obj first.
	newScope.Insert(expr.Ident.Name, SelfObj)
	obj, err := expr.Expr.Eval(newScope)
	if err != nil {
		return nil, err
	}
	if obj == NilObj {
		return nil, fmt.Errorf("Can't bind nil object to symbol.")
	}
	// Put it into symbol table.
	sc.Insert(expr.Ident.Name, obj)
	return NilObj, nil
}

func (expr *FuncExpr) Eval(sc *Scope) (*Object, error) {
	params := make([]string, 0, len(expr.Params))
	for _, ident := range expr.Params {
		params = append(params, ident.Name)
	}
	unresolvedNames := make(map[string]bool)
	expr.collectUnresolvedNames(NewScope(nil), unresolvedNames)
	closure := NewScope(nil)
	funcObj := createFunc(closure, params, expr.Expr)
	// Capture all the unresolved names from current scope.
	for name, _ := range unresolvedNames {
		obj := sc.Lookup(name)
		if obj == nil {
			return nil, fmt.Errorf("uncaptured name %q", name)
		}
		if obj == SelfObj {
			closure.Insert(name, funcObj)
		} else {
			closure.Insert(name, obj)
		}
	}
	return funcObj, nil
}

func (expr *ExprList) Eval(sc *Scope) (*Object, error) {
	objects := make([]*Object, 0, len(expr.Exprs))
	for _, e := range expr.Exprs {
		obj, err := e.Eval(sc)
		if err != nil {
			return nil, err
		}
		objects = append(objects, obj)
	}
	return createList(objects), nil
}

func (expr *CallExpr) Eval(sc *Scope) (*Object, error) {
	obj, err := expr.Fun.Eval(sc)
	if err != nil {
		return nil, err
	}
	if obj.Kind != Func {
		return nil, fmt.Errorf("The object is not a function object.")
	}
	funObj := obj.Value.(*FuncValue)
	numParams := len(funObj.Params)
	numArgs := len(expr.Args.Exprs)
	if numParams != numArgs {
		return nil, fmt.Errorf("Wrong number of arguments(%d), expect %d", numArgs, numParams)
	}
	argList, err := expr.Args.Eval(sc)
	if err != nil {
		return nil, err
	}
	args := argList.Value.([]*Object)
	// Binding arguments to function's closure.
	for idx, param := range funObj.Params {
		funObj.Closure.Insert(param, args[idx])
	}
	return funObj.Body.Eval(funObj.Closure)
}

func (expr *DoExpr) Eval(sc *Scope) (*Object, error) {
	obj, err := expr.Exprs.Eval(sc)
	if err != nil {
		return nil, err
	}
	objects := obj.Value.([]*Object)
	if len(objects) == 0 {
		return NilObj, nil
	}
	return objects[len(objects)-1], nil
}

func (expr *IfExpr) Eval(sc *Scope) (*Object, error) {
	cond, err := expr.Cond.Eval(sc)
	if err != nil {
		return nil, err
	}
	if cond.Kind != Boolean {
		return nil, fmt.Errorf("expression in if must return boolean")
	}
	condRes := cond.Value.(bool)
	if condRes {
		return expr.Then.Eval(sc)
	} else {
		return expr.Else.Eval(sc)
	}
}

func (expr *BinaryOp) Eval(sc *Scope) (*Object, error) {
	left, err := expr.Left.Eval(sc)
	if err != nil {
		return nil, err
	}
	right, err := expr.Right.Eval(sc)
	if err != nil {
		return nil, err
	}
	if left.Kind != right.Kind {
		return nil, fmt.Errorf("left operand and right operand have different types.")
	}
	if left.Kind != Double {
		return nil, fmt.Errorf("You can only compare double numbers.")
	}
	v1 := left.Value.(float64)
	v2 := right.Value.(float64)
	switch expr.Op {
	case token.LT:
		return createBoolean(v1 < v2), nil
	case token.LE:
		return createBoolean(v1 <= v2), nil
	case token.GT:
		return createBoolean(v1 > v2), nil
	case token.GE:
		return createBoolean(v1 >= v2), nil
	case token.EQ:
		return createBoolean(v1 == v2), nil
	}
	return nil, fmt.Errorf("invalid op %q", token.TokenName(expr.Op))
}

func (expr *MultiOp) Eval(sc *Scope) (*Object, error) {
	list, err := expr.Exprs.Eval(sc)
	if err != nil {
		return nil, err
	}
	objects := list.Value.([]*Object)
	if len(objects) == 0 {
		return NilObj, nil
	}
	var opFun func(v1, v2 float64) float64
	switch expr.Op {
	case token.ADD:
		opFun = func(v1, v2 float64) float64 { return v1 + v2 }
	case token.MULT:
		opFun = func(v1, v2 float64) float64 { return v1 * v2 }
	case token.DIV:
		opFun = func(v1, v2 float64) float64 { return v1 / v2 }
	case token.SUB:
		opFun = func(v1, v2 float64) float64 { return v1 - v2 }
	}
	if objects[0].Kind != Double {
		return nil, fmt.Errorf("operand must be double numbers")
	}
	operand := objects[0].Value.(float64)
	for _, obj := range objects[1:] {
		if obj.Kind != Double {
			return nil, fmt.Errorf("operand must be double numbers")
		}
		v := obj.Value.(float64)
		operand = opFun(operand, v)
	}
	return createDouble(operand), nil
}

func (expr *BindExpr) Eval(sc *Scope) (*Object, error) {
	obj, err := expr.Value.Eval(sc)
	if err != nil {
		return nil, err
	}
	sc.Insert(expr.Ident.Name, obj)
	return NilObj, nil
}

func (expr *LetExpr) Eval(sc *Scope) (*Object, error) {
	newScope := NewScope(sc)
	for _, binding := range expr.Bindings {
		_, err := binding.Eval(newScope)
		if err != nil {
			return nil, err
		}
	}
	return expr.Body.Eval(newScope)
}

func (expr *NilExpr) collectUnresolvedNames(sc *Scope, names map[string]bool) {
	panic("why tring to find unresolved names in nil expression")
}

// collectUnresolvedNames implementation.
func (expr *IdentExpr) collectUnresolvedNames(sc *Scope, names map[string]bool) {
	if sc.Lookup(expr.Name) == nil {
		// if the name doesn't exit in current scope, add to unresolved set.
		names[expr.Name] = true
	}
}

func (expr *NumExpr) collectUnresolvedNames(sc *Scope, names map[string]bool) {
	// Do nothing.
}

func (expr *BooleanExpr) collectUnresolvedNames(sc *Scope, names map[string]bool) {
	// Do nothing.
}

func (expr *DefExpr) collectUnresolvedNames(sc *Scope, names map[string]bool) {
	// def should be only allowed to in outer-most scope, it shouldn't be called.
	panic("why trying to find unresolved names in def expression.")
}

func (expr *DefnExpr) collectUnresolvedNames(sc *Scope, names map[string]bool) {
	// defn should be only allowed to in outer-most scope, it shouldn't be called.
	panic("why trying to find unresolved names in def expression.")
}

func (expr *FuncExpr) collectUnresolvedNames(sc *Scope, names map[string]bool) {
	newScope := NewScope(sc)
	for _, param := range expr.Params {
		// Add parameter names to scope, these don't need to be marked as unresolved inside of function.
		newScope.Insert(param.Name, NilObj)
	}
	expr.Expr.collectUnresolvedNames(newScope, names)
}

func (expr *ExprList) collectUnresolvedNames(sc *Scope, names map[string]bool) {
	for _, expr := range expr.Exprs {
		expr.collectUnresolvedNames(sc, names)
	}
}

func (expr *CallExpr) collectUnresolvedNames(sc *Scope, names map[string]bool) {
	expr.Fun.collectUnresolvedNames(sc, names)
	expr.Args.collectUnresolvedNames(sc, names)
}

func (expr *DoExpr) collectUnresolvedNames(sc *Scope, names map[string]bool) {
	expr.Exprs.collectUnresolvedNames(sc, names)
}

func (expr *IfExpr) collectUnresolvedNames(sc *Scope, names map[string]bool) {
	expr.Cond.collectUnresolvedNames(sc, names)
	expr.Then.collectUnresolvedNames(sc, names)
	expr.Else.collectUnresolvedNames(sc, names)
}

func (expr *BinaryOp) collectUnresolvedNames(sc *Scope, names map[string]bool) {
	expr.Left.collectUnresolvedNames(sc, names)
	expr.Right.collectUnresolvedNames(sc, names)
}

func (expr *MultiOp) collectUnresolvedNames(sc *Scope, names map[string]bool) {
	expr.Exprs.collectUnresolvedNames(sc, names)
}

func (expr *BindExpr) collectUnresolvedNames(sc *Scope, names map[string]bool) {
	expr.Value.collectUnresolvedNames(sc, names)
	sc.Insert(expr.Ident.Name, NilObj)
}

func (expr *LetExpr) collectUnresolvedNames(sc *Scope, names map[string]bool) {
	newScope := NewScope(sc)
	for _, binding := range expr.Bindings {
		binding.collectUnresolvedNames(newScope, names)
	}
	expr.Body.collectUnresolvedNames(newScope, names)
}
