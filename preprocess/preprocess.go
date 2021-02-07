package preprocess

import (
	"dxkite.cn/go-c11"
	"dxkite.cn/go-c11/scanner"
	"dxkite.cn/go-c11/token"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type MacroDecl interface {
	decl()
}

type MacroVal struct {
	Name string
	Body []token.Token
}

type MacroFunc struct {
	Name     string
	Params   []string
	Ellipsis bool // ...
	Body     []token.Token
}

type HandlerFn func(tok token.Token) []token.Token

// MacroVal Handler
type MacroHandler struct {
	Name    string
	Handler HandlerFn
}

func (m *MacroVal) decl()     {}
func (m *MacroFunc) decl()    {}
func (m *MacroHandler) decl() {}

// 宏定义
type MacroInfo struct {
	Decl  MacroDecl // 定义信息
	Index int       // 定义的优先级
}

type Token struct {
	Pos token.Position
	Typ token.Type
	Lit string
	// 展开优先级
	Index int
	// 展开
	Expand token.Token
}

func (t *Token) Position() token.Position {
	return t.Pos
}

func (t *Token) Type() token.Type {
	return t.Typ
}
func (t *Token) Literal() string {
	return t.Lit
}

// 解析环境
type Context struct {
	Val     map[string]*MacroInfo // 宏定义
	Inc     []string              // 文件目录
	counter int                   // __COUNTER__
	once    map[string]struct{}   // #pragma once

	// MacroIndex
	idx int
	// 当前元素
	lit string
	pos token.Position
	tok token.Type
	cur token.Token
	in  scanner.MultiScanner
	rcd bool
	tks []token.Token
	err go_c11.ErrorList
	// 条件栈
	cdt *ConditionStack
}

// 设置输入
func New(s scanner.Scanner) *Context {
	c := &Context{}
	c.in = scanner.NewMultiScan(s)
	c.err.Reset()
	c.Val = map[string]*MacroInfo{}
	c.cdt = NewConditionStack()
	c.next()
	return c
}

// 测试 once
func (c *Context) Scan() (t token.Token) {
	for t == nil {
		if c.tok == token.EOF {
			break
		}
		if tokIsMacro(c.cur) {
			c.doMacro()
			continue
		}
		// 宏展开
		if c.doExpand(c.cur) {
			c.next()
			continue
		}
		// 普通token
		t = c.cur
		c.next()
	}
	return t
}

func (c *Context) Error() *go_c11.ErrorList {
	return &c.err
}

func tokIsMacro(tok token.Token) bool {
	return tok.Position().Column == 1 && tok.Literal() == "#"
}

// 获取下一个
func (c *Context) next() {
	if c.rcd && c.cur != nil {
		c.tks = append(c.tks, c.cur)
	}
	c.cur = c.in.Scan()
	if c.cur != nil {
		c.lit = c.cur.Literal()
		c.tok = c.cur.Type()
		c.pos = c.cur.Position()
	} else {
		c.tok = token.EOF
		c.lit = ""
	}
}

// 获取下一个非空token
func (c *Context) nextToken() {
	for {
		c.next()
		if c.tok != token.WHITESPACE {
			break
		}
	}
}

// 测试 once
func (c *Context) onceContain(p string) bool {
	_, ok := c.once[p]
	return ok
}

// 写入 pragma once
func (c *Context) pragmaOnce(p string) {
	c.once[p] = struct{}{}
}

func (c *Context) DefineVal(name, value string) error {
	if tok, err := scanner.ScanString("<build-in>", value); err != nil {
		return err
	} else {
		c.Define(name, &MacroVal{
			Name: name,
			Body: tok,
		})
	}
	return nil
}

func (c *Context) Define(name string, val MacroDecl) {
	c.idx++
	def := &MacroInfo{val, c.idx}
	c.Val[name] = def
}

func (c *Context) DefineHandler(name string, val HandlerFn) {
	c.Define(name, &MacroHandler{
		Name:    name,
		Handler: val,
	})
	return
}

func (c *Context) IsDefined(name string) bool {
	_, ok := c.Val[name]
	return ok
}

func (c *Context) Init() {
	c.DefineHandler("__FILE__", c.fileFn)
	c.DefineHandler("__LINE__", c.lineFn)
	c.DefineHandler("__COUNTER__", c.counterFn)
	_ = c.DefineVal("__DATE__", strconv.QuoteToGraphic(time.Now().Format("Jan 02 2006")))
	_ = c.DefineVal("__TIME__", strconv.QuoteToGraphic(time.Now().Format("15:04:05")))
}

func (c *Context) counterFn(tok token.Token) []token.Token {
	val := &Token{
		Pos: tok.Position(),
		Typ: token.INT,
		Lit: strconv.Itoa(c.counter),
	}
	c.counter++
	return []token.Token{val}
}

func (c *Context) fileFn(tok token.Token) []token.Token {
	val := &Token{
		Pos: tok.Position(),
		Typ: token.STRING,
		Lit: strconv.QuoteToGraphic(tok.Position().Filename),
	}
	return []token.Token{val}
}

func (c *Context) lineFn(tok token.Token) []token.Token {
	val := &Token{
		Pos: tok.Position(),
		Typ: token.STRING,
		Lit: strconv.Itoa(tok.Position().Line),
	}
	return []token.Token{val}
}

func (c *Context) doMacro() {
	c.nextToken()
	switch c.lit {
	case "if":
		c.next()
		cdt := c.evalConstExpr()
		c.expectEndMacro()
		if cdt {
			c.cdt.Push(IN_THEN)
		} else {
			// 跳到下一个分支
			c.cdt.Push(IN_ELSE)
			c.skipUtilElse()
		}
	case "ifdef":
		c.doIfDefine(true)
	case "ifndef":
		c.doIfDefine(false)
	case "elif":
		c.next()
		if c.cdt.Top() == IN_ELSE {
			cdt := c.evalConstExpr()
			c.expectEndMacro()
			if cdt {
				c.cdt.Pop()
				c.cdt.Push(IN_THEN)
			} else {
				// 跳到下一个分支
				c.skipUtilElse()
			}
		} else if c.cdt.Top() == IN_THEN {
			// 直接跳到结尾
			c.skipUtilCdt("endif")
			c.next() // endif
			c.expectEndMacro()
		} else {
			c.err.Add(c.pos, "unexpected #else")
		}
	case "else":
		c.next()
		c.expectEndMacro()
		if c.cdt.Top() == IN_THEN {
			c.skipUtilCdt("endif")
			c.next()
			c.expectEndMacro()
		} else {
			c.err.Add(c.pos, "unexpected #else")
		}
	case "endif":
		if c.cdt.Top() == IN_THEN || c.cdt.Top() == IN_ELSE {
			c.cdt.Pop()
			c.next() // endif
			c.expectEndMacro()
		} else {
			c.err.Add(c.pos, "unexpected #endif")
		}
	case "define":
		c.doDefine()
	case "include":
	case "pragma":
	case "line":
	case "error":
	default:
	}
}

// 展开宏
func (c *Context) doExpand(tok token.Token) bool {
	if tok.Type() != token.IDENT {
		return false
	}
	if tks, ok := c.Expand(tok); ok {
		c.push(tks)
		return true
	}
	return false
}

// 展开宏
func (c *Context) push(tok []token.Token) {
	c.in.Push(scanner.NewArrayScan(tok))
}

// 处理宏展开
func (c *Context) Expand(tok token.Token) ([]token.Token, bool) {
	if tok.Type() != token.IDENT {
		return nil, false
	}
	name := tok.Literal()
	if v, ok := c.Val[name]; ok {
		// 宏未定义时不展开
		if t, ok := tok.(*Token); ok && v.Index > t.Index {
			return nil, false
		}
		switch val := v.Decl.(type) {
		case *MacroVal:
			return c.ExpandVal(tok, v.Index, val.Body), true
		case *MacroHandler:
			return c.ExpandVal(tok, v.Index, val.Handler(tok)), true
		case *MacroFunc:
			// 忽略非函数式宏
			if n := c.peekNext(); n.Literal() != "(" {
				return nil, false
			}
			// 处理
		}
	}
	return nil, false
}

// 展开宏
func (c *Context) ExpandVal(v token.Token, idx int, tks []token.Token) []token.Token {
	ex := make([]token.Token, len(tks))
	col := 0
	pos := v.Position()
	for i, tok := range tks {
		if i == 0 {
			col = tok.Position().Column
		}
		t := &Token{
			Pos: token.Position{
				Filename: pos.Filename,
				Line:     pos.Line,
				Column:   pos.Column + tok.Position().Column - col,
			},
			Typ:    tok.Type(),
			Lit:    tok.Literal(),
			Index:  idx,
			Expand: v,
		}
		ex[i] = t
	}
	return ex
}

// peek 下一个非空 token
func (c *Context) peekNext() token.Token {
	n := 1
	for {
		v := c.peek(n)
		if len(v) < n {
			break
		}
		if v[n-1].Type() != token.WHITESPACE {
			return v[n-1]
		}
	}
	return &Token{
		Pos: token.Position{},
		Typ: token.EOF,
		Lit: "",
	}
}

// peek 下一个 token
func (c *Context) peekOne() token.Token {
	if ps, ok := c.in.(scanner.PeekScanner); ok {
		return ps.PeekOne()
	}
	p := scanner.NewPeekScan(c.in)
	tok := p.PeekOne()
	c.in = scanner.NewMultiScan(p)
	return tok
}

func (c *Context) peek(offset int) []token.Token {
	if ps, ok := c.in.(scanner.PeekScanner); ok {
		return ps.Peek(offset)
	}
	p := scanner.NewPeekScan(c.in)
	tok := p.Peek(offset)
	c.in = scanner.NewMultiScan(p)
	return tok
}

func (c *Context) expectIdent() string {
	if c.tok == token.IDENT {
		lit := c.lit
		c.next()
		return lit
	}
	c.err.Add(c.pos, fmt.Sprintf("expect token %s got %s", token.IDENT, c.tok))
	return ""
}

func (c *Context) expectEndMacro() {
	if c.isMacroEnd() {
		c.nextToken()
		return
	}
	c.err.Add(c.pos, fmt.Sprintf("expect end macro got %s", c.tok))
}

// 宏结尾
func (c *Context) isMacroEnd() bool {
	return c.tok == token.NEWLINE || c.tok == token.EOF
}

// 跳过无法到达的代码
func (c *Context) skipUtilCdt(names ...string) []token.Token {
	cdt := 0
	tks := make([]token.Token, 2)
	for {
		c.next()
		if c.tok == token.EOF {
			c.err.Add(c.pos, fmt.Sprintf("expect %s, got EOF", strings.Join(names, ",")))
			break
		}
		if tokIsMacro(c.cur) {
			tks[0] = c.cur
			c.nextToken()
			switch c.lit {
			case "if", "ifndef", "ifdef":
				cdt++
			default:
				if cdt == 0 {
					for _, name := range names {
						if name == c.lit {
							tks[1] = c.cur
							return tks
						}
					}
				}
				if c.lit == "endif" {
					cdt--
				}
			}
		}
	}
	return tks[0:0]
}

// #ifdef #ifndef
func (c *Context) doIfDefine(want bool) {
	c.nextToken()
	ident := c.expectIdent()
	cdt := c.IsDefined(ident)
	if cdt == want {
		c.cdt.Push(IN_THEN)
	} else {
		c.cdt.Push(IN_ELSE)
		c.skipUtilElse()
	}
	c.expectEndMacro()
}

// skip to #else/#elif
func (c *Context) skipUtilElse() {
	m := c.skipUtilCdt("elif", "else")
	if c.lit == "elif" {
		c.next()  // elif
		c.push(m) // push back
	} else {
		c.next() // else
	}
}

func (c *Context) evalConstExpr() bool {
	tks := []token.Token{}
	for {
		if c.tok == token.EOF || c.tok == token.NEWLINE {
			break
		}
		if c.tok != token.WHITESPACE {
			tks = append(tks, c.cur)
		}
		c.next()
	}
	// TODO parse expr
	return false
}

func (c *Context) doDefine() {
	c.nextToken()
	ident := c.expectIdent()
	if n := c.peekOne(); n.Literal() == "(" {
		c.doDefineFunc(ident)
	} else {
		c.doDefineVal(ident)
	}
	c.expectEndMacro()
}

func (c *Context) doDefineVal(ident string) {
	tks := make([]token.Token, 0)
	c.nextToken()
	for !c.isMacroEnd() {
		tks = append(tks, c.cur)
		c.nextToken()
	}
	c.Define(ident, &MacroVal{
		Name: ident,
		Body: tks,
	})
}

func (c *Context) doDefineFunc(ident string) {

}
