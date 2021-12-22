package preprocess

import (
	"dxkite.cn/c/scanner"
	"dxkite.cn/c/token"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

type processor struct {
	ctx       *Context
	cur       token.Token
	r         scanner.Scanner
	err       ErrorList
	ignoreErr bool
}

// 设置输入
func newProcessor(ctx *Context, s scanner.Scanner, ignoreErr bool) *processor {
	e := &processor{}
	e.r = s
	e.ctx = ctx
	e.next()
	e.ignoreErr = ignoreErr
	e.err = ErrorList{}
	return e
}

func (p *processor) Error() ErrorList {
	return p.err
}

func (p *processor) Scan() (t token.Token) {
	for t == nil {
		if p.cur.Type() == token.EOF {
			return p.cur
		}
		// 宏定义
		if isMacroTok(p.cur) {
			p.doMacro()
			continue
		}
		// 宏展开
		if p.expand(p.cur) {
			continue
		}
		// 普通token
		t = p.cur
		if !p.ignoreErr && p.err.Len() > 0 {
			t = &scanner.IllegalToken{
				Token: &scanner.Token{
					Pos: p.cur.Position(),
					Typ: p.cur.Type(),
					Lit: p.cur.Literal(),
				},
				Err: p.err,
			}
			return
		}
		p.next()
	}
	return t
}

// 处理宏展开
func (p *processor) expand(tok token.Token) bool {
	// 非标识符不展开
	if tok.Type() != token.IDENT {
		return false
	}

	// 不允许递归展开
	if t, ok := tok.(*Token); ok && t.ExpandFrom(tok) {
		return false
	}

	name := tok.Literal()
	if v, ok := p.ctx.Val[name]; ok {
		switch val := v.(type) {
		case *MacroVal:
			tks := p.expandVal(tok, val.Body)
			d := calcDelta([]token.Token{tok}, tks)
			p.deltaLine(d)
			p.push(tks)
			p.next()
			return true
		case *MacroHandler:
			tks := p.expandVal(tok, val.Handler(tok))
			d := calcDelta([]token.Token{tok}, tks)
			p.deltaLine(d)
			p.push(tks)
			p.next()
			return true
		case *MacroFunc:
			// 忽略非函数式宏
			if n := p.peekNext(); n.Literal() != "(" {
				return false
			}
			if t, tks, ok := p.expandFunc(tok, val); ok {
				d := calcDelta(t, tks)
				p.deltaLine(d)
				p.push(tks)
				p.next()
				return true
			}
		}
	}
	return false
}

// 处理 ##
func (p *processor) expandVal(v token.Token, body []token.Token) []token.Token {
	return p.expandMacroBody(v, body, nil)
}

// 展开函数
func (p *processor) expandFunc(tok token.Token, val *MacroFunc) ([]token.Token, []token.Token, bool) {
	c := p.startCache()
	p.nextToken() // ident
	params, ok := p.readParameters(val)
	if !ok {
		c.Restore()
		return nil, nil, false
	}
	p.push([]token.Token{p.cur})
	total := c.GetClear()
	fmt.Println(tok.Literal(), "param", inlineTokenString(total))
	body := p.expandMacroBody(tok, val.Body, params)
	fmt.Println(tok.Literal(), inlineTokenString(val.Body), "=>", inlineTokenString(body))
	return total, body, true
}

func (p *processor) expandMacroBodyToken(tok, ident token.Token, params map[string][]token.Token, afterHashHash, followHashHash bool) (exp []token.Token) {
	// 调整展开位置
	ident = newTokenPos(ident, token.Position{
		Filename: tok.Position().Filename,
		Line:     tok.Position().Line,
		Column:   ident.Position().Column,
	})

	if ident.Type() != token.IDENT {
		return []token.Token{ident}
	}

	if cur, ok := params[ident.Literal()]; ok {
		fmt.Println("cur", ident.Literal(), "=>", inlineTokenString(cur))
		if afterHashHash {
			exp = append(exp, cur[0])
			cur = cur[1:]
		}
		nc := len(cur)
		var tail token.Token
		if nc > 0 && followHashHash {
			tail = cur[nc-1]
			cur = cur[:nc-1]
		}
		delta := 0
		// 展开其他部分
		if len(cur) > 0 {
			tokens, _ := scanner.ScanToken(newProcessor(p.ctx, newExpandMock(tok, scanner.NewArrayScan(cur)), p.ignoreErr))
			fmt.Println("expand", inlineTokenString(cur), "=>", inlineTokenString(tokens))
			exp = append(exp, tokens...)
			delta = calcDelta(cur, tokens)
		}
		if tail != nil {
			t := &Token{
				Pos: token.Position{
					Filename: tail.Position().Filename,
					Line:     tail.Position().Line,
					Column:   tail.Position().Column + delta,
				},
				Typ:    tail.Type(),
				Lit:    tail.Literal(),
				Expand: tok,
			}
			exp = append(exp, t)
		}
		expandTokAt(ident, exp)
		return
	}

	if afterHashHash || followHashHash {
		return []token.Token{ident}
	}

	// 普通展开
	exp, _ = scanner.ScanToken(newProcessor(p.ctx, newExpandMock(tok, scanner.NewArrayScan([]token.Token{ident})), p.ignoreErr))
	expandTokAt(ident, exp)
	return
}

func expandTokAt(tok token.Token, tks []token.Token) {
	pos := tok.Position()
	base := tks[0].Position()
	for i := range tks {
		offset := tks[i].Position().Column - base.Column
		tks[i] = newTokenPos(tks[i], token.Position{
			Filename: pos.Filename,
			Line:     pos.Line,
			Column:   offset + pos.Column,
		})
	}
}

// 展开参数
func (p *processor) expandMacroBody(tok token.Token, body []token.Token, params map[string][]token.Token) []token.Token {
	var exp []token.Token
	n := len(body)
	if n == 0 {
		return exp
	}

	tks := copyTokenSlice(body)
	col := tks[0].Position().Column
	pos := tok.Position()
	// 展开处理
	for i := 0; i < n; i++ {
		typ := tks[i].Type()
		lit := tks[i].Literal()
		offset := tks[i].Position().Column - col
		// 如果是标识符
		if tks[i].Type() == token.IDENT {
			afterHashHash := i > 0 && tks[i-1].Literal() == "##"
			followHashHash := i+1 < n && tks[i+1].Literal() == "##"
			tokens := p.expandMacroBodyToken(tok, tks[i], params, afterHashHash, followHashHash)
			exp = append(exp, tokens...)
			d := calcDelta([]token.Token{tks[i]}, tokens)
			if d != 0 {
				columnDelta(tks[i+1:], d)
			}
			continue
		}

		if tks[i].Literal() == "##" && i+1 < n && len(exp) > 0 {
			followHashHash := i+2 < n && tks[i+2].Literal() == "##"
			tokens := p.expandMacroBodyToken(tok, tks[i+1], params, true, followHashHash)
			tail := len(exp) - 1
			beforeTok := exp[tail]
			afterTok := tokens[0]
			lit = beforeTok.Literal() + afterTok.Literal()
			if !isValidToken(lit) {
				typ = token.ILLEGAL
				p.addErr(p.cur.Position(), "invalid ## operator between %s and %s", beforeTok.Literal(), afterTok.Literal())
			}

			exp[tail] = &Token{
				Pos:    exp[tail].Position(),
				Typ:    typ,
				Lit:    lit,
				Expand: tok,
			}
			exp = append(exp, tokens[1:]...)

			beforeLen := tokenLen([]token.Token{beforeTok, afterTok})
			afterLen := len(lit) + tokenLen(tokens[1:])
			columnDelta(tks[i+1:], afterLen-beforeLen)
			i++
			continue
		}

		// 处理 #
		if params != nil && tks[i].Literal() == "#" && i+1 < n && tks[i+1].Type() == token.IDENT {
			// # 操作
			name := tks[i+1].Literal()
			beforeLen := tokenLen([]token.Token{tks[i], tks[i+1]})
			typ = token.STRING
			v := params[name]
			lit = strconv.QuoteToGraphic(relativeTokenString(v))
			afterLen := len(lit)
			columnDelta(tks[i+1:], afterLen-beforeLen)
			i++
		}

		t := &Token{
			Pos: token.Position{
				Filename: pos.Filename,
				Line:     pos.Line,
				Column:   pos.Column + offset,
			},
			Typ:    typ,
			Lit:    lit,
			Expand: tok,
		}
		exp = append(exp, t)
	}
	return exp
}

func (p *processor) readParameters(val *MacroFunc) (map[string][]token.Token, bool) {
	p.expectPunctuator("(")
	params := map[string][]token.Token{}
	i := 0
	n := len(val.Params)
	for !p.isMacroEnd() && p.cur.Literal() != ")" {
		if len(params) < n {
			pp := p.readParameter()
			params[val.Params[i]] = pp
			p.punctuator(",", i+1 < n)
		} else if val.Ellipsis {
			params["__VA_ARGS__"] = p.readEllipsisParameter()
		} else {
			p.addErr(p.cur.Position(), "expect params %d got %d", n, i)
		}
		i++
	}
	p.expectPunctuator(")")
	if len(params) < n {
		p.addErr(p.cur.Position(), "requires %d arguments, but only %d given", n, len(params))
		return nil, false
	}
	return params, true
}

// 读取参数
func (p *processor) readParameter() (tks []token.Token) {
	paren := 0
	for !p.isMacroEnd() && p.cur.Literal() != "," && p.cur.Literal() != ")" {
		tks = append(tks, p.cur)
		if p.cur.Literal() == "(" {
			paren++
		}
		p.next()
		if p.cur.Literal() == ")" && paren != 0 {
			tks = append(tks, p.cur)
			paren--
			p.next()
		}
	}
	return
}

// 读取参数
func (p *processor) readEllipsisParameter() (tks []token.Token) {
	paren := 0
	for !p.isMacroEnd() && p.cur.Literal() != ")" {
		tks = append(tks, p.cur)
		if p.cur.Literal() == "(" {
			paren++
		}
		p.next()
		if p.cur.Literal() == ")" && paren != 0 {
			tks = append(tks, p.cur)
			paren--
			p.next()
		}
	}
	return
}

// 获取下一个
func (p *processor) next() token.Token {
	p.cur = p.r.Scan()
	return p.cur
}

func (p *processor) addErr(pos token.Position, msg string, args ...interface{}) {
	p.err.Add(pos, msg, args...)
}

func (p *processor) addIfErr(err *Error) {
	if err != nil {
		p.err.Add(err.Pos, err.Msg)
	}
}

// 获取下一个非空token
func (p *processor) nextToken() token.Token {
	for {
		p.next()
		if p.cur.Type() != token.WHITESPACE {
			break
		}
	}
	return p.cur
}

// 获取下一个非空token
func (p *processor) skipWhitespace() token.Token {
	for {
		if p.cur.Type() != token.WHITESPACE {
			break
		}
		p.next()
	}
	return p.cur
}

func (p *processor) doMacro() {
	p.nextToken()
	switch p.cur.Literal() {
	case "if":
		p.next()
		cdt := p.evalConstExpr()
		p.expectEndMacro()
		if cdt {
			p.ctx.Push(IN_THEN)
		} else {
			// 跳到下一个分支
			p.ctx.Push(IN_ELSE)
			p.skipUtilElse()
		}
	case "ifdef":
		p.doIfDefine(true)
	case "ifndef":
		p.doIfDefine(false)
	case "elif":
		p.next()
		if p.ctx.Top() == IN_ELSE {
			cdt := p.evalConstExpr()
			p.expectEndMacro()
			if cdt {
				p.ctx.Pop()
				p.ctx.Push(IN_THEN)
			} else {
				// 跳到下一个分支
				p.skipUtilElse()
			}
		} else if p.ctx.Top() == IN_THEN {
			// 直接跳到结尾
			p.skipUtilCdt("endif")
			p.next() // endif
			p.expectEndMacro()
		} else {
			p.addErr(p.cur.Position(), "unexpected #else")
		}
	case "else":
		p.next()
		p.expectEndMacro()
		if p.ctx.Top() == IN_THEN {
			p.skipUtilCdt("endif")
			p.next()
			p.expectEndMacro()
		} else {
			p.addErr(p.cur.Position(), "unexpected #else")
		}
	case "endif":
		if p.ctx.Top() == IN_THEN || p.ctx.Top() == IN_ELSE {
			p.ctx.Pop()
			p.next() // endif
			p.expectEndMacro()
		} else {
			p.addErr(p.cur.Position(), "unexpected #endif")
		}
	case "define":
		p.doDefine()
	case "undef":
		p.doUndef()
	case "include":
		p.doInclude()
	case "pragma":
		p.doPragma()
	case "line":
		p.doLine()
	case "error":
		p.doError()
	default:
		p.skipEndMacro()
	}
}

// 重新压入token
func (p *processor) push(body []token.Token) {
	p.pushScanner(scanner.NewArrayScan(body))
}

// 重新压入token
func (p *processor) pushScanner(s scanner.Scanner) {
	var r scanner.MultiScanner
	if ps, ok := p.r.(scanner.MultiScanner); ok {
		r = ps
	} else {
		r = scanner.NewMultiScan(p.r)
	}
	r.Push(s)
	p.r = r
}

// 重新计算行内偏移量
func (p *processor) deltaLine(delta int) {
	c := p.startCache()
	p.skipEndMacro()
	arr := c.GetClear()
	columnDelta(arr, delta)
	p.push(arr)
}

// peek 下一个非空 token
func (p *processor) peekNext() token.Token {
	n := 1
	for {
		v := p.peek(n)
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

func (p *processor) peek(offset int) []token.Token {
	if ps, ok := p.r.(scanner.PeekScanner); ok {
		return ps.Peek(offset)
	}
	pp := scanner.NewPeekScan(p.r)
	tok := pp.Peek(offset)
	p.r = pp
	return tok
}

func (p *processor) expectIdent() string {
	if p.cur.Type() == token.IDENT {
		lit := p.cur.Literal()
		p.next()
		return lit
	}
	p.addErr(p.cur.Position(), fmt.Sprintf("expect token %s got %s", token.IDENT, p.cur.Type()))
	return ""
}

func (p *processor) expectPunctuator(lit string) {
	p.punctuator(lit, true)
}

func (p *processor) punctuator(lit string, require bool) {
	if p.cur.Type() == token.PUNCTUATOR && lit == p.cur.Literal() {
		p.nextToken()
		return
	}

	if require {
		p.addErr(p.cur.Position(), fmt.Sprintf("expect punctuator %s got %s", lit, p.cur.Literal()))
	}
}

func (p *processor) expectEndMacro() {
	if p.isMacroEnd() {
		p.nextToken()
		return
	}
	p.addErr(p.cur.Position(), fmt.Sprintf("expect end macro got %s", p.cur.Type()))
}

// 宏结尾
func (p *processor) isMacroEnd() bool {
	return p.cur.Type() == token.NEWLINE || p.cur.Type() == token.EOF
}

// 跳过无法到达的代码
func (p *processor) skipUtilCdt(names ...string) []token.Token {
	cdt := 0
	tks := make([]token.Token, 2)
	for {
		p.next()
		if p.cur.Type() == token.EOF {
			p.addErr(p.cur.Position(), fmt.Sprintf("expect %s, got EOF", strings.Join(names, ",")))
			break
		}
		if isMacroTok(p.cur) {
			tks[0] = p.cur
			p.nextToken()
			switch p.cur.Literal() {
			case "if", "ifndef", "ifdef":
				cdt++
			default:
				if cdt == 0 {
					for _, name := range names {
						if name == p.cur.Literal() {
							tks[1] = p.cur
							return tks
						}
					}
					if p.cur.Literal() == "endif" {
						return tks[0:0]
					}
				}
				if p.cur.Literal() == "endif" {
					cdt--
				}
			}
		}
	}
	return tks[0:0]
}

// #ifdef #ifndef
func (p *processor) doIfDefine(want bool) {
	p.nextToken()
	ident := p.expectIdent()
	cdt := p.ctx.IsDefined(ident)
	if cdt == want {
		p.ctx.Push(IN_THEN)
	} else {
		p.ctx.Push(IN_ELSE)
		p.skipUtilElse()
	}
	p.expectEndMacro()
}

// skip to #else/#elif
func (p *processor) skipUtilElse() {
	m := p.skipUtilCdt("elif", "else")
	if p.cur.Literal() == "elif" {
		p.next()  // elif
		p.push(m) // push back
	} else {
		p.next() // else
	}
}

func (p *processor) evalConstExpr() bool {
	var tks []token.Token
	for {
		if p.cur.Type() == token.EOF || p.cur.Type() == token.NEWLINE {
			break
		}
		if p.cur.Type() != token.WHITESPACE {
			tks = append(tks, p.cur)
		}
		p.next()
	}
	exp := newProcessor(p.ctx, scanner.NewArrayScan(tks), p.ignoreErr)
	expand, err := scanner.ScanToken(exp)
	if err != nil {
		p.addErr(p.cur.Position(), "invalid const-expr %s", inlineTokenString(tks))
	}
	return EvalConstExpr(p.ctx, expand)
}

func (p *processor) doDefine() {
	p.nextToken()
	ident := p.expectIdent()

	if p.ctx.IsDefined(ident) {
		p.addErr(p.cur.Position(), "duplicate define of %s", ident)
		p.skipEndMacro()
		return
	}

	if p.cur.Literal() == "(" {
		p.doDefineFunc(ident)
	} else {
		p.doDefineVal(ident)
	}
	p.expectEndMacro()
}

func (p *processor) skipEndMacro() {
	for !p.isMacroEnd() {
		p.nextToken()
	}
}

func (p *processor) doUndef() {
	p.nextToken()
	ident := p.expectIdent()
	delete(p.ctx.Val, ident)
	p.skipEndMacro()
	p.expectEndMacro()
}

func (p *processor) doDefineVal(ident string) {
	var tks []token.Token
	p.skipWhitespace()

	for !p.isMacroEnd() {
		tks = append(tks, p.cur)
		p.nextToken()
	}

	if err := p.ctx.DefineVal(ident, tks); err != nil {
		p.err.Add(err.(*Error).Pos, err.(*Error).Msg)
	}
}

func (p *processor) doDefineFunc(ident string) {
	var tks []token.Token
	var params []string

	p.expectPunctuator("(")

	elp := false
	for !p.isMacroEnd() && p.cur.Literal() != ")" {
		if p.cur.Literal() == "..." {
			elp = true
			p.nextToken()
			break
		} else if p.cur.Type() == token.IDENT {
			params = append(params, p.cur.Literal())
			p.nextToken()
			p.punctuator(",", false)
		} else {
			p.addErr(p.cur.Position(), fmt.Sprintf("expect ident, got %s <%s>", p.cur.Type(), p.cur.Literal()))
			break
		}
	}

	p.expectPunctuator(")")

	for !p.isMacroEnd() {
		tks = append(tks, p.cur)
		p.nextToken()
	}
	if err := p.ctx.DefineFunc(ident, params, elp, tks); err != nil {
		p.err.Add(err.(*Error).Pos, err.(*Error).Msg)
	}
}

func (p *processor) doInclude() {
	// include "file"
	if p.peekNext().Type() == token.STRING {
		p.nextToken()
		f, err := strconv.Unquote(p.cur.Literal())
		if err != nil {
			p.addErr(p.cur.Position(), "invalid include string %s", p.cur.Literal())
		}
		p.nextToken()
		p.skipEndMacro() // 跳到换行
		p.includeFile(f)
		return
	}

	// include <file>
	if p.peekNext().Literal() == "<" {
		p.nextToken()
		p.expectPunctuator("<")
		var f []token.Token
		for !p.isMacroEnd() && p.cur.Literal() != ">" {
			f = append(f, p.cur)
			p.next()
		}
		p.expectPunctuator(">")
		p.includeFile(relativeTokenString(f))
		return
	}

	// 不进行多次展开
	if _, ok := p.peekNext().(*Token); ok {
		p.addErr(p.nextToken().Position(), "invalid #include")
		return
	}

	// #include token...
	// 展开include后重新处理
	c := p.startCache()
	p.skipEndMacro()
	tks := c.GetClear()
	exp := newProcessor(p.ctx, scanner.NewArrayScan(tks), p.ignoreErr)
	expand, err := scanner.ScanToken(exp)
	if err != nil {
		p.addErr(p.cur.Position(), "invalid include string %s", inlineTokenString(tks))
	}
	p.push(expand)
	p.doInclude()
}

func (p *processor) searchFile(s string, tok token.Token) (string, bool) {
	return p.ctx.SearchFile(s, filepath.Dir(tok.Position().Filename))
}

func (p *processor) includeFile(s string) {
	if fn, ok := p.searchFile(s, p.cur); ok {
		if p.ctx.onceContain(fn) {
			p.skipEndMacro()
			p.expectEndMacro()
			return
		}
		sc, err := scanner.NewFileScan(fn)
		if err != nil {
			p.ctx.AddError(p.cur.Position(), "include %s: %s", fn, err.Error())
			return
		}
		p.push([]token.Token{p.cur})
		p.pushScanner(sc)
		p.next()
	} else {
		p.ctx.AddError(p.cur.Position(), "file not found %s", s)
	}
}
func (p *processor) startCache() scanner.CachedScanner {
	r := scanner.NewCachedScanner(p.r)
	r.Start()
	p.r = r
	return r
}

func (p *processor) doPragma() {
	c := p.startCache()
	p.skipEndMacro()
	pp := c.GetClear()
	for _, v := range pp {
		// 支持 pragma once 指令
		if v.Literal() == "once" {
			pp, _ := filepath.Abs(v.Position().Filename)
			p.ctx.pragmaOnce(pp)
		}
	}
	p.expectEndMacro()
}

func (p *processor) doLine() {
	c := p.startCache()
	file := p.cur.Position().Filename
	line := p.cur.Position().Line
	p.skipEndMacro()
	v := c.GetClear()
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
		p.r = mockLine(p.r, p.cur.Position(), file, line)
	}
	p.expectEndMacro()
}

func (p *processor) doError() {
	pos := p.cur.Position()
	p.next() // error
	c := p.startCache()
	p.skipEndMacro()
	msg := c.GetClear()
	p.expectEndMacro()
	p.addErr(pos, inlineTokenString(msg))
}

type lineDirective struct {
	scanner.Scanner
	scope string
	delta int
	file  string
}

func (ld *lineDirective) Scan() token.Token {
	t := ld.Scanner.Scan()
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

func mockLine(s scanner.Scanner, pos token.Position, file string, line int) scanner.Scanner {
	ld := &lineDirective{}
	ld.delta = line - pos.Line
	ld.file = file
	ld.scope = pos.Filename
	ld.Scanner = s
	return ld
}

type expandMock struct {
	r   scanner.Scanner
	tok token.Token
}

func (em *expandMock) Scan() token.Token {
	t := em.r.Scan()
	return &Token{
		Pos:    t.Position(),
		Typ:    t.Type(),
		Lit:    t.Literal(),
		Expand: em.tok,
	}
}

func newExpandMock(tok token.Token, r scanner.Scanner) scanner.Scanner {
	return &expandMock{r, tok}
}

// NewScanner 预处理扫描器
func NewScanner(ctx *Context, s scanner.Scanner, ignore bool) scanner.Scanner {
	return scanner.NewTokenScan(newProcessor(ctx, s, ignore))
}
