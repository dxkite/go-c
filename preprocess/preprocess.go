package preprocess

import (
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
	Body []*token.Token
}

type MacroFunc struct {
	Name     string
	Params   []string
	Ellipsis bool // ...
	Body     []*token.Token
}

type HandlerFn func(tok *token.Token) []*token.Token

// MacroVal Handler
type MacroHandler struct {
	Name    string
	Handler HandlerFn
}

func (m *MacroVal) decl()     {}
func (m *MacroFunc) decl()    {}
func (m *MacroHandler) decl() {}

// 解析环境
type Context struct {
	Val     map[string]MacroDecl // 宏定义
	Inc     []string             // 文件目录
	counter int                  // __COUNTER__
	once    map[string]struct{}  // #pragma once

	// 当前元素
	lit string
	pos token.Position
	tok token.Type
	cur *token.Token
	in  scanner.MultiScanner
	rcd bool
	tks []*token.Token
	err scanner.ErrorList
	// 条件栈
	cdt *ConditionStack
}

// 设置输入
func New(s scanner.Scanner) *Context {
	c := &Context{}
	c.in = scanner.NewMultiScan(s)
	c.err.Reset()
	c.Val = map[string]MacroDecl{}
	c.cdt = NewConditionStack()
	c.next()
	return c
}

// 测试 once
func (c *Context) Scan() (t *token.Token) {
	for t == nil {
		if c.tok == token.EOF {
			break
		}
		if tokIsMacro(c.cur) {
			c.doMacro()
			continue
		}
		// 宏展开
		if c.doExtract(c.cur) {
			c.next()
			continue
		}
		// 普通token
		t = c.cur
		c.next()
	}
	return t
}

func (c *Context) Error() *scanner.ErrorList {
	return &c.err
}

func tokIsMacro(tok *token.Token) bool {
	return tok.Position.Column == 1 && tok.Lit == "#"
}

// 获取下一个
func (c *Context) next() {
	if c.rcd && c.cur != nil {
		c.tks = append(c.tks, c.cur)
	}
	c.cur = c.in.Scan()
	if c.cur != nil {
		c.lit = c.cur.Lit
		c.tok = c.cur.Type
		c.pos = c.cur.Position
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
		c.Val[name] = &MacroVal{
			Name: name,
			Body: tok,
		}
	}
	return nil
}

func (c *Context) Define(name string, val MacroDecl) error {
	c.Val[name] = val
	return nil
}

func (c *Context) DefineHandler(name string, val HandlerFn) error {
	c.Val[name] = &MacroHandler{
		Name:    name,
		Handler: val,
	}
	return nil
}

func (c *Context) IsDefined(name string) bool {
	_, ok := c.Val[name]
	return ok
}

func (p *Context) Init() {
	_ = p.DefineHandler("__FILE__", p.fileFn)
	_ = p.DefineHandler("__LINE__", p.lineFn)
	_ = p.DefineHandler("__COUNTER__", p.counterFn)
	_ = p.DefineVal("__DATE__", strconv.QuoteToGraphic(time.Now().Format("Jan 02 2006")))
	_ = p.DefineVal("__TIME__", strconv.QuoteToGraphic(time.Now().Format("15:04:05")))
}

func (p *Context) counterFn(tok *token.Token) []*token.Token {
	val := &token.Token{
		Position: tok.Position,
		Type:     token.INT,
		Lit:      strconv.Itoa(p.counter),
	}
	p.counter++
	return []*token.Token{val}
}

func (p *Context) fileFn(tok *token.Token) []*token.Token {
	val := &token.Token{
		Position: tok.Position,
		Type:     token.STRING,
		Lit:      strconv.QuoteToGraphic(tok.Position.Filename),
	}
	return []*token.Token{val}
}

func (p *Context) lineFn(tok *token.Token) []*token.Token {
	val := &token.Token{
		Position: tok.Position,
		Type:     token.STRING,
		Lit:      strconv.Itoa(tok.Position.Line),
	}
	return []*token.Token{val}
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
	case "include":
	case "define":
	case "pragma":
	case "line":
	case "error":
	default:
	}
}

// 展开宏
func (c *Context) doExtract(tok *token.Token) bool {
	if tok.Type != token.IDENT {
		return false
	}
	if tks, ok := c.Extract(tok); ok {
		c.push(tks)
		return true
	}
	return false
}

// 展开宏
func (c *Context) push(tok []*token.Token) {
	c.in.Push(scanner.NewArrayScan(tok))
}

// 处理宏展开
func (c *Context) Extract(tok *token.Token) ([]*token.Token, bool) {
	if v, ok := c.Val[tok.Lit]; ok {
		switch val := v.(type) {
		case *MacroVal:
			return c.ExtractAt(tok.Position, val.Body), true
		case *MacroHandler:
			return c.ExtractAt(tok.Position, val.Handler(tok)), true
		case *MacroFunc:
			// 忽略非函数式宏
			if n := c.peekNext(); n.Lit != "(" {
				return nil, false
			}
			// 处理
		}
	}
	return nil, false
}

// 在某位置展开宏
func (c *Context) ExtractAt(pos token.Position, tks []*token.Token) []*token.Token {
	ex := make([]*token.Token, len(tks))
	col := 0
	for i, tok := range tks {
		if i == 0 {
			col = tok.Position.Column
		}
		t := &token.Token{
			Position: token.Position{
				Filename: pos.Filename,
				Line:     pos.Line,
				Column:   pos.Column + tok.Position.Column - col,
			},
			Type: tok.Type,
			Lit:  tok.Lit,
		}
		ex[i] = t
	}
	return ex
}

// peek 下一个非空 token
func (c *Context) peekNext() *token.Token {
	n := 1
	for {
		v := c.peek(n)
		if len(v) < n {
			break
		}
		if v[n-1].Type != token.WHITESPACE {
			return v[n-1]
		}
	}
	return &token.Token{
		Position: token.Position{},
		Type:     token.EOF,
		Lit:      "",
	}
}

// peek 下一个 token
func (c *Context) peekOne() *token.Token {
	if _, ok := c.in.(scanner.PeekScanner); !ok {
		c.in = scanner.NewMultiScan(scanner.NewPeekScan(c.in))
	}
	return c.in.(scanner.PeekScanner).PeekOne()
}

func (c *Context) peek(offset int) []*token.Token {
	if _, ok := c.in.(scanner.PeekScanner); !ok {
		c.in = scanner.NewMultiScan(scanner.NewPeekScan(c.in))
	}
	return c.in.(scanner.PeekScanner).Peek(offset)
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
	if c.tok == token.NEWLINE || c.tok == token.EOF {
		c.nextToken()
		return
	}
	c.err.Add(c.pos, fmt.Sprintf("expect end macro got %s", c.tok))
}

// 跳过无法到达的代码
func (c *Context) skipUtilCdt(names ...string) []*token.Token {
	cdt := 0
	tks := make([]*token.Token, 2)
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
	tks := []*token.Token{}
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
