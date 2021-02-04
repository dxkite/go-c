package preprocess

import (
	"dxkite.cn/go-c11/scanner"
	"dxkite.cn/go-c11/token"
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
	lit     string
	pos     token.Position
	tok     token.Type
	cur     *token.Token
	in, out scanner.MultiScanner
	rcd     bool
	tks     []*token.Token
	err     scanner.ErrorList
}

// 设置输入
func New(s scanner.Scanner) *Context {
	c := &Context{}
	c.in = scanner.NewMultiScan(s)
	c.out = scanner.NewMultiScan()
	c.err.Reset()
	return c
}

// 测试 once
func (c *Context) Scan() (t *token.Token) {
	for t == nil {
		t = c.out.Scan()
		// 无输出
		if t == nil {
			c.next()
			if c.tok == token.EOF {
				break
			}
			if tokIsMacro(c.cur) {
				// do Macro
				continue
			}
			if c.tok == token.IDENT && c.IsDefined(c.cur) {
				// do extract
				continue
			}
			t = c.cur
		}
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

func (c *Context) IsDefined(tok *token.Token) bool {
	_, ok := c.Val[tok.Lit]
	return ok
}

// 展开宏
func (c *Context) Extract(tok *token.Token) []*token.Token {
	if v, ok := c.Val[tok.Lit]; ok {
		switch val := v.(type) {
		case *MacroVal:
			return c.ExtractAt(tok.Position, val.Body)
		case *MacroFunc:
		case *MacroHandler:

		}
	}
	return []*token.Token{tok}
}

// 在某位置展开宏
func (c *Context) ExtractAt(pos token.Position, tks []*token.Token) []*token.Token {
	ex := []*token.Token{}
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
		ex = append(ex, t)
	}
	return ex
}
