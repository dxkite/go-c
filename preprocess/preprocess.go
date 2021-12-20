package preprocess

import (
	"dxkite.cn/c/scanner"
	"dxkite.cn/c/token"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Recorder struct {
	offset []int
	tks    []token.Token
}

func (r Recorder) Enable() bool {
	return len(r.offset) > 0
}

func (r *Recorder) Push(tok token.Token) {
	r.tks = append(r.tks, tok)
}

func (r *Recorder) Start() {
	r.offset = append(r.offset, len(r.tks))
}

func (r *Recorder) EndGet() []token.Token {
	lof := len(r.offset)
	if lof > 0 {
		idx := lof - 1
		start := r.offset[idx]
		r.offset = r.offset[:idx]
		return copyTokenSlice(r.tks[start:])
	}
	r.tks = r.tks[0:0]
	return r.tks
}

func NewRecorder() *Recorder {
	return &Recorder{
		offset: []int{},
		tks:    []token.Token{},
	}
}

type Expander struct {
	ctx       *Context
	cur       token.Token
	in        scanner.MultiScanner
	rcd       *Recorder
	err       ErrorList
	ignoreErr bool
}

// 设置输入
func NewExpander(ctx *Context, s scanner.Scanner, ignore bool) *Expander {
	e := &Expander{}
	e.in = scanner.NewMultiScan(s)
	e.ctx = ctx
	e.rcd = NewRecorder()
	e.next()
	e.ignoreErr = ignore
	e.err = ErrorList{}
	return e
}

func (e *Expander) Error() ErrorList {
	return e.err
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
		if !e.ignoreErr && e.err.Len() > 0 {
			t = &scanner.IllegalToken{
				Token: &scanner.Token{
					Pos: e.cur.Position(),
					Typ: e.cur.Type(),
					Lit: e.cur.Literal(),
				},
				Err: e.err,
			}
			return
		}
		e.next()
	}
	return t
}

func tokIsMacro(tok token.Token) bool {
	// 展开后的#不作为宏符号
	if v, ok := tok.(*Token); ok && v.Expand != nil {
		return false
	}
	return tok.Position().Column == 1 && tok.Literal() == "#"
}

// 获取下一个
func (e *Expander) next() token.Token {
	if e.rcd.Enable() && e.cur != nil {
		e.rcd.Push(e.cur)
	}
	e.cur = e.in.Scan()
	return e.cur
}

func (e *Expander) addErr(pos token.Position, msg string, args ...interface{}) {
	e.err.Add(pos, msg, args...)
}

func (e *Expander) startRecord() {
	e.rcd.Start()
}

