package errors

import (
	"dxkite.cn/c/token"
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// 扫描错误
type Error struct {
	Pos    token.Position
	Msg    string
	Code   ErrCode
	Params []interface{}
}

// 错误信息
func (e Error) Error() string {
	t := fmt.Sprintf("在 %s 文件的第%d行%d列: ", e.Pos.Filename, e.Pos.Line, e.Pos.Column)
	if e.Code != ErrUnKnown {
		s := t + e.Code.String()
		if len(e.Params) > 0 {
			return fmt.Sprintf(s, e.Params...)
		}
		return s
	}
	return t + e.Msg
}

// 错误列表
type ErrorList []*Error

// 添加一个错误
func (p *ErrorList) Add(pos token.Position, msg string, args ...interface{}) {
	msg = fmt.Sprintf(msg, args...)
	*p = append(*p, &Error{Pos: pos, Msg: msg})
}

func (p *ErrorList) AddErr(err *Error) {
	*p = append(*p, err)
}

func (p *ErrorList) AddErrMsg(pos token.Position, code ErrCode, params ...interface{}) {
	*p = append(*p, New(pos, code, params...))
}

func (p *ErrorList) AddStdErr(pos token.Position, err error) {
	if v, ok := err.(*Error); ok {
		*p = append(*p, v)
	} else {
		p.Add(pos, err.Error())
	}
}

func New(pos token.Position, code ErrCode, params ...interface{}) *Error {
	return &Error{Pos: pos, Code: code, Params: params}
}

func NewStd(pos token.Position, err error) *Error {
	return &Error{Pos: pos, Msg: err.Error()}
}

func NewMsg(pos token.Position, msg string, args ...interface{}) *Error {
	msg = fmt.Sprintf(msg, args...)
	return &Error{Pos: pos, Msg: msg}
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
		return "无错误"
	case 1:
		return p[0].Error()
	}
	return fmt.Sprintf("%s (共 %d 个错误)", p[0], len(p)-1)
}

func (p ErrorList) SaveFile(filename string) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	je := json.NewEncoder(f)
	je.SetEscapeHTML(false)
	je.SetIndent("", "    ")
	if err := je.Encode(p); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}

func (p *ErrorList) LoadFromFile(filename string) error {
	f, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	je := json.NewDecoder(f)
	if err := je.Decode(&p); err != nil {
		_ = f.Close()
		return err
	} else {
		if p == nil {
			p = &ErrorList{}
		}
		return f.Close()
	}
}

type ErrorType int

const (
	ErrTypeError ErrorType = iota
	ErrTypeWarning
)

// 错误回调
type ErrorHandler func(pos token.Position, typ ErrorType, code ErrCode, params ...interface{})
