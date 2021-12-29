package scanner

import (
	"bufio"
	"bytes"
	"dxkite.cn/c/errors"
	"dxkite.cn/c/token"
	"io"
	"unicode"
	"unicode/utf8"
)

// 普通token
type Token struct {
	Pos token.Position
	Typ token.Type
	Lit string
}

type Option struct {
	// 全角符号转半角符号
	PunctuatorFullWidthToHalfWidth bool
}

// 非法token
type IllegalToken struct {
	*Token
	Err error
}

// 全角符号转换的普通符号
type FullWidthPunctuatorToken struct {
	*Token
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

func (t *Token) String() string {
	return token.String(t)
}

func NewScan(filename string, r io.Reader, option *Option) Scanner {
	s := &scanner{}
	if option == nil {
		option = &Option{}
	}
	s.opt = option
	s.new(filename, r)
	return s
}

func NewStringScan(filename string, code string, option *Option) Scanner {
	return NewScan(filename, bytes.NewBufferString(code), option)
}

type scanner struct {
	filename  string
	r         *bufio.Reader
	ch        rune
	offset    int
	rdOffset  int
	line, col int
	rcd       bool
	lit       string
	err       error
	opt       *Option
}

// new
func (s *scanner) new(filename string, r io.Reader) {
	s.r = bufio.NewReader(r)
	s.ch = ' '
	s.offset = 0
	s.filename = filename
	s.line = 1
	s.col = 0
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
			s.markErr(errors.ErrReadFile, err)
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
		s.markErr(errors.ErrReadFile, err)
		return 0
	}
	if err == io.EOF {
		return 0
	}
	return buf[0]
}

func (s *scanner) markErr(code errors.ErrCode, params ...interface{}) {
	s.err = errors.New(s.curPos(), code, params...)
}

func (s *scanner) peekN(n int) string {
	buf, err := s.r.Peek(n)
	if err != nil && err != io.EOF {
		s.markErr(errors.ErrReadFile, err)
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

func (s *scanner) Scan() token.Token {
	t := &Token{}
	t.Pos = s.curPos()
	t.Typ = token.ILLEGAL
	s.err = nil
	fullWidth := false
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
		case "auto", "break", "case", "char", "const", "continue", "default", "do", "double", "else",
			"enum", "extern", "float", "for", "goto", "if", "inline", "int", "long", "register",
			"restrict", "return", "short", "signed", "sizeof", "static", "struct", "switch", "typedef", "union",
			"unsigned", "void", "volatile", "while", "_Bool", "_Complex", "_Imaginary":
			t.Typ = token.KEYWORD
		}
	case s.nextIsNumber():
		t.Typ, t.Lit = s.scanNumber()
	default:
		if lit, n, full, ok := s.nextIsPunctuator(); ok {
			t.Lit = lit
			t.Typ = token.PUNCTUATOR
			fullWidth = full
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
	if s.err != nil {
		return &IllegalToken{
			Token: t,
			Err:   s.err,
		}
	}
	if fullWidth {
		return &FullWidthPunctuatorToken{Token: t}
	}
	return t
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
	return s.scanQuote(errors.ErrScanUncloseChar, '\'')
}

// 扫描字符串
func (s *scanner) scanString() string {
	return s.scanQuote(errors.ErrScanUncloseString, '"')
}

// 扫描字符串
func (s *scanner) scanQuote(err errors.ErrCode, quote rune) string {
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
		s.markErr(err)
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
				s.markErr(errors.ErrScanHexFormat, s.ch)
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
			s.markErr(errors.ErrScanUnicodeFormat, s.ch)
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
				s.markErr(errors.ErrScanUncloseComment)
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

func toHalfWidthPunctuator(s string) (string, bool) {
	r := make([]rune, len(s))
	t := false
	for i, v := range s {
		if v >= 0xff01 && v <= 0xff5e {
			v -= 0xff00 - 0x20
			t = true
		}
		r[i] = v
	}
	return string(r), t
}

func (s *scanner) nextIsPunctuator() (string, int, bool, bool) {
	tok := s.peekCN(string(s.ch), 3)
	trans := false
	if s.opt.PunctuatorFullWidthToHalfWidth {
		if v, ok := toHalfWidthPunctuator(tok); ok {
			tok = v
			trans = true
		}
	}
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
				return v, i, trans, true
			}
			return ch, i, trans, true
		}
	}
	return "", 0, trans, false
}
