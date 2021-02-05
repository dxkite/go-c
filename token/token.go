package token

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"unicode/utf8"
)

// 文件内位置
type Position struct {
	// 文件路径
	Filename string
	// 行,列
	Line, Column int
}

// token
type Token struct {
	Position Position
	Type     Type
	Lit      string
}

func SaveJson(filename string, tks []*Token) error {
	buf := &bytes.Buffer{}
	je := json.NewEncoder(buf)
	je.SetEscapeHTML(false)
	je.SetIndent("", "    ")
	if err := je.Encode(tks); err != nil {
		return err
	} else {
		return ioutil.WriteFile(filename, buf.Bytes(), os.ModePerm)
	}
}

func LoadJson(filename string) ([]*Token, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	tks := []*Token{}
	if err := json.Unmarshal(buf, &tks); err != nil {
		return nil, err
	} else {
		return tks, nil
	}
}

func String(tks []*Token) string {
	str := ""
	col := 1
	line := 1
	for _, tok := range tks {
		// 换行
		if tok.Type == NEWLINE {
			line++
			col = 1
			str += "\n"
			continue
		}

		// 行
		if tok.Position.Line != line {
			if d := tok.Position.Line - line; d > 0 {
				str += strings.Repeat("\n", d)
				line = tok.Position.Line
				col = 1
			}
		}

		// 列
		if tok.Position.Column != col {
			if d := tok.Position.Column - col; d > 0 {
				str += strings.Repeat(" ", d)
				col = tok.Position.Column
			}
		}

		col = col + utf8.RuneCountInString(tok.Lit)
		str += tok.Lit
	}
	return str
}
