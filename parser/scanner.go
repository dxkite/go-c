package parser

import (
	"dxkite.cn/c/scanner"
	"dxkite.cn/c/token"
)

type fileScanner struct {
	r    scanner.PeekScanner
	file string
}

type Token struct {
	Pos token.Position
	Typ token.Type
	Lit string
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

func newFileScanner(file string, r scanner.PeekScanner) *fileScanner {
	return &fileScanner{file: file, r: r}
}

func (f *fileScanner) Scan() token.Token {
	if t := f.r.PeekOne(); t.Position().Filename == f.file {
		return f.r.Scan()
	}
	return &Token{
		Pos: token.Position{},
		Typ: token.EOF,
		Lit: "",
	}
}
