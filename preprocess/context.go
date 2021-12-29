package preprocess

import (
	"dxkite.cn/c/errors"
	"dxkite.cn/c/scanner"
	"dxkite.cn/c/token"
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

// MacroHandler MacroVal Handler
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

// Context 解析环境
type Context struct {
	Val     map[string]MacroDecl // 宏定义
	Inc     []string             // 文件目录
	counter int                  // __COUNTER__
	once    map[string]struct{}  // #pragma once
	cdt     *ConditionStack      // 条件栈
	err     errors.ErrorList     // 错误信息
}

// NewContext 创建宏处理环境
func NewContext() *Context {
	c := &Context{}
	c.Val = map[string]MacroDecl{}
	c.cdt = NewConditionStack()
	c.once = map[string]struct{}{}
	c.err = errors.ErrorList{}
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

func (c *Context) DefineVal(name string, body []token.Token) *errors.Error {
	if err := checkValidHashHashExpr(body); err != nil {
		return err
	}
	c.Define(name, &MacroVal{
		Name: name,
		Body: body,
	})
	return nil
}

func (c *Context) DefineValStr(name, value string) *errors.Error {
	if tok, err := scanner.ScanString("<build-in>", value, nil); err != nil {
		return errors.NewStd(token.Position{}, err)
	} else {
		return c.DefineVal(name, tok)
	}
}

func (c *Context) Define(name string, val MacroDecl) {
	c.Val[name] = val
}

func (c *Context) DefineFunc(name string, params []string, ellipsis bool, body []token.Token) *errors.Error {
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

func (c *Context) Error() errors.ErrorList {
	return c.err
}

// AddError 添加错误
func (c *Context) AddError(err *errors.Error) {
	c.err.AddErr(err)
}

func (c *Context) AddErrorMsg(pos token.Position, code errors.ErrCode, args ...interface{}) {
	c.err.AddErrMsg(pos, code, args...)
}

// Top 栈顶
func (c *Context) Top() Condition {
	return c.cdt.Top()
}

// Push 压入栈
func (c *Context) Push(cdt Condition) {
	c.cdt.Push(cdt)
}

// Pop 弹出栈
func (c *Context) Pop() Condition {
	return c.cdt.Pop()
}

// SearchFile 查找文件
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
