package scanner

import (
	"bytes"
	"dxkite.cn/c/token"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func encodingJson(tks []token.Token) ([]byte, error) {
	buf := &bytes.Buffer{}
	je := json.NewEncoder(buf)
	je.SetEscapeHTML(false)
	je.SetIndent("", "  ")
	if err := je.Encode(tks); err != nil {
		return nil, err
	} else {
		return buf.Bytes(), nil
	}
}

func (tk *IllegalToken) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		*Token
		Err string
	}{Token: tk.Token, Err: tk.Err.Error()})
}

func saveJson(filename string, tks []token.Token) error {
	if b, err := encodingJson(tks); err != nil {
		return err
	} else {
		return ioutil.WriteFile(filename, b, os.ModePerm)
	}
}

func exists(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func Test_scanner_next(t *testing.T) {
	s := &scanner{}
	code := "a中♥\r\n\ra\\\nb\na\\\r\nb"
	s.new("", bytes.NewBufferString(code))
	tests := []struct {
		Lit rune
		Pos token.Position
	}{
		{'a', token.Position{Filename: "", Line: 1, Column: 1}},
		{'中', token.Position{Filename: "", Line: 1, Column: 2}},
		{'♥', token.Position{Filename: "", Line: 1, Column: 3}},
		{'\n', token.Position{Filename: "", Line: 1, Column: 4}},
		{'\n', token.Position{Filename: "", Line: 2, Column: 1}},
		{'a', token.Position{Filename: "", Line: 3, Column: 1}},
		{'b', token.Position{Filename: "", Line: 4, Column: 1}},
		{'\n', token.Position{Filename: "", Line: 4, Column: 2}},
		{'a', token.Position{Filename: "", Line: 5, Column: 1}},
		{'b', token.Position{Filename: "", Line: 6, Column: 1}},
	}

	for _, tt := range tests {
		t.Run(string(tt.Lit), func(t *testing.T) {
			ch, pos := s.next()
			if ch != tt.Lit {
				t.Errorf("want %v got %v", tt.Lit, ch)
			}
			if !reflect.DeepEqual(tt.Pos, pos) {
				t.Errorf("want %+v got %+v", tt.Pos, pos)
			}
		})
	}
}

func TestScanFile(t *testing.T) {
	if err := filepath.Walk("testdata/", func(p string, info os.FileInfo, err error) error {
		ext := filepath.Ext(p)
		if ext == ".c" {
			t.Run(p, func(t *testing.T) {
				tks, err := ScanFile(p)

				expect := p + ".expect.json"
				if !exists(expect) {
					if err := saveJson(expect, tks); err != nil {
						t.Errorf("saveJson error = %v", err)
						return
					}
				}

				want, err := ioutil.ReadFile(expect)
				if err != nil {
					t.Errorf("ReadFile() error = %v", err)
					return
				}

				got, _ := encodingJson(tks)
				if !bytes.Equal(want, got) {
					t.Errorf("ScanFile() got = %s, want %s", got, want)
				}
			})
		}
		return nil
	}); err != nil {
		t.Error(err)
	}
}

func Test_scanner_scanQuote(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		wantErr bool
	}{
		{
			"success", `'1' '\123' '\x12' '\u1234' '\U12345678'`, false,
		},
		{
			"error", `'12'`, true,
		},
		{
			"error", `'\u12'`, true,
		},
		{
			"error", `'\1234'`, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ScanString("", tt.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("ScanString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("got error %v", err)
		})
	}
}