func (e *Expander) endRecord() []token.Token {
	return e.rcd.EndGet()
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
			e.addErr(e.cur.Position(), "unexpected #else")
		}
	case "else":
		e.next()
		e.expectEndMacro()
		if e.ctx.Top() == IN_THEN {
			e.skipUtilCdt("endif")
			e.next()
			e.expectEndMacro()
		} else {
			e.addErr(e.cur.Position(), "unexpected #else")
		}
	case "endif":
		if e.ctx.Top() == IN_THEN || e.ctx.Top() == IN_ELSE {
			e.ctx.Pop()
			e.next() // endif
			e.expectEndMacro()
		} else {
			e.addErr(e.cur.Position(), "unexpected #endif")
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
		e.doLine()
	case "error":
		e.doError()
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
			d := calcDelta([]token.Token{tok}, tks)
			e.next() // 跳到当前token之后
			e.deltaLine(d)
			e.push(tks)
			e.next()
			return true
		case *MacroHandler:
			tks := e.ExpandVal(tok, val.Handler(tok), nil)
			d := calcDelta([]token.Token{tok}, tks)
			e.next() // 跳到当前token之后
			e.deltaLine(d)
			e.push(tks)
			e.next()
			return true
		case *MacroFunc:
			// 忽略非函数式宏
			if n := e.peekNext(); n.Literal() != "(" {
				return false
			}
			// 处理
			if total, tks, ok := e.ExpandFunc(tok, val); ok {
				d := calcDelta(total, tks)
				e.deltaLine(d)
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

// 后续修改偏移
func columnDelta(tks []token.Token, delta int) []token.Token {
	for i := range tks {
		tks[i] = newDeltaToken(tks[i], delta)
	}
	return tks
}

func tokenLen(tks []token.Token) int {
	size := len(tks)
	if size == 0 {
		return 0
	}
	t := tks[size-1]
	start := tks[0].Position().Column
	end := t.Position().Column + utf8.RuneCountInString(t.Literal())
	return end - start
}

// 计算展开偏移
func calcDelta(before, after []token.Token) (delta int) {
	return tokenLen(after) - tokenLen(before)
}

// 重新计算行内偏移量
func (e *Expander) deltaLine(delta int) {
	e.startRecord()
	e.skipEndMacro()
	end := e.cur
	arr := e.endRecord()
	arr = append(arr, end)
	columnDelta(arr, delta)
	e.push(arr)
}

func newTokenPos(t token.Token, pos token.Position) token.Token {
	switch v := t.(type) {
	case *Token:
		v.Pos = pos
		return v
	case *scanner.Token:
		v.Pos = pos
		return v
	default:
		return &scanner.Token{
			Pos: pos,
			Typ: t.Type(),
			Lit: t.Literal(),
		}
	}
}

func newDeltaToken(t token.Token, d int) token.Token {
	p := t.Position()
	p.Column += d
	return newTokenPos(t, p)
}

func copyToken(t token.Token) token.Token {
	if t == nil {
		return t
	}
	switch v := t.(type) {
	case *Token:
		return &Token{
			Pos: token.Position{
				Filename: t.Position().Filename,
				Line:     t.Position().Line,
				Column:   t.Position().Column,
			},
			Typ:    t.Type(),
			Lit:    t.Literal(),
			Expand: copyToken(v.Expand),
		}
	default:
		return &scanner.Token{
			Pos: token.Position{
				Filename: t.Position().Filename,
				Line:     t.Position().Line,
				Column:   t.Position().Column,
			},
			Typ: t.Type(),
			Lit: t.Literal(),
		}
	}
}

func copyTokenSlice(tks []token.Token) (cpy []token.Token) {
	for i := range tks {
		cpy = append(cpy, copyToken(tks[i]))
	}
	return
}

// 展开宏
func (e *Expander) ExpandVal(v token.Token, body []token.Token, params map[string][]token.Token) []token.Token {
	var ex []token.Token
	lt := len(body)
	if lt < 0 {
		return ex
	}
	tks := copyTokenSlice(body)
	col := tks[0].Position().Column
	pos := v.Position()
	// 展开处理
	for i := 0; i < lt; i++ {
		// 展开参数
		if tks[i].Type() == token.IDENT && params != nil {
			if v, ok := params[tks[i].Literal()]; ok {
				t := newTokenPos(tks[i], token.Position{
					Filename: pos.Filename,
					Line:     pos.Line,
					Column:   pos.Column + tks[i].Position().Column - col,
				})
				extAt := e.ExpandVal(t, v, nil)
				vv := scanner.NewArrayScan(extAt)
				exp := NewExpander(e.ctx, vv, e.ignoreErr)
				expTks, _ := scanner.ScanToken(exp)
				d := calcDelta([]token.Token{tks[i]}, expTks)
				columnDelta(tks[i+1:], d)
				ex = append(ex, expTks...)
				continue
			}
		}
		// # ## 操作
		typ := tks[i].Type()
		lit := tks[i].Literal()
		offset := tks[i].Position().Column - col
		if i+1 < lt && tks[i+1].Literal() == "##" {
			// ## 操作
			tok := tks[i]
			i += 2
			lit = tok.Literal() + tks[i].Literal()
			if !isValidToken(lit) {
				typ = token.ILLEGAL
				e.addErr(e.cur.Position(), "invalid ## operator between %s and %s", tok.Literal(), tks[i].Literal())
			}
			before := tks[i-2 : i+1]
			beforeLen := tokenLen(before)
			afterLen := len(lit)
			columnDelta(tks[i+1:], afterLen-beforeLen)
		} else if tks[i].Literal() == "#" &&
			params != nil && i+1 < lt && tks[i+1].Type() == token.IDENT {
			// # 操作
			name := tks[i+1].Literal()
			typ = token.STRING
			if v, ok := params[name]; ok {
				i++
				lit = strconv.QuoteToGraphic(relativeTokenString(v))
				before := tks[i-1 : i+1]
				beforeLen := tokenLen(before)
				afterLen := len(lit)
				columnDelta(tks[i+1:], afterLen-beforeLen)
			}
		}
		t := &Token{
			Pos: token.Position{
				Filename: pos.Filename,
				Line:     pos.Line,
				// 展开后的位置
				// 初始位置 + 相对偏移
				Column: pos.Column + offset,
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
	if tks, err := scanner.ScanToken(s); err != nil {
		return false
	} else {
		return len(tks) == 1
	}
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
	e.addErr(e.cur.Position(), fmt.Sprintf("expect token %s got %s", token.IDENT, e.cur.Type()))
	return ""
}

func (e *Expander) expectPunctuator(lit string) {
	e.punctuator(lit, true)
}

func (e *Expander) punctuator(lit string, require bool) {
	if e.cur.Type() == token.PUNCTUATOR && lit == e.cur.Literal() {
		e.nextToken()
		return
	}

	if require {
		e.addErr(e.cur.Position(), fmt.Sprintf("expect punctuator %s got %s", lit, e.cur.Literal()))
	}
}

func (e *Expander) expectEndMacro() {
	if e.isMacroEnd() {
		e.nextToken()
		return
	}
	e.addErr(e.cur.Position(), fmt.Sprintf("expect end macro got %s", e.cur.Type()))
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
			e.addErr(e.cur.Position(), fmt.Sprintf("expect %s, got EOF", strings.Join(names, ",")))
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
	var tks []token.Token
	for {
		if e.cur.Type() == token.EOF || e.cur.Type() == token.NEWLINE {
			break
		}
		if e.cur.Type() != token.WHITESPACE {
			tks = append(tks, e.cur)
		}
		e.next()
	}
	exp := NewExpander(e.ctx, scanner.NewArrayScan(tks), e.ignoreErr)
	expand, err := scanner.ScanToken(exp)
	if err != nil {
		e.addErr(e.cur.Position(), "invalid const-expr %s", inlineTokenString(tks))
	}
	return EvalConstExpr(e.ctx, expand)
}

func (e *Expander) doDefine() {
	e.nextToken()
	ident := e.expectIdent()

	if e.ctx.IsDefined(ident) {
		e.addErr(e.cur.Position(), "duplicate define of %s", ident)
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
	var tks []token.Token
	e.skipWhitespace()

	for !e.isMacroEnd() {
		tks = append(tks, e.cur)
		e.nextToken()
	}

	if pos, err := checkValidMacroExpr(tks); err != nil {
		e.addErr(pos, err.Error())
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
	var tks []token.Token
	var params []string

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
			e.addErr(e.cur.Position(), fmt.Sprintf("expect ident, got %s <%s>", e.cur.Type(), e.cur.Literal()))
			break
		}
	}

	e.expectPunctuator(")")

	for !e.isMacroEnd() {
		tks = append(tks, e.cur)
		e.nextToken()
	}

	if pos, err := checkValidMacroExpr(tks); err != nil {
		e.addErr(pos, err.Error())
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
func (e *Expander) ExpandFunc(tok token.Token, val *MacroFunc) ([]token.Token, []token.Token, bool) {
	e.startRecord()
	e.nextToken()
	params, ok := e.readParameters(val)
	total := e.endRecord()
	// 参数错误不解析
	if !ok {
		total = append(total, e.cur)
		e.push(total)
		e.next()
		return nil, nil, false
	}
	return total, e.ExpandVal(tok, val.Body, params), true
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
			e.addErr(e.cur.Position(), "expect params %d got %d", lp, i)
		}
		i++
	}
	e.expectPunctuator(")")
	if len(params) < lp {
		e.addErr(e.cur.Position(), "requires %d arguments, but only %d given", lp, len(params))
		return nil, false
	}
	return params, true
}

// 读取参数
func (e *Expander) readParameter() []token.Token {
	e.startRecord()
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
	return e.endRecord()
}

// 读取参数
func (e *Expander) readEllipsisParameter() []token.Token {
	e.startRecord()
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
	return e.endRecord()
}

func (e *Expander) doInclude() {
	// include "file"
	if e.peekNext().Type() == token.STRING {
		e.nextToken()
		p, err := strconv.Unquote(e.cur.Literal())
		if err != nil {
			e.addErr(e.cur.Position(), "invalid include string %s", e.cur.Literal())
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
		var p []token.Token
		for !e.isMacroEnd() && e.cur.Literal() != ">" {
			p = append(p, e.cur)
			e.next()
		}
		e.expectPunctuator(">")
		e.includeFile(relativeTokenString(p))
		return
	}

	// 不进行多次展开
	if _, ok := e.nextToken().(*Token); ok {
		e.addErr(e.nextToken().Position(), "invalid #include")
		return
	}

	// 展开include后重新处理
	e.startRecord()
	e.skipEndMacro()
	p := e.endRecord()
	exp := NewExpander(e.ctx, scanner.NewArrayScan(p), e.ignoreErr)
	expand, err := scanner.ScanToken(exp)
	if err != nil {
		e.addErr(e.cur.Position(), "invalid include string %s", inlineTokenString(p))
	}
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
		sc, err := scanner.NewFileScan(fn)
		if err != nil {
			e.ctx.AddError(e.cur.Position(), "include %s: %s", fn, err.Error())
			return
		}
		e.push([]token.Token{e.cur})
		e.in.Push(sc)
		e.next()
	} else {
		e.ctx.AddError(e.cur.Position(), "file not found %s", s)
	}
}

func (e *Expander) doPragma() {
	e.startRecord()
	e.skipEndMacro()
	p := e.endRecord()
	for _, v := range p {
		// 支持 pragma once 指令
		if v.Literal() == "once" {
			pp, _ := filepath.Abs(v.Position().Filename)
			e.ctx.pragmaOnce(pp)
		}
	}
	e.expectEndMacro()
}

func (e *Expander) doLine() {
	e.startRecord()
	file := e.cur.Position().Filename
	line := e.cur.Position().Line
	e.skipEndMacro()
	v := e.endRecord()
	enable := false
	for _, item := range v {
		if item.Type() == token.INT {
			enable = true
			line, _ = strconv.Atoi(item.Literal())
		}
		if item.Type() == token.STRING {
			enable = true
			file, _ = strconv.Unquote(item.Literal())
		}
	}
	if enable {
		e.in = mockLine(e.in, e.cur.Position(), file, line)
	}
	e.expectEndMacro()
}

func (e *Expander) doError() {
	pos := e.cur.Position()
	e.next() // error
	e.skipWhitespace()
	e.startRecord()
	e.skipEndMacro()
	msg := e.endRecord()
	e.expectEndMacro()
	e.addErr(pos, inlineTokenString(msg))
}

type lineDirective struct {
	scanner.MultiScanner
	scope string
	delta int
	file  string
}

func (ld *lineDirective) Scan() token.Token {
	t := ld.MultiScanner.Scan()
	if t.Position().Filename == ld.scope {
		tt := &Token{
			Pos:    t.Position(),
			Typ:    t.Type(),
			Lit:    t.Literal(),
			Expand: nil,
		}
		tt.Pos.Filename = ld.file
		tt.Pos.Line = tt.Pos.Line + ld.delta
		return tt
	}
	return t
}

func mockLine(s scanner.MultiScanner, pos token.Position, file string, line int) scanner.MultiScanner {
	ld := &lineDirective{}
	ld.delta = line - pos.Line
	ld.file = file
	ld.scope = pos.Filename
	ld.MultiScanner = s
	return ld
}

type preprocess struct {
	scanner.Scanner
}

// 预处理扫描器
func NewScanner(ctx *Context, s scanner.Scanner, ignore bool) scanner.Scanner {
	exp := NewExpander(ctx, s, ignore)
	pre := &preprocess{scanner.NewTokenScan(exp)}
	return pre
}
