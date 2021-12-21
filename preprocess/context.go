package preprocess

import (
	"dxkite.cn/c/scanner"
	"dxkite.cn/c/token"
	"os"
	"path"
	"path/filepath"
	"strconv"
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

type Token struct {
	Pos    token.Position
	Typ    token.Type
	Lit    string
	Expand token.Token // 父级展开
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

func (t *Token) ExpandFrom(tok token.Token) bool {
	var exp token.Token
	// 获取展开历史
	if t, ok := tok.(*Token); ok && t.Expand != nil {
		exp = t.Expand
	} else {
		return false
	}
	// 不递归展开
	for {
		if exp.Literal() == tok.Literal() {
			return true
		}
		if t, ok := exp.(*Token); ok && t.Expand != nil {
			exp = t.Expand
		} else {
			return false
		}
	}
}

// 解析环境
type Context struct {
	Val     map[string]MacroDecl // 宏定义
	Inc     []string             // 文件目录
	counter int                  // __COUNTER__
	once    map[string]struct{}  // #pragma once
	cdt     *ConditionStack      // 条件栈
	err     ErrorList
}

func NewContext() *Context {
	c := &Context{}
	c.Val = map[string]MacroDecl{}
	c.cdt = NewConditionStack()
	c.once = map[string]struct{}{}
	return c
}

// 测试 once
func (c *Context) onceContain(p string) bool {
	pp, _ := filepath.Abs(p)
	_, ok := c.once[pp]
	return ok
}

// 写入 pragma once
func (c *Context) pragmaOnce(p string) {
	c.once[p] = struct{}{}
}

func (c *Context) DefineVal(name string, body []token.Token) error {
	if err := checkValidHashHashExpr(body); err != nil {
		return err
	}
	c.Define(name, &MacroVal{
		Name: name,
		Body: body,
	})
	return nil
}

func (c *Context) DefineValStr(name, value string) error {
	if tok, err := scanner.ScanString("<build-in>", value); err != nil {
		return err
	} else {
		return c.DefineVal(name, tok)
	}
}

func (c *Context) Define(name string, val MacroDecl) {
	c.Val[name] = val
}

// 检查是否可用
func checkValidHashHashExpr(tks []token.Token) error {
	if len(tks) > 0 {
		beg := 0
		end := len(tks) - 1
		err := "'##' cannot appear at either end of a macro expansion"
		if tks[beg].Literal() == "##" {
			return &Error{tks[beg].Position(), err}
		}
		if tks[end].Literal() == "##" {
			return &Error{tks[end].Position(), err}
		}
	}
	return nil
}

// 检查是否可用
func checkValidHashExpr(params []string, tks []token.Token) error {
	err := "'#' must follow a macro parameter"
	m := map[string]bool{}
	for _, v := range params {
		m[v] = true
	}
	n := len(tks)
	for i := 0; i < n; i++ {
		if tks[i].Literal() == "#" {
			if !(i+1 < n && tks[i+1].Type() == token.IDENT && m[tks[i+1].Literal()]) {
				return &Error{tks[i].Position(), err}
			}
		}
	}
	return nil
}

func (c *Context) DefineFunc(name string, params []string, ellipsis bool, body []token.Token) error {
	if err := checkValidHashExpr(params, body); err != nil {
		return err
	}
	if err := checkValidHashHashExpr(body); err != nil {
		return err
	}
	c.Define(name, &MacroFunc{
		Name:     name,
		Params:   params,
		Ellipsis: ellipsis,
		Body:     body,
	})
	return nil
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
	_ = c.DefineValStr("__DATE__", strconv.QuoteToGraphic(time.Now().Format("Jan 02 2006")))
	_ = c.DefineValStr("__TIME__", strconv.QuoteToGraphic(time.Now().Format("15:04:05")))
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

func (c *Context) Error() ErrorList {
	return c.err
}

// 添加错误
func (c *Context) AddError(pos token.Position, msg string, args ...interface{}) {
	c.err.Add(pos, msg, args...)
}

// 栈顶
func (c *Context) Top() Condition {
	return c.cdt.Top()
}

// 压入栈
func (c *Context) Push(cdt Condition) {
	c.cdt.Push(cdt)
}

// 弹出栈
func (c *Context) Pop() Condition {
	return c.cdt.Pop()
}

func exists(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// 查找文件
func (c *Context) SearchFile(name string, cur string) (string, bool) {
	if p := path.Join(cur, name); exists(p) {
		return p, true
	}
	for _, rp := range c.Inc {
		if p := path.Join(rp, name); exists(p) {
			return p, true
		}
	}
	return "", false
}
