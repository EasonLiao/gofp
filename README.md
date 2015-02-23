# gofp

A toy interpreter for a functional programming language similar to Clojure/Lisp.

example:
```
> (defn accum [n] (if (= n 0) 0 (+ n (accum (- n 1)))))

ast.DefnExpr {
.  Ident: ast.IdentExpr {
.  .  Name: "accum"
.  }
.  Expr: ast.FuncExpr {
.  .  Params: []*ast.IdentExpr (len = 1) {
.  .  .  0: ast.IdentExpr {
.  .  .  .  Name: "n"
.  .  .  }
.  .  }
.  .  Expr: ast.IfExpr {
.  .  .  Cond: ast.BinaryOp {
.  .  .  .  Op: "="
.  .  .  .  Left: ast.IdentExpr {
.  .  .  .  .  Name: "n"
.  .  .  .  }
.  .  .  .  Right: ast.NumExpr {
.  .  .  .  .  Value: 0.000000
.  .  .  .  }
.  .  .  }
.  .  .  Then: ast.NumExpr {
.  .  .  .  Value: 0.000000
.  .  .  }
.  .  .  Else: ast.MultiOp {
.  .  .  .  Op: "+"
.  .  .  .  Exprs: ast.ExprList {
.  .  .  .  .  Exprs: []ast.Expr (len = 2) {
.  .  .  .  .  .  0: ast.IdentExpr {
.  .  .  .  .  .  .  Name: "n"
.  .  .  .  .  .  }
.  .  .  .  .  .  1: ast.CallExpr {
.  .  .  .  .  .  .  Fun: ast.IdentExpr {
.  .  .  .  .  .  .  .  Name: "accum"
.  .  .  .  .  .  .  }
.  .  .  .  .  .  .  Args: ast.ExprList {
.  .  .  .  .  .  .  .  Exprs: []ast.Expr (len = 1) {
.  .  .  .  .  .  .  .  .  0: ast.MultiOp {
.  .  .  .  .  .  .  .  .  .  Op: "-"
.  .  .  .  .  .  .  .  .  .  Exprs: ast.ExprList {
.  .  .  .  .  .  .  .  .  .  .  Exprs: []ast.Expr (len = 2) {
.  .  .  .  .  .  .  .  .  .  .  .  0: ast.IdentExpr {
.  .  .  .  .  .  .  .  .  .  .  .  .  Name: "n"
.  .  .  .  .  .  .  .  .  .  .  .  }
.  .  .  .  .  .  .  .  .  .  .  .  1: ast.NumExpr {
.  .  .  .  .  .  .  .  .  .  .  .  .  Value: 1.000000
.  .  .  .  .  .  .  .  .  .  .  .  }
.  .  .  .  .  .  .  .  .  .  .  }
.  .  .  .  .  .  .  .  .  .  }
.  .  .  .  .  .  .  .  .  }
.  .  .  .  .  .  .  .  }
.  .  .  .  .  .  .  }
.  .  .  .  .  .  }
.  .  .  .  .  }
.  .  .  .  }
.  .  .  }
.  .  }
.  }
}
&{Nil <nil>}

> (accum 100)

ast.CallExpr {
.  Fun: ast.IdentExpr {
.  .  Name: "accum"
.  }
.  Args: ast.ExprList {
.  .  Exprs: []ast.Expr (len = 1) {
.  .  .  0: ast.NumExpr {
.  .  .  .  Value: 100.000000
.  .  .  }
.  .  }
.  }
}
&{Double 5050}
```
