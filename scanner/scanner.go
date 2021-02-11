package scanner

import (
	"bufio"
	"bytes"
	"dxkite.cn/go-c11"
	"dxkite.cn/go-c11/token"
	"fmt"
	"io"
	"os"
	"unicode"
	"unicode/utf8"
)

type Scanner interface {
	Scan() (t token.Token)
	Error() *go_c11.ErrorList
}

type PeekScanner interface {
	Scanner
	Peek(offset int) []token.Token
	PeekOne() token.Token
}

func NewScan(filename string, r io.Reader) Scanner {
	s := &scanner{}
	s.init(filename, r)
	return s
}

func NewStringScan(filename string, code string) Scanner {
	s := &scanner{}
	s.init(filename, bytes.NewBufferString(code))
	return s
}

type scanner struct {
	filename  string
	r         *bufio.Reader
	ch        rune
	offset    int
	rdOffset  int
	line, col int
	err       go_c11.ErrorList
	rcd       bool
	lit       string
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

// init
func (s *scanner) init(filename string, r io.Reader) {
	s.r = bufio.NewReader(r)
	s.ch = ' '
	s.offset = 0
	s.filename = filename
	s.line = 1
	s.col = 0
	s.err.Reset()
	s.nextRune()
}

func (s *scanner) record() {
	s.rcd = true
	s.lit = ""
}

func (s *scanner) literal() string {
	s.rcd = false
	return s.lit
}

// 获取下一个字符
func (s *scanner) next() (rune, token.Position) {
	cur := s.ch
	p := s.curPos()
	s.nextRune()

	if s.ch == '\\' && (s.peekN(2) == "\r\n" || s.peek() == '\n') {
		s.nextRune()
		s.nextRune()
	}

	if cur != -1 && s.rcd {
		s.lit += string(cur)
	}
	return cur, p
}

func (s *scanner) nextRune() {
	if s.ch == '\n' {
		s.line++
		s.col = 1
	} else {
		s.col++
	}
	s.offset = s.rdOffset

	ch, w, err := s.r.ReadRune()
	if err != nil {
		if err != io.EOF {
			s.error(s.curPos(), fmt.Sprintf("reader error %s", err))
		}
		s.ch = -1
		return
	}

	if ch == '\r' {
		if s.peek() == '\n' {
			w++
			_, _ = s.r.ReadByte()
		}
		ch = '\n'
	}

	s.rdOffset += w
	s.ch = ch
	return
}

func (s *scanner) peek() byte {
	buf, err := s.r.Peek(1)
	if err != nil && err != io.EOF {
		s.error(s.curPos(), fmt.Sprintf("reader error %s", err))
		return 0
	}
	return buf[0]
}

func (s *scanner) peekN(n int) string {
	buf, err := s.r.Peek(n)
	if err != nil && err != io.EOF {
		s.error(s.curPos(), fmt.Sprintf("reader error %s", err))
		return ""
	}
	return string(buf)
}

func (s *scanner) peekCN(ch string, n int) string {
	return ch + s.peekN(n)
}

func (s *scanner) curPos() token.Position {
	return token.Position{
		Filename: s.filename,
		Line:     s.line,
		Column:   s.col,
	}
}

func (s *scanner) error(p token.Position, msg string) {
	s.err.Add(p, msg)
}

func (s *scanner) Scan() token.Token {
	t := &Token{}
	t.Pos = s.curPos()
	t.Typ = token.ILLEGAL
	switch ch := s.ch; {
	case isWhitespace(ch):
		t.Typ = token.WHITESPACE
		t.Lit = " "
		s.skipWhitespace()
	case ch == '/' && (s.peek() == '/' || s.peek() == '*'):
		t.Typ = token.WHITESPACE
		t.Lit = " "
		s.skipComment()
	case s.nextIsChar(ch):
		t.Typ = token.CHAR
		t.Lit = s.scanChar()
	case s.nextIsString(ch):
		t.Typ = token.STRING
		t.Lit = s.scanString()
	case isLetter(ch):
		t.Typ = token.IDENT
		t.Lit = s.scanIdentifier()
		switch t.Lit {
		case "auto", "break", "case", "char",
			"const", "continue", "default", "do",
			"double", "else", "enum", "extern", "float",
			"for", "goto", "if", "inline", "int", "long",
			"register", "restrict", "return", "short", "signed",
			"sizeof", "static", "struct", "switch", "typedef",
			"union", "unsigned", "void", "volatile", "while",
			"_Alignas", "_Alignof", "_Atomic", "_Bool", "_Complex",
			"_Generic", "_Imaginary", "_Noreturn", "_Static_assert",
			"_Thread_local":
			t.Typ = token.KEYWORD
		}
	case s.nextIsNumber():
		t.Typ, t.Lit = s.scanNumber()
	default:
		if lit, n, ok := s.nextIsPunctuator(); ok {
			t.Lit = lit
			t.Typ = token.PUNCTUATOR
			for n > 0 {
				n--
				s.next()
			}
		} else {
			s.next()
			switch ch {
			case -1:
				t.Typ = token.EOF
			case '\n':
				t.Typ = token.NEWLINE
				t.Lit = "\n"
			default:
				t.Lit = string(ch)
				t.Typ = token.TEXT
			}
		}
	}
	return t
}

func (s *scanner) Error() *go_c11.ErrorList {
	return &s.err
}

func isWhitespace(ch rune) bool {
	switch ch {
	case ' ', '\t', '\r':
		return true
	default:
		return false
	}
}

func (s *scanner) skipWhitespace() bool {
	c := 0
	for isWhitespace(s.ch) {
		c++
		s.next()
	}
	return c > 0
}

// 扫描标识符
func (s *scanner) scanIdentifier() string {
	s.record()
	for isLetter(s.ch) || isDigit(s.ch) {
		s.next()
	}
	return s.literal()
}

func lower(ch rune) rune { return ('a' - 'A') | ch }

func isDecimal(ch rune) bool { return '0' <= ch && ch <= '9' }

func isHex(ch rune) bool { return '0' <= ch && ch <= '9' || 'a' <= lower(ch) && lower(ch) <= 'f' }

func isOct(ch rune) bool { return '0' <= ch && ch <= '7' }

func isLetter(ch rune) bool {
	return 'a' <= lower(ch) && lower(ch) <= 'z' || ch == '_' || ch >= utf8.RuneSelf && unicode.IsLetter(ch)
}

func isDigit(ch rune) bool {
	return isDecimal(ch) || ch >= utf8.RuneSelf && unicode.IsDigit(ch)
}

func (s *scanner) nextIsString(ch rune) bool {
	// "u8" | "u" | "U" | "L"
	if ch == 'u' && s.peekN(2) == `8"` {
		return true
	}
	switch ch {
	case 'u', 'U', 'L':
		return s.peek() == '"'
	}
	return ch == '"'
}

func (s *scanner) nextIsChar(ch rune) bool {
	// [ "L" | "u" | "U" ] "'" c-char-sequence "'"
	switch ch {
	case 'u', 'U', 'L':
		return s.peek() == '\''
	}
	return ch == '\''
}

// 扫描字符串
func (s *scanner) scanChar() string {
	return s.scanQuote("char", '\'')
}

// 扫描字符串
func (s *scanner) scanString() string {
	return s.scanQuote("string", '"')
}

// 扫描字符串
func (s *scanner) scanQuote(name string, quote rune) string {
	s.record()
	for s.ch != quote {
		s.next()
	}
	s.next()
	for s.ch > 0 && s.ch != '\n' && s.ch != quote {
		if s.ch == '\\' {
			s.next()
			s.scanEscape()
		} else {
			s.next()
		}
		if quote == '\'' {
			break
		}
	}
	if s.ch != quote {
		s.error(s.curPos(), fmt.Sprintf("unclosed %s lit %c", name, quote))
	} else {
		s.next()
	}
	return s.literal()
}

// 扫描字符串
func (s *scanner) scanEscape() bool {
	switch s.ch {
	// simple-escape-sequence
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', '"', '\'':
		s.next()
		return true
	case '0', '1', '2', '3', '4', '5', '6', '7':
		s.next()
		n := 2
		for n > 0 && isOct(s.ch) {
			n--
			s.next()
		}
	case 'x':
		s.next()
		n := 2
		for n > 0 {
			if !isHex(s.ch) {
				s.error(s.curPos(), fmt.Sprintf("unexpected %c in hex escape", s.ch))
				return false
			}
			n--
			s.next()
		}
	case 'u', 'U':
		return s.scanUniversalEscape()
	}
	return true
}

func (s *scanner) scanUniversalEscape() bool {
	n := 4
	if s.ch == 'U' {
		n = 8
	}
	s.next() // skip u
	for n > 0 {
		if !isHex(s.ch) {
			s.error(s.curPos(), fmt.Sprintf("unexpected %c in unicode escape", s.ch))
			return false
		}
		n--
		s.next()
	}
	return true
}

func (s *scanner) skipComment() {
	s.next()
	if s.ch == '/' {
		s.next()
		for s.ch != '\n' && s.ch >= 0 {
			s.next()
		}
		return
	}

	if s.ch == '*' {
		s.next()
		for {
			if s.ch < 0 {
				s.error(s.curPos(), "comment not terminated")
				return
			}
			s.next()
			if s.ch == '*' && s.peek() == '/' {
				break
			}
		}
		s.next() // *
		s.next() // /
		return
	}
}

func (s *scanner) nextIsNumber() bool {
	if isDigit(s.ch) {
		return true
	}
	if s.ch == '.' && isDigit(rune(s.peek())) {
		return true
	}
	return false
}

func (s *scanner) scanNumber() (token.Type, string) {
	s.record()
	s.next()
	typ := token.INT

	base := 10

	switch s.ch {
	case '0':
		if lower(rune(s.peek())) == 'x' {
			base = 16
		} else {
			base = 8
		}
	}

	s.scanNumberBase(base)
	if s.ch == '.' {
		typ = token.FLOAT
		s.next()
		s.scanNumberBase(base)
	}

	if ch := lower(s.ch); ch == 'e' || ch == 'p' {
		typ = token.FLOAT
		s.next()
		if s.ch == '+' || s.ch == '-' {
			s.next()
		}
		s.scanNumberBase(10)
	}

	if ch := lower(s.ch); ch == 'l' || ch == 'f' {
		typ = token.FLOAT
		s.next()
	} else {
		s.scanIntSuffix()
	}

	return typ, s.literal()
}

func (s *scanner) scanIntSuffix() {
	// u
	// ul
	// ull
	if lower(s.ch) == 'u' {
		s.next()
		n := 2
		for n > 0 && lower(s.ch) == 'l' {
			s.next()
			n--
		}
	}

	// ll l llu lu
	if lower(s.ch) == 'l' {
		if lower(s.ch) == 'l' {
			s.next()
		}
		if lower(s.ch) == 'u' {
			s.next()
		}
	}
	return
}

func (s *scanner) scanNumberBase(base int) {
	if base <= 10 {
		for isDecimal(s.ch) {
			s.next()
		}
	} else {
		for isHex(s.ch) {
			s.next()
		}
	}
	return
}

var mp = map[string]string{
	"<:":   "[",
	":>":   "]",
	"<%":   "{",
	"%>":   "}",
	"%:":   "#",
	"%:%:": "##",
}

func (s *scanner) nextIsPunctuator() (string, int, bool) {
	tok := s.peekCN(string(s.ch), 3)
	for i := len(tok); i > 0; i-- {
		switch ch := tok[:i]; ch {
		case "...", ".", ",", "?", ":", ";",
			"[", "]", "(", ")", "{", "}", "~",
			"->", "--", "-=", "-",
			"++", "+=", "+",
			"&=", "&&", "&",
			"*=", "*",
			"!", "!=",
			"==", "=",
			"^=", "^",
			"/=", "/",
			"%=", "%:%:", "%:", "%",
			"||", "|=", "|",
			"<<=", ">>=", "<<", ">>", "<:", ":>", "<%", "%>", "<=", ">=", "<", ">",
			"##", "#":
			if v, ok := mp[ch]; ok {
				return v, i, true
			}
			return ch, i, true
		}
	}
	return "", 0, false
}

func ScanToken(s Scanner) []token.Token {
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
	}
	return tks
}

