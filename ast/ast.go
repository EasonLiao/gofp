package ast

import (
	"github.com/easonliao/gofp/token"
)

type Expr interface {
	exprNode()
	ExprName() string
}

type (
	IdentExpr struct {
		Name string
	}

	NumExpr struct {
		Value float64
	}

	DefExpr struct {
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
)

func (*IdentExpr) exprNode() {}
func (*NumExpr) exprNode()   {}
func (*DefExpr) exprNode()   {}
func (*FuncExpr) exprNode()  {}
func (*ExprList) exprNode()  {}
func (*CallExpr) exprNode()  {}
func (*DoExpr) exprNode()    {}
func (*IfExpr) exprNode()    {}
func (*BinaryOp) exprNode()  {}
func (*MultiOp) exprNode()   {}

func (*IdentExpr) ExprName() string { return "IdentExpr" }
func (*NumExpr) ExprName() string   { return "NumExpr" }
func (*DefExpr) ExprName() string   { return "DefExpr" }
func (*FuncExpr) ExprName() string  { return "FuncExpr" }
func (*ExprList) ExprName() string  { return "ExprList" }
func (*CallExpr) ExprName() string  { return "CallExpr" }
func (*DoExpr) ExprName() string    { return "DoExpr" }
func (*IfExpr) ExprName() string    { return "IfExpr" }
func (*BinaryOp) ExprName() string  { return "BinaryOP" }
func (*MultiOp) ExprName() string   { return "MultiOp" }
