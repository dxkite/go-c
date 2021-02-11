package preprocess

import (
	"dxkite.cn/go-c11"
	"dxkite.cn/go-c11/scanner"
	"dxkite.cn/go-c11/token"
	"errors"
	"fmt"
	"path"
	"path/filepath"
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
	err     go_c11.ErrorList
}

func NewContext() *Context {
	c := &Context{}
	c.err.Reset()
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
	c.Val[name] = val
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

func (c *Context) Error() *go_c11.ErrorList {
	return &c.err
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

// 查找文件
func (c *Context) SearchFile(name string, cur string) (string, bool) {
	if p := path.Join(cur, name); go_c11.Exists(p) {
		return p, true
	}
	for _, rp := range c.Inc {
		if p := path.Join(rp, name); go_c11.Exists(p) {
			return p, true
		}
	}
	return "", false
}

type Expander struct {
	ctx *Context
	cur token.Token
	in  scanner.MultiScanner
	rcd bool
	tks []token.Token
	err go_c11.ErrorList
}

// 设置输入
func NewExpander(ctx *Context, s scanner.Scanner) *Expander {
	e := &Expander{}
	e.in = scanner.NewMultiScan(s)
	e.ctx = ctx
	e.next()
	return e
}

func (e *Expander) Scan() (t token.Token) {
	for t == nil {
		if e.cur.Type() == token.EOF {
			return e.cur
		}
		if tokIsMacro(e.cur) {
			e.doMacro()
			continue
		}
		// 宏展开
		if e.Expand(e.cur) {
			continue
		}
		// 普通token
		t = e.cur
		e.next()
	}
	return t
}

func (c *Expander) Error() *go_c11.ErrorList {
	return &c.err
}

func tokIsMacro(tok token.Token) bool {
	// 展开后的#不作为宏符号
	if _, ok := tok.(*Token); ok {
		return false
	}
	return tok.Position().Column == 1 && tok.Literal() == "#"
}

// 获取下一个
func (e *Expander) next() token.Token {
	if e.rcd && e.cur != nil {
		e.tks = append(e.tks, e.cur)
	}
	e.cur = e.in.Scan()
	return e.cur
}

func (e *Expander) record() {
	e.rcd = true
	e.tks = e.tks[0:0]
}

func (e *Expander) arr() []token.Token {
	pp := []token.Token{}
	pp = append(pp, e.tks...)
	return pp
}

// 获取下一个非空token
func (e *Expander) nextToken() token.Token {
	for {
		e.next()
		if e.cur.Type() != token.WHITESPACE {
			break
		}
	}
	return e.cur
}

// 获取下一个非空token
func (e *Expander) skipWhitespace() token.Token {
	for {
		if e.cur.Type() != token.WHITESPACE {
			break
		}
		e.next()
	}
	return e.cur
}

func (e *Expander) doMacro() {
	e.nextToken()
	switch e.cur.Literal() {
	case "if":
		e.next()
		cdt := e.evalConstExpr()
		e.expectEndMacro()
		if cdt {
			e.ctx.Push(IN_THEN)
		} else {
			// 跳到下一个分支
			e.ctx.Push(IN_ELSE)
			e.skipUtilElse()
		}
	case "ifdef":
		e.doIfDefine(true)
	case "ifndef":
		e.doIfDefine(false)
	case "elif":
		e.next()
		if e.ctx.Top() == IN_ELSE {
			cdt := e.evalConstExpr()
			e.expectEndMacro()
			if cdt {
				e.ctx.Pop()
				e.ctx.Push(IN_THEN)
			} else {
				// 跳到下一个分支
				e.skipUtilElse()
			}
		} else if e.ctx.Top() == IN_THEN {
			// 直接跳到结尾
			e.skipUtilCdt("endif")
			e.next() // endif
			e.expectEndMacro()
		} else {
			e.err.Add(e.cur.Position(), "unexpected #else")
		}
	case "else":
		e.next()
		e.expectEndMacro()
		if e.ctx.Top() == IN_THEN {
			e.skipUtilCdt("endif")
			e.next()
			e.expectEndMacro()
		} else {
			e.err.Add(e.cur.Position(), "unexpected #else")
		}
	case "endif":
		if e.ctx.Top() == IN_THEN || e.ctx.Top() == IN_ELSE {
			e.ctx.Pop()
			e.next() // endif
			e.expectEndMacro()
		} else {
			e.err.Add(e.cur.Position(), "unexpected #endif")
		}
	case "define":
		e.doDefine()
	case "undef":
		e.doUndef()
	case "include":
		e.doInclude()
	case "pragma":
		e.doPragma()
	case "line":

	case "error":
	default:
		e.skipEndMacro()
	}
}

// 重新压入token
func (e *Expander) push(tok []token.Token) {
	e.in.Push(scanner.NewArrayScan(tok))
}

// 处理宏展开
func (e *Expander) Expand(tok token.Token) bool {
	// 非标识符不展开
	if tok.Type() != token.IDENT {
		return false
	}

	// 不允许递归展开
	if t, ok := tok.(*Token); ok && t.ExpandFrom(tok) {
		return false
	}

	name := tok.Literal()
	if v, ok := e.ctx.Val[name]; ok {

		switch val := v.(type) {
		case *MacroVal:
			tks := e.ExpandVal(tok, val.Body, nil)
			e.push(tks)
			e.next()
			return true
		case *MacroHandler:
			tks := e.ExpandVal(tok, val.Handler(tok), nil)
			e.push(tks)
			e.next()
			return true
		case *MacroFunc:
			// 忽略非函数式宏
			if n := e.peekNext(); n.Literal() != "(" {
				return false
			}
			// 处理
			if tks, ok := e.ExpandFunc(tok, val); ok {
				tks = append(tks, e.cur)
				e.push(tks)
				e.next()
				return true
			} else {
				return false
			}
		}
	}
	return false
}

// 展开宏
func (e *Expander) ExpandVal(v token.Token, tks []token.Token, params map[string][]token.Token) []token.Token {
	ex := make([]token.Token, 0)
	col := 0
	pos := v.Position()
	ps := scanner.NewPeekScan(scanner.NewTokenScan(scanner.NewArrayScan(tks)))
	f := true

	for tok := ps.Scan(); tok.Type() != token.EOF; tok = ps.Scan() {
		if f {
			col = tok.Position().Column
			f = false
		}

		typ := tok.Type()
		lit := tok.Literal()

		// ## 操作
		if ps.PeekOne().Literal() == "##" {
			ps.Scan() //##
			nxt := ps.Scan()
			lit = tok.Literal() + nxt.Literal()
			if !isValidToken(lit) {
				typ = token.ILLEGAL
				e.err.Add(e.cur.Position(), "invalid ## operator between %s and %s", lit, nxt.Literal())
			}
		}

		if tok.Literal() == "#" &&
			params != nil && ps.PeekOne().Type() == token.IDENT {
			name := ps.PeekOne().Literal()
			typ = token.STRING
			if v, ok := params[name]; ok {
				ps.Scan()
				lit = strconv.QuoteToGraphic(token.RelativeString(v))
			}
		}

		if tok.Type() == token.IDENT && params != nil {
			if v, ok := params[tok.Literal()]; ok {
				vv := scanner.NewArrayScan(v)
				exp := NewExpander(e.ctx, vv)
				tks := scanner.ScanToken(exp)
				ex = append(ex, tks...)
				if exp.err.Len() > 0 {
					e.err.Merge(exp.err)
				}
				continue
			}
		}

		t := &Token{
			Pos: token.Position{
				Filename: pos.Filename,
				Line:     pos.Line,
				Column:   pos.Column + tok.Position().Column - col,
			},
			Typ:    typ,
			Lit:    lit,
			Expand: v,
		}
		ex = append(ex, t)
	}
	return ex
}

func isValidToken(lit string) bool {
	s := scanner.NewStringScan("<runtime>", lit)
	tok := scanner.ScanToken(s)
	if len(tok) > 1 {
		return false
	}
	if s.Error().Len() > 0 {
		return false
	}
	return true
}

// peek 下一个非空 token
func (e *Expander) peekNext() token.Token {
	n := 1
	for {
		v := e.peek(n)
		if len(v) < n {
			break
		}
		if v[n-1].Type() != token.WHITESPACE {
			return v[n-1]
		}
		n++
	}
	return &Token{
		Pos: token.Position{},
		Typ: token.EOF,
		Lit: "",
	}
}

func (e *Expander) peek(offset int) []token.Token {
	if ps, ok := e.in.(scanner.PeekScanner); ok {
		return ps.Peek(offset)
	}
	p := scanner.NewPeekScan(e.in)
	tok := p.Peek(offset)
	e.in = scanner.NewMultiScan(p)
	return tok
}

func (e *Expander) expectIdent() string {
	if e.cur.Type() == token.IDENT {
		lit := e.cur.Literal()
		e.next()
		return lit
	}
	e.err.Add(e.cur.Position(), fmt.Sprintf("expect token %s got %s", token.IDENT, e.cur.Type()))
	return ""
}

func (e *Expander) expectPunctuator(lit string) {
	e.punctuator(lit, true)
}

func (e *Expander) punctuator(lit string, require bool) {
	if e.cur.Type() == token.PUNCTUATOR && lit == e.cur.Literal() {
		e.nextToken()
	}
	if require {
		e.err.Add(e.cur.Position(), fmt.Sprintf("expect punctuator %s got %s", lit, e.cur.Literal()))
	}
}

func (e *Expander) expectEndMacro() {
	if e.isMacroEnd() {
		e.nextToken()
		return
	}
	e.err.Add(e.cur.Position(), fmt.Sprintf("expect end macro got %s", e.cur.Type()))
}

// 宏结尾
func (e *Expander) isMacroEnd() bool {
	return e.cur.Type() == token.NEWLINE || e.cur.Type() == token.EOF
}

// 跳过无法到达的代码
func (e *Expander) skipUtilCdt(names ...string) []token.Token {
	cdt := 0
	tks := make([]token.Token, 2)
	for {
		e.next()
		if e.cur.Type() == token.EOF {
			e.err.Add(e.cur.Position(), fmt.Sprintf("expect %s, got EOF", strings.Join(names, ",")))
			break
		}
		if tokIsMacro(e.cur) {
			tks[0] = e.cur
			e.nextToken()
			switch e.cur.Literal() {
			case "if", "ifndef", "ifdef":
				cdt++
			default:
				if cdt == 0 {
					for _, name := range names {
						if name == e.cur.Literal() {
							tks[1] = e.cur
							return tks
						}
					}
					if e.cur.Literal() == "endif" {
						return tks[0:0]
					}
				}
				if e.cur.Literal() == "endif" {
					cdt--
				}
			}
		}
	}
	return tks[0:0]
}

// #ifdef #ifndef
func (e *Expander) doIfDefine(want bool) {
	e.nextToken()
	ident := e.expectIdent()
	cdt := e.ctx.IsDefined(ident)
	if cdt == want {
		e.ctx.Push(IN_THEN)
	} else {
		e.ctx.Push(IN_ELSE)
		e.skipUtilElse()
	}
	e.expectEndMacro()
}

// skip to #else/#elif
func (e *Expander) skipUtilElse() {
	m := e.skipUtilCdt("elif", "else")
	if e.cur.Literal() == "elif" {
		e.next()  // elif
		e.push(m) // push back
	} else {
		e.next() // else
	}
}

func (e *Expander) evalConstExpr() bool {
	tks := []token.Token{}
	for {
		if e.cur.Type() == token.EOF || e.cur.Type() == token.NEWLINE {
			break
		}
		if e.cur.Type() != token.WHITESPACE {
			tks = append(tks, e.cur)
		}
		e.next()
	}
	// TODO parse expr
	return false
}

func (e *Expander) doDefine() {
	e.nextToken()
	ident := e.expectIdent()

	if e.ctx.IsDefined(ident) {
		e.err.Add(e.cur.Position(), "duplicate define of %s", ident)
		e.skipEndMacro()
		return
	}

	if e.cur.Literal() == "(" {
		e.doDefineFunc(ident)
	} else {
		e.doDefineVal(ident)
	}
	e.expectEndMacro()
}

func (e *Expander) skipEndMacro() {
	for !e.isMacroEnd() {
		e.nextToken()
	}
}

func (e *Expander) doUndef() {
	e.nextToken()
	ident := e.expectIdent()
	delete(e.ctx.Val, ident)
	e.skipEndMacro()
	e.expectEndMacro()
}

func (e *Expander) doDefineVal(ident string) {
	tks := make([]token.Token, 0)
	e.skipWhitespace()

	for !e.isMacroEnd() {
		tks = append(tks, e.cur)
		e.nextToken()
	}

	if pos, err := checkValidMacroExpr(tks); err != nil {
		e.err.Add(pos, err.Error())
		return
	}

	e.ctx.Define(ident, &MacroVal{
		Name: ident,
		Body: tks,
	})
}

// 检查是否可用
func checkValidMacroExpr(tks []token.Token) (token.Position, error) {
	if len(tks) > 0 {
		beg := 0
		end := len(tks) - 1
		err := errors.New("'##' cannot appear at either end of a macro expansion")
		if tks[beg].Literal() == "##" {
			return tks[beg].Position(), err
		}
		if tks[end].Literal() == "##" {
			return tks[end].Position(), err
		}
	}
	return token.Position{}, nil
}

func (e *Expander) doDefineFunc(ident string) {
	tks := make([]token.Token, 0)
	params := make([]string, 0)
	e.expectPunctuator("(")

	elp := false
	for !e.isMacroEnd() && e.cur.Literal() != ")" {
		if e.cur.Literal() == "..." {
			elp = true
			e.nextToken()
			break
		} else if e.cur.Type() == token.IDENT {
			params = append(params, e.cur.Literal())
			e.nextToken()
			e.punctuator(",", false)
		} else {
			e.err.Add(e.cur.Position(), fmt.Sprintf("expect ident, got %s <%s>", e.cur.Type(), e.cur.Literal()))
			break
		}
	}

	e.expectPunctuator(")")

	for !e.isMacroEnd() {
		tks = append(tks, e.cur)
		e.nextToken()
	}

	if pos, err := checkValidMacroExpr(tks); err != nil {
		e.err.Add(pos, err.Error())
		return
	}

	e.ctx.Define(ident, &MacroFunc{
		Name:     ident,
		Params:   params,
		Ellipsis: elp,
		Body:     tks,
	})
}

// 展开函数
func (e *Expander) ExpandFunc(tok token.Token, val *MacroFunc) ([]token.Token, bool) {
	e.record()
	e.nextToken()
	params, ok := e.readParameters(val)
	scan := e.arr()
	// 参数错误不解析
	if !ok {
		scan = append(scan, e.cur)
		e.push(scan)
		e.next()
		return nil, false
	}
	return e.ExpandVal(tok, val.Body, params), true
}

func (e *Expander) readParameters(val *MacroFunc) (map[string][]token.Token, bool) {
	e.expectPunctuator("(")
	params := map[string][]token.Token{}
	i := 0
	lp := len(val.Params)
	for !e.isMacroEnd() && e.cur.Literal() != ")" {
		if len(params) < lp {
			p := e.readParameter()
			params[val.Params[i]] = p
			e.punctuator(",", i+1 < lp)
		} else if val.Ellipsis {
			params["__VA_ARGS__"] = e.readEllipsisParameter()
		} else {
			e.err.Add(e.cur.Position(), "expect params %d got %d", lp, i)
		}
		i++
	}
	e.expectPunctuator(")")
	if len(params) < lp {
		e.err.Add(e.cur.Position(), "requires %d arguments, but only %d given", lp, len(params))
		return nil, false
	}
	return params, true
}

// 读取参数
func (e *Expander) readParameter() []token.Token {
	e.record()
	paren := 0
	for !e.isMacroEnd() && e.cur.Literal() != "," && e.cur.Literal() != ")" {
		if e.cur.Literal() == "(" {
			paren++
		}
		if e.cur.Literal() == ")" && paren != 0 {
			paren--
			e.nextToken()
		}
		e.nextToken()
	}
	return e.arr()
}

// 读取参数
func (e *Expander) readEllipsisParameter() []token.Token {
	e.record()
	paren := 0
	for !e.isMacroEnd() && e.cur.Literal() != ")" {
		if e.cur.Literal() == "(" {
			paren++
		}
		if e.peekNext().Literal() == ")" && paren > 0 {
			paren--
			e.nextToken()
		}
		e.nextToken()
	}
	return e.arr()
}

func (e *Expander) doInclude() {
	// include "file"
	if e.peekNext().Type() == token.STRING {
		e.nextToken()
		p, err := strconv.Unquote(e.cur.Literal())
		if err != nil {
			e.err.Add(e.cur.Position(), "invalid include string %s", e.cur.Literal())
		}
		e.nextToken()
		e.skipEndMacro() // 跳到换行
		e.includeFile(p)
		return
	}

	// include <file>
	if e.peekNext().Literal() == "<" {
		e.nextToken()
		e.expectPunctuator("<")
		p := []token.Token{}
		for !e.isMacroEnd() && e.cur.Literal() != ">" {
			p = append(p, e.cur)
			e.next()
		}
		e.expectPunctuator(">")
		e.includeFile(token.RelativeString(p))
		return
	}

	// 不进行多次展开
	if _, ok := e.nextToken().(*Token); ok {
		e.err.Add(e.nextToken().Position(), "invalid #include")
		return
	}

	// 展开include后重新处理
	e.record()
	e.skipEndMacro()
	p := e.arr()
	exp := NewExpander(e.ctx, scanner.NewArrayScan(p))
	expand := scanner.ScanToken(exp)
	e.push(expand)
	e.doInclude()
}

func (e *Expander) searchFile(s string, tok token.Token) (string, bool) {
	return e.ctx.SearchFile(s, filepath.Dir(tok.Position().Filename))
}

func (e *Expander) includeFile(s string) {
	if fn, ok := e.searchFile(s, e.cur); ok {

		if e.ctx.onceContain(fn) {
			e.skipEndMacro()
			e.expectEndMacro()
			return
		}

		sc := scanner.NewFileScan(fn)
		e.push([]token.Token{e.cur})
		e.in.Push(sc)
		e.next()
	} else {
		e.ctx.AddError(e.cur.Position(), "file not found %s", s)
	}
}

func (e *Expander) doPragma() {
	e.record()
	e.skipEndMacro()
	p := e.arr()
	for _, v := range p {
		// 支持 pragma once 指令
		if v.Literal() == "once" {
			pp, _ := filepath.Abs(v.Position().Filename)
			e.ctx.pragmaOnce(pp)
		}
	}
}

type preprocess struct {
	scanner.Scanner
}

// 预处理扫描器
func NewPreprocess(ctx *Context, s scanner.Scanner) scanner.Scanner {
	exp := NewExpander(ctx, s)
	pre := &preprocess{scanner.NewTokenScan(exp)}
	return pre
}
