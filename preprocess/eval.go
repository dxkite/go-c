package preprocess

import (
	"dxkite.cn/c/errors"
	"dxkite.cn/c/token"
	"strconv"
)

func Eval(ctx *Context, expr Expr) bool {
	e := &Evaluator{ctx: ctx, err: errors.ErrorList{}}
	return e.eval(expr) > 0
}

type Evaluator struct {
	ctx *Context
	err errors.ErrorList
}

// 获取数据类型
func (e *Evaluator) valueOf(tok token.Token) int64 {
	switch tok.Type() {
	case token.INT:
		return e.evalInt(tok)
	case token.CHAR:
		return e.evalChar(tok)
	case token.IDENT:
		if e.ctx.IsDefined(tok.Literal()) {
			return 1
		}
	}
	return 0
}

// 一元运算
func (e *Evaluator) evalUnaryExpr(expr *UnaryExpr) int64 {
	v := e.eval(expr.X)
	switch expr.Op.Literal() {
	case "~":
		return ^v
	case "!":
		if !(v > 0) {
			return 1
		}
	case "-":
		return -v
	case "+":
		return +v
	case "defined":
		return v
	}
	return 0
}

// 一元运算
func (e *Evaluator) evalBinaryExpr(expr *BinaryExpr) int64 {
	x, y := e.eval(expr.X), e.eval(expr.Y)
	switch expr.Op.Literal() {
	case "||":
		if (x > 0) || (y > 0) {
			return 1
		}
	case "&&":
		if (x > 0) && (y > 0) {
			return 1
		}
	case "|":
		return x | y
	case "^":
		return x ^ y
	case "&":
		return x & y
	case "!=":
		if x != y {
			return 1
		}
	case "==":
		if x == y {
			return 1
		}
	case ">":
		if x > y {
			return 1
		}
	case "<":
		if x < y {
			return 1
		}
	case ">=":
		if x >= y {
			return 1
		}
	case "<=":
		if x <= y {
			return 1
		}
	case ">>":
		return x >> y
	case "<<":
		return x << y
	case "+":
		return x + y
	case "-":
		return x - y
	case "*":
		return x * y
	case "/":
		return x / y
	case "%":
		return x % y
	}
	return 0
}

// 解析数字
func (e *Evaluator) evalInt(tok token.Token) int64 {
	v, err := strconv.ParseInt(tok.Literal(), 0, 64)
	if err != nil {
		e.addErr(tok.Position(), "error parse int %s", err.Error())
	}
	return v
}

// 解析数字（浮点数）
func (e *Evaluator) evalChar(tok token.Token) int64 {
	if c, ok := parseChar(tok.Literal()); !ok {
		e.addErr(tok.Position(), "error parse char %s", tok.Literal())
	} else {
		return int64(c)
	}
	return 0
}

func parseChar(ch string) (uint8, bool) {
	l := len(ch)
	if l <= 2 {
		return 0, false
	}
	if ch[0] != '\'' && ch[l-1] != '\'' {
		return 0, false
	}
	ch = ch[1 : l-1] // ''
	// x
	if len(ch) == 1 {
		return ch[0], true
	} else if ch[0] == '\\' {
		// \000
		// \xff
		i := 1
		base := 8
		switch ch[1] {
		case 'x':
			i = 2
			base = 16
		case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', '"', '\'':
			return ch[1], true
		}
		x := 0
		for i < len(ch) {
			x = x*base + digitVal(rune(ch[i]))
			i++
		}
		return uint8(x), true
	}
	return 0, false
}

func lower(ch rune) rune { return ('a' - 'A') | ch }
func digitVal(ch rune) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch - '0')
	case 'a' <= lower(ch) && lower(ch) <= 'f':
		return int(lower(ch) - 'a' + 10)
	}
	return 16
}

func (e *Evaluator) addErr(pos token.Position, msg string, args ...interface{}) {
	e.err.Add(pos, msg, args...)
}

func (e *Evaluator) eval(expr Expr) int64 {
	switch v := expr.(type) {
	case *IdentLit:
		return e.valueOf(v.Token)
	case *IntLit:
		return e.valueOf(v.Token)
	case *UnaryExpr:
		return e.evalUnaryExpr(v)
	case *BinaryExpr:
		return e.evalBinaryExpr(v)
	case *CondExpr:
		x := e.eval(v.X)
		if x > 0 {
			return e.eval(v.Then)
		}
		return e.eval(v.Else)
	case *ParenExpr:
		return e.eval(v.X)
	}
	return 0
}
