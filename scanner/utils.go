package scanner

import (
	"bytes"
	"dxkite.cn/go-c11/token"
	"encoding/json"
	"io/ioutil"
	"os"
)

func EncodingJson(tks []token.Token) ([]byte, error) {
	buf := &bytes.Buffer{}
	je := json.NewEncoder(buf)
	je.SetEscapeHTML(false)
	je.SetIndent("", "    ")
	if err := je.Encode(tks); err != nil {
		return nil, err
	} else {
		return buf.Bytes(), nil
	}
}

func SaveJson(filename string, tks []token.Token) error {
	if b, err := EncodingJson(tks); err != nil {
		return err
	} else {
		return ioutil.WriteFile(filename, b, os.ModePerm)
	}
}
