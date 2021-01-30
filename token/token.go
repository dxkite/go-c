package token

// 文件内位置
type Position struct {
	// 文件路径
	Filename string
	// 行,列
	Line, Column int
	// 完整偏移量
	Offset int
}

// token
type Token struct {
	Position Position
	Type     Type
	Lit      string
}
