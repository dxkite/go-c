package scanner

import (
	"bytes"
	"dxkite.cn/c/token"
	"io"
	"os"
)

type Scanner interface {
	// 获取下一个token
	Scan() token.Token
}

type PeekScanner interface {
	Scanner
	Peek(offset int) []token.Token
	PeekOne() token.Token
}

func ScanToken(s Scanner) ([]token.Token, error) {
	tks := make([]token.Token, 0)
	for {
		tok := s.Scan()
		if tok.Type() == token.EOF {
			break
		}
		// 合并连续空白
		if i := len(tks) - 1; tok.Type() == token.WHITESPACE && i >= 0 && tks[i].Type() == token.WHITESPACE {
			continue
		}
		tks = append(tks, tok)
		if v, ok := tok.(*IllegalToken); ok {
			return tks, v.Err
		}
	}
	return tks, nil
}

func Scan(filename string, r io.Reader) ([]token.Token, error) {
	s := NewScan(filename, r)
	return ScanToken(s)
}

func ScanString(name, code string) ([]token.Token, error) {
	return Scan(name, bytes.NewBufferString(code))
}

func ScanFile(filename string) ([]token.Token, error) {
	f, er := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if er != nil {
		return nil, er
	}
	defer func() { _ = f.Close() }()
	return Scan(filename, f)
}

type arrayScanner struct {
	arr []token.Token
	off int
}

func NewArrayScan(tok []token.Token) Scanner {
	return &arrayScanner{
		arr: tok,
		off: 0,
	}
}

func (a *arrayScanner) Scan() (t token.Token) {
	of := a.off
	a.off++
	if of < len(a.arr) {
		return a.arr[of]
	}
	return &Token{
		Typ: token.EOF,
	}
}

type peekScanner struct {
	c   []token.Token
	s   Scanner
	off int
}

func NewPeekScan(s Scanner) PeekScanner {
	return &peekScanner{
		c:   []token.Token{},
		s:   s,
		off: 0,
	}
}

func (s *peekScanner) Scan() (t token.Token) {
	of := s.off
	s.off++
	if of < len(s.c) {
		return s.c[of]
	}
	s.c = s.c[0:0]
	s.off = 0
	return s.s.Scan()
}

func (s *peekScanner) Peek(n int) (t []token.Token) {
	of := s.off
	t = []token.Token{}
	lc := len(s.c)
	for n > len(t) && of < lc {
		t = append(t, s.c[of])
		of++
	}
	for n > len(t) {
		r := s.s.Scan()
		t = append(t, r)
		s.c = append(s.c, r)
		if r.Type() == token.EOF {
			break
		}
	}
	return
}

func (s *peekScanner) PeekOne() token.Token {
	return s.Peek(1)[0]
}

type MultiScanner interface {
	Scanner
	Push(s Scanner)
}

type multiScanner struct {
	s   []Scanner
	cur int
}

func NewMultiScan(s ...Scanner) MultiScanner {
	return &multiScanner{
		s:   s,
		cur: len(s) - 1,
	}
}

func (ms *multiScanner) Scan() (t token.Token) {
	for ms.cur >= 0 {
		t = ms.s[ms.cur].Scan()
		if t.Type() != token.EOF {
			break
		}
		if ms.cur > 0 {
			ms.cur--
		} else if ms.cur == 0 {
			return t
		}
	}
	return
}

func (ms *multiScanner) Push(s Scanner) {
	ms.cur++
	if ms.cur >= len(ms.s) {
		ms.s = append(ms.s, s)
	} else {
		ms.s[ms.cur] = s
	}
}

type tokenScanner struct {
	Scanner
}

// 扫描字符串
// 跳过空白符
func NewTokenScan(s Scanner) Scanner {
	return &tokenScanner{s}
}

func (ts *tokenScanner) Scan() (t token.Token) {
	for t = ts.Scanner.Scan(); t.Type() == token.WHITESPACE; t = ts.Scanner.Scan() {
		// next
	}
	return t
}

type CachedScanner interface {
	Scanner
	Start()
	Clear()
	GetClear() (v []token.Token)
	Restore()
}

type cachedScanner struct {
	enable bool
	cached []token.Token
	r      Scanner
}

func NewCachedScanner(r Scanner) CachedScanner {
	c := &cachedScanner{r: r}
	c.cached = []token.Token{}
	return c
}

func (r *cachedScanner) Scan() (t token.Token) {
	t = r.r.Scan()
	if r.enable && t.Type() != token.EOF {
		r.cached = append(r.cached, t)
	}
	return
}

func (r *cachedScanner) Start() {
	r.enable = true
}

// 清空缓存
func (r *cachedScanner) Clear() {
	r.enable = false
	r.cached = r.cached[0:0]
}

// 清空缓并获取内容
func (r *cachedScanner) GetClear() (v []token.Token) {
	r.enable = false
	v = append(v, r.cached...)
	r.cached = r.cached[0:0]
	return
}

// 重置指针位置
func (r *cachedScanner) Restore() {
	r.enable = false
	r.r = NewMultiScan(NewArrayScan(r.cached), r.r)
}

type fileScanner struct {
	Scanner
	closed bool
	f      io.ReadCloser
	name   string
	eof    token.Token
}

func NewFileScan(filename string) (Scanner, error) {
	s := &fileScanner{}
	s.name = filename
	f, er := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if er != nil {
		return nil, er
	}
	s.f = f
	s.Scanner = NewScan(filename, f)
	return s, nil
}

func (s *fileScanner) Scan() token.Token {
	t := s.Scanner.Scan()
	if t.Type() == token.EOF {
		if err := s.f.Close(); err == nil {
			s.closed = true
		} else {
			return &IllegalToken{
				Token: &Token{
					Pos: t.Position(),
					Typ: t.Type(),
					Lit: t.Literal(),
				},
				Err: err,
			}
		}
		return t
	}
	if s.closed {
		return s.eof
	}
	return t
}
