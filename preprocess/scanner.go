package preprocess

import (
	"dxkite.cn/c/scanner"
	"dxkite.cn/c/token"
)

// 扫描器
// 自动去除后置空白
type noSpaceScanner struct {
	r scanner.PeekScanner
}

// 扫描字符串
// 跳过空白符
func NewScan(s scanner.Scanner) scanner.Scanner {
	return &noSpaceScanner{scanner.NewPeekScan(s)}
}

func (ts *noSpaceScanner) Scan() (t token.Token) {
	t = ts.r.Scan()
	space := false
	for ts.r.PeekOne().Type() == token.WHITESPACE {
		space = true
		ts.r.Scan() // whitespace
	}
	return &Token{
		Pos:   t.Position(),
		Typ:   t.Type(),
		Lit:   t.Literal(),
		Space: space,
	}
}
