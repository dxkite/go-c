package preprocess

import (
	"dxkite.cn/c/scanner"
	"dxkite.cn/c/token"
	"fmt"
)

const (
	LowestPrec = 0 // 最低优先级
	BinaryPrec = 2
	UnaryPrec  = 12 // 一元操作符
)

func Precedence(lit string) int {
	switch lit {
	case "?":
		return 1
	case "||":
		return 2
	case "&&":
		return 3
	case "|":
		return 4
	case "^":
		return 5
	case "&":
		return 6
	case "!=", "==":
		return 7
	case ">", "<", ">=", "<=":
		return 8
	case ">>", "<<":
		return 9
	case "+", "-":
		return 10
	case "*", "/", "%":
		return 11
	}
	return LowestPrec
}

func EvalConstExpr(ctx *Context, tks []token.Token) bool {
	p := NewParser(ctx, scanner.NewArrayScan(tks))
	return Eval(ctx, p.ParseExpr())
}

type (
	BadExpr struct {
		token.Token
	}
	IdentLit struct {
		token.Token
	}
	IntLit struct {
		token.Token
	}
	// 一元运算
	UnaryExpr struct {
		Op token.Token // 操作类型
		X  Expr        // 操作的表达式
	}
	// 二元运算
	BinaryExpr struct {
		X  Expr        // 左值
		Op token.Token // 操作类型
		Y  Expr        // 右值
	}
	// 条件表达式 logical-OR-expression ? expression : conditional-expression
	CondExpr struct {
		X    Expr
		Op   token.Token // 操作类型
		Then Expr        // 左值
		Else Expr        // 右值
	}
	// 括号表达式
	ParenExpr struct {
		Lparen token.Token // "("
		X      Expr        // 表达式值
		Rparen token.Token // ")"
	}
)

type Expr interface {
	expr()
}

func (*BadExpr) expr()    {}
func (*IdentLit) expr()   {}
func (*IntLit) expr()     {}
func (*UnaryExpr) expr()  {}
func (*BinaryExpr) expr() {}
func (*CondExpr) expr()   {}
func (*ParenExpr) expr()  {}

func ExprString(expr Expr) string {
	switch e := expr.(type) {
	case *BadExpr:
		return fmt.Sprintf("bad{%s}", e.Literal())
	case *IdentLit:
		return e.Literal()
	case *IntLit:
		return e.Literal()
	case *UnaryExpr:
		return fmt.Sprintf("(%s (%s))", e.Op.Literal(), ExprString(e.X))
	case *BinaryExpr:
		return fmt.Sprintf("(%s %s %s)", ExprString(e.X), e.Op.Literal(), ExprString(e.Y))
	case *CondExpr:
		return fmt.Sprintf("(%s?%s:%s)", ExprString(e.X), ExprString(e.Then), ExprString(e.Else))
	case *ParenExpr:
		return fmt.Sprintf("(%s)", ExprString(e.X))
	}
	return "unknown{}"
}

type Parser struct {
	cur token.Token
	r   scanner.Scanner
	ctx *Context
}

func NewParser(ctx *Context, r scanner.Scanner) *Parser {
	p := &Parser{
		r: scanner.NewTokenScan(r),
	}
	p.ctx = ctx
	p.next()
	return p
}

func (p *Parser) ParseExpr() (expr Expr) {
	expr = p.parseExpr()
	if p.cur.Type() != token.EOF {
		p.addErr(p.cur.Position(), "unexpect token %s", p.cur.Literal())
	}
	return
}

func (p *Parser) addErr(pos token.Position, msg string, args ...interface{}) {
	p.ctx.AddError(pos, "preprocess expr: "+msg, args...)
}

func (p *Parser) parseCondExpr() Expr {
	x := p.parseBinaryExpr(BinaryPrec)
	if p.cur.Type() == token.PUNCTUATOR && p.cur.Literal() == "?" {
		op := p.cur
		p.next()
		then := p.parseExpr()
		p.exceptPunctuator(":")
		el := p.parseCondExpr()
		return &CondExpr{
			X:    x,
			Op:   op,
			Then: then,
			Else: el,
		}
	}
	return x
}

func (p *Parser) exceptPunctuator(lit string) (t token.Token) {
	t = p.cur
	if p.cur.Type() == token.PUNCTUATOR && p.cur.Literal() == lit {
		p.next()
		return
	}
	p.addErr(p.cur.Position(), "expect %s got %s", lit, p.cur.Literal())
	return
}

func (p *Parser) parseOpExpr(prec int) (expr Expr) {
	if prec >= UnaryPrec {
		return p.parseUnaryExpr()
	} else {
		return p.parseBinaryExpr(prec)
	}
}

func (p *Parser) parseBinaryExpr(prec int) Expr {
	x := p.parseOpExpr(prec + 1)
	if p.cur.Type() == token.PUNCTUATOR && Precedence(p.cur.Literal()) >= prec {
		op := p.cur
		p.next()
		y := p.parseOpExpr(prec + 1)
		return &BinaryExpr{
			X:  x,
			Op: op,
			Y:  y,
		}
	}
	return x
}

// 	( ("-" / "+" / "~" / "defined" ) parseTermExpr )
func (p *Parser) parseUnaryExpr() Expr {
	if p.cur.Type() == token.PUNCTUATOR && litIn(p.cur.Literal(), []string{"+", "-", "~", "!"}) {
		op := p.cur
		p.next()
		x := p.parseTermExpr()
		return &UnaryExpr{
			Op: op,
			X:  x,
		}
	}

	if p.cur.Type() == token.IDENT && p.cur.Literal() == "defined" {
		return p.parseDefined()
	}

	return p.parseTermExpr()
}

// defined ( id )
func (p *Parser) parseDefined() (expr Expr) {
	op := p.cur
	p.next()
	left := 0
	for p.cur.Type() == token.PUNCTUATOR && p.cur.Literal() == "(" {
		p.next()
		left++
	}
	if p.cur.Type() != token.IDENT {
		p.addErr(p.cur.Position(), "expected ident got %s", p.cur.Type())
	}
	x := p.cur
	p.next()
	for left > 0 {
		p.exceptPunctuator(")")
		left--
	}
	return &UnaryExpr{
		Op: op,
		X:  &IdentLit{x},
	}
}

func litIn(lit string, arr []string) bool {
	for _, v := range arr {
		if lit == v {
			return true
		}
	}
	return false
}

// "(" expr ")" | number | char
func (p *Parser) parseTermExpr() (expr Expr) {
	if p.cur.Type() == token.PUNCTUATOR && p.cur.Literal() == "(" {
		l := p.cur
		p.next()
		x := p.parseExpr()
		r := p.exceptPunctuator(")")
		return &ParenExpr{
			Lparen: l,
			X:      x,
			Rparen: r,
		}
	}

	if v := p.parseValueExpr(); v != nil {
		return v
	}

	p.addErr(p.cur.Position(), "unexpected %s %s", p.cur.Type(), p.cur.Literal())
	expr = &BadExpr{p.cur}
	p.next()
	return
}

func (p *Parser) parseValueExpr() Expr {
	switch c, t := p.cur, p.cur.Type(); {
	case t == token.IDENT:
		p.next()
		return &IdentLit{c}
	case t == token.CHAR:
		p.next()
		return &IntLit{c}
	case t == token.INT:
		p.next()
		return &IntLit{c}
	}
	return nil
}

func (p *Parser) parseExpr() Expr {
	return p.parseCondExpr()
}

// 获取下一个Token
func (p *Parser) next() token.Token {
	p.cur = p.r.Scan()
	return p.cur
}
