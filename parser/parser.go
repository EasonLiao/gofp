package parser

import (
	"fmt"
	"strconv"

	"github.com/easonliao/gofp/ast"
	"github.com/easonliao/gofp/scanner"
	"github.com/easonliao/gofp/token"
)

func ParseExpr(src []byte) (ast.Expr, error) {
	var p parser
	p.init(src)
	expr := p.parseExpr()
	p.match(token.EOF)
	return expr, p.err
}

type parser struct {
	sc  scanner.Scanner
	tok token.Token
	lit string
	err error
}

func (p *parser) init(src []byte) {
	p.sc.Init(src)
	p.next()
}

func (p *parser) parseExpr() ast.Expr {
	if p.err != nil {
		return nil
	}
	if p.tok == token.LPAREN {
		// If an expression starts with '(' it's a function call unless the first token after '(' is
		// a keyword like 'if', 'fn', 'do', 'def'.
		defer p.match(token.RPAREN)
		p.next()
		switch p.tok {
		case token.FN:
			return p.parseFun()
		case token.IF:
			return p.parseIf()
		case token.DO:
			return p.parseDoBlock()
		case token.DEF:
			return p.parseDef()
		case token.DEFN:
			return p.parseDefn()
		case token.LET:
			return p.parseLet()
		case token.ADD, token.SUB, token.MULT, token.DIV:
			return p.parseMultiOp()
		case token.LT, token.GT, token.LE, token.GE, token.EQ:
			return p.parseBinaryOp()
		case token.IDENT, token.LPAREN:
			// It's a function call.
			return p.parseCallExpr()
		}
	} else {
		// The first token of an expression is not '(', it can only be num or identifier.
		switch p.tok {
		case token.NUM:
			return p.parseNum()
		case token.IDENT:
			return p.parseIdent()
		case token.TRUE:
			p.next()
			return &ast.BooleanExpr{Bool: true}
		case token.FALSE:
			p.next()
			return &ast.BooleanExpr{Bool: false}
		case token.EOF:
			return &ast.NilExpr{}
		}
	}
	p.errorf("unexpected token %s", token.TokenName(p.tok))
	return nil
}

func (p *parser) parseIdent() *ast.IdentExpr {
	if p.err != nil {
		return nil
	}
	lit := p.lit
	p.match(token.IDENT)
	return &ast.IdentExpr{Name: lit}
}

func (p *parser) parseNum() ast.Expr {
	if p.err != nil {
		return nil
	}
	lit := p.lit
	p.match(token.NUM)
	value, err := strconv.ParseFloat(lit, 64)
	if err != nil {
		p.errorf(err.Error())
		return nil
	}
	return &ast.NumExpr{Value: value}
}

func (p *parser) parseFun() ast.Expr {
	if p.err != nil {
		return nil
	}
	parameters := make([]*ast.IdentExpr, 0)
	p.match(token.FN)
	p.match(token.LBRACK)
	for p.tok == token.IDENT && p.err == nil {
		parameters = append(parameters, p.parseIdent())
	}
	p.match(token.RBRACK)
	body := p.parseExpr()
	return &ast.FuncExpr{Params: parameters, Expr: body}
}

func (p *parser) parseIf() *ast.IfExpr {
	if p.err != nil {
		return nil
	}
	p.match(token.IF)
	cond := p.parseExpr()
	then := p.parseExpr()
	else_ := p.parseExpr()
	return &ast.IfExpr{Cond: cond, Then: then, Else: else_}
}

func (p *parser) parseCallExpr() *ast.CallExpr {
	if p.err != nil {
		return nil
	}
	fun := p.parseExpr()
	args := p.parseExprList()
	return &ast.CallExpr{Fun: fun, Args: args}
}

func (p *parser) parseDoBlock() *ast.DoExpr {
	if p.err != nil {
		return nil
	}
	p.match(token.DO)
	exprs := p.parseExprList()
	return &ast.DoExpr{Exprs: exprs}
}

func (p *parser) parseExprList() *ast.ExprList {
	exprs := make([]ast.Expr, 0)
	for p.err == nil && p.canStartExpr() {
		exprs = append(exprs, p.parseExpr())
	}
	return &ast.ExprList{Exprs: exprs}
}

func (p *parser) parseDef() *ast.DefExpr {
	if p.err != nil {
		return nil
	}
	p.match(token.DEF)
	ident := p.parseIdent()
	expr := p.parseExpr()
	return &ast.DefExpr{Ident: ident, Expr: expr}
}

func (p *parser) parseDefn() *ast.DefnExpr {
	if p.err != nil {
		return nil
	}
	p.match(token.DEFN)
	ident := p.parseIdent()
	p.match(token.LBRACK)
	parameters := make([]*ast.IdentExpr, 0)
	for p.tok == token.IDENT && p.err == nil {
		parameters = append(parameters, p.parseIdent())
	}
	p.match(token.RBRACK)
	body := p.parseExpr()
	fnExpr := &ast.FuncExpr{Params: parameters, Expr: body}
	return &ast.DefnExpr{Ident: ident, Expr: fnExpr}
}

func (p *parser) parseBinaryOp() *ast.BinaryOp {
	if p.err != nil {
		return nil
	}
	op := p.tok
	p.next()
	left := p.parseExpr()
	right := p.parseExpr()
	return &ast.BinaryOp{Op: op, Left: left, Right: right}
}

func (p *parser) parseMultiOp() *ast.MultiOp {
	if p.err != nil {
		return nil
	}
	op := p.tok
	p.next()
	exprs := p.parseExprList()
	return &ast.MultiOp{Op: op, Exprs: exprs}
}

func (p *parser) next() {
	p.tok, p.lit, p.err = p.sc.Next()
}

func (p *parser) match(tok token.Token) {
	if p.err != nil {
		return
	}
	if p.tok != tok {
		p.errorf("Expecting token %s while get %s", token.TokenName(tok), token.TokenName(p.tok))
		return
	}
	p.next()
}

func (p *parser) parseLet() *ast.LetExpr {
	if p.err != nil {
		return nil
	}
	p.match(token.LET)
	p.match(token.LBRACK)
	bindings := make([]*ast.BindExpr, 0, 1)
	bindings = append(bindings, p.parseBindingPair())
	for p.err == nil && p.tok == token.IDENT {
		bindings = append(bindings, p.parseBindingPair())
	}
	p.match(token.RBRACK)
	return &ast.LetExpr{Bindings: bindings, Body: p.parseExpr()}
}

func (p *parser) parseBindingPair() *ast.BindExpr {
	ident := p.parseIdent()
	expr := p.parseExpr()
	return &ast.BindExpr{Ident: ident, Value: expr}
}

// check whether current token can be a start of an expression.
func (p *parser) canStartExpr() bool {
	if p.tok == token.LPAREN || p.tok == token.IDENT || p.tok == token.NUM {
		return true
	}
	return false
}

func (p *parser) errorf(format string, a ...interface{}) {
	p.err = fmt.Errorf(format, a...)
}
