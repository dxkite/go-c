package go_c11

import (
	"fmt"
	"sort"

	"dxkite.cn/go-c11/token"
)

// 扫描错误
type Error struct {
	Pos token.Position
	Msg string
}

// 错误信息
func (e Error) Error() string {
	return e.Msg
}

// 错误列表
type ErrorList []*Error

// 添加一个错误
func (p *ErrorList) Add(pos token.Position, msg string) {
	*p = append(*p, &Error{pos, msg})
}

// 合并错误
func (p *ErrorList) Merge(err ErrorList) {
	*p = append(*p, err...)
}

// 清空错误
func (p *ErrorList) Reset() { *p = (*p)[0:0] }

// 排序接口
func (p ErrorList) Len() int      { return len(p) }
func (p ErrorList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p ErrorList) Less(i, j int) bool {
	if p[i].Pos.Filename != p[j].Pos.Filename {
		return len(p[i].Pos.Filename) < len(p[j].Pos.Filename)
	}
	return p[i].Pos.Column < p[j].Pos.Column
}

// 排序输出
func (p ErrorList) Sort() {
	sort.Sort(p)
}

// 输出错误
func (p ErrorList) Error() string {
	switch len(p) {
	case 0:
		return "no errors"
	case 1:
		return p[0].Error()
	}
	return fmt.Sprintf("%s (and %d more errors)", p[0], len(p)-1)
}
