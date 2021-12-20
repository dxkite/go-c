package preprocess

import (
	"dxkite.cn/c/token"
	"strings"
	"unicode/utf8"
)

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