func Scan(filename string, r io.Reader) ([]token.Token, *go_c11.ErrorList) {
	s := NewScan(filename, r)
	tks := ScanToken(s)
	return tks, s.Error()
}

func ScanString(name, code string) ([]token.Token, *go_c11.ErrorList) {
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
	err go_c11.ErrorList
}

func NewArrayScan(tok []token.Token) Scanner {
	return &arrayScanner{
		arr: tok,
		off: 0,
		err: go_c11.ErrorList{},
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

func (a *arrayScanner) Error() *go_c11.ErrorList {
	return &a.err
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

func (s *peekScanner) Error() *go_c11.ErrorList {
	return s.s.Error()
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
			err := ms.s[ms.cur].Error()
			ms.cur--
			ms.s[ms.cur].Error().Merge(*err)
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

func (ms *multiScanner) Error() *go_c11.ErrorList {
	return ms.s[ms.cur].Error()
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

type fileScanner struct {
	Scanner
	closed bool
	f      io.ReadCloser
	name   string
	eof    token.Token
}

func NewFileScan(filename string) Scanner {
	s := &fileScanner{}
	s.name = filename
	f, er := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if er != nil {
		s.Error().Add(token.Position{}, "open file error %s", filename)
	}
	s.f = f
	s.Scanner = NewScan(filename, f)
	return s
}

func (s *fileScanner) Scan() token.Token {
	t := s.Scanner.Scan()
	if t.Type() == token.EOF {
		if err := s.f.Close(); err == nil {
			s.closed = true
		} else {
			s.Error().Add(token.Position{}, "close file error %s", s.name)
		}
		return t
	}
	if s.closed {
		return s.eof
	}
	return t
}
