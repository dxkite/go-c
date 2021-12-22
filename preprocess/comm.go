package preprocess

import (
	"dxkite.cn/c/scanner"
	"dxkite.cn/c/token"
	"os"
	"strings"
	"unicode/utf8"
)

// 宏可变参数 __VA_ARGS__
const MacroParameterVarArgs = "__VA_ARGS__"

func tokenString(tks []token.Token) string {
	str := ""
	col := 1
	line := 1
	for _, tok := range tks {
		// 换行
		if tok.Type() == token.NEWLINE {
			line++
			col = 1
			str += "\n"
			continue
		}

		// 行
		if tok.Position().Line != line {
			if d := tok.Position().Line - line; d > 0 {
				str += strings.Repeat("\n", d)
				line = tok.Position().Line
				col = 1
			}
		}

		// 列
		if tok.Position().Column != col {
			if d := tok.Position().Column - col; d > 0 {
				str += strings.Repeat(" ", d)
				col = tok.Position().Column
			}
		}

		col = col + utf8.RuneCountInString(tok.Literal())
		str += tok.Literal()
	}
	return str
}

func relativeTokenString(tks []token.Token) string {
	return strings.TrimSpace(tokenString(tks))
}

func inlineTokenString(tks []token.Token) string {
	return strings.ReplaceAll(relativeTokenString(tks), "\n", "")
}

func isMacroTok(tok token.Token) bool {
	// 展开后的#不作为宏符号
	if v, ok := tok.(*Token); ok && v.Expand != nil {
		return false
	}
	return tok.Position().Column == 1 && tok.Literal() == "#"
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

func isValidToken(lit string) bool {
	s := scanner.NewStringScan("<runtime>", lit)
	if tks, err := scanner.ScanToken(s); err != nil {
		return false
	} else {
		return len(tks) == 1
	}
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
	m[MacroParameterVarArgs] = true
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
