package token

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
)

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

func SaveJson(filename string, tks []*Token) error {
	buf := &bytes.Buffer{}
	je := json.NewEncoder(buf)
	je.SetEscapeHTML(false)
	je.SetIndent("", "    ");
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
