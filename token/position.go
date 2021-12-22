package token

import "fmt"

// 文件内位置
type Position struct {
	// 文件路径
	Filename string
	// 行,列
	Line, Column int
}

func (p *Position) String() string {
	return fmt.Sprintf("%s:%d:%d", p.Filename, p.Line, p.Column)
}
