package scanner

import (
	"bytes"
	go_c11 "dxkite.cn/c11"
	"dxkite.cn/c11/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func Test_scanner_next(t *testing.T) {
	s := &scanner{}
	code := "a中♥\r\n\ra\\\nb\na\\\r\nb"
	s.init("", bytes.NewBufferString(code))
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

				expectErr := p + ".err.json"

				if err != nil {
					if errorList, ok := err.(*c11.ErrorList); ok {
						if !c11.Exists(expectErr) {
							if err := errorList.SaveFile(expectErr); err != nil {
								t.Errorf("SaveFile error = %v", err)
								return
							}
						}
						var loadError c11.ErrorList
						if err := loadError.LoadFromFile(expectErr); err != nil {
							t.Errorf("LoadFromFile error = %v", err)
						} else if !reflect.DeepEqual(loadError, *errorList) {
							t.Errorf("ScanFile() got = %v, want %v", loadError, errorList)
						}
					} else {
						t.Errorf("ScanFile error = %v", err)
						return
					}
				}

				expect := p + ".expect.json"
				if !c11.Exists(expect) {
					if err := SaveJson(expect, tks); err != nil {
						t.Errorf("SaveJson error = %v", err)
						return
					}
				}

				want, err := ioutil.ReadFile(expect)
				if err != nil {
					t.Errorf("ReadFile() error = %v", err)
					return
				}

				got, _ := EncodingJson(tks)
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
		})
	}
}

func Test_peekScanner_Peek(t *testing.T) {
	s := NewStringScan("", "int float u'1' '\\123' '\\x12' '\\u1234' '\\U12345678'")
	p := NewPeekScan(s)
	pn := p.Peek(3)
	for i, item := range pn {
		r := p.Scan()
		if !reflect.DeepEqual(item, r) {
			t.Errorf("PeekScanError(%d): want %v got %v", i, item, r)
		}
	}
}

func Test_multiScanner_Scan(t *testing.T) {
	s1 := NewStringScan("s1", "int '")
	s2 := NewStringScan("s2", "float u'1' a")
	s3 := NewStringScan("s3", "abc=12")
	s := NewMultiScan(s2, s1)
	p := NewPeekScan(s)
	pn := p.Peek(4)

	for i := 0; i < 2; i++ {
		r := p.Scan()
		if !reflect.DeepEqual(pn[i], r) {
			t.Errorf("PeekScanError(%d): want %v got %v", i, pn[i], r)
		}
	}

	ss := NewMultiScan(p)
	ss.Push(s3)

	for i := 0; i < 3; i++ {
		// scan abc=12
		ss.Scan()
	}

	for i := 2; i < 4; i++ {
		r := ss.Scan()
		if !reflect.DeepEqual(pn[i], r) {
			t.Errorf("PeekScanError(%d): want %v got %v", i, pn[i], r)
		}
	}
}
