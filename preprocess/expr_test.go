package preprocess

import (
	"bytes"
	"dxkite.cn/c/scanner"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func encodingExprJson(expr Expr) ([]byte, error) {
	buf := &bytes.Buffer{}
	je := json.NewEncoder(buf)
	je.SetEscapeHTML(false)
	je.SetIndent("", "  ")
	if err := je.Encode(expr); err != nil {
		return nil, err
	} else {
		return buf.Bytes(), nil
	}
}

func saveExprJson(filename string, expr Expr) error {
	if b, err := encodingExprJson(expr); err != nil {
		return err
	} else {
		return ioutil.WriteFile(filename, b, os.ModePerm)
	}
}

func TestParser_ParseExpr(t *testing.T) {
	if err := filepath.Walk("testdata/expr", func(p string, info os.FileInfo, err error) error {
		ext := filepath.Ext(p)
		if ext == ".c" {
			t.Run(p, func(t *testing.T) {
				timeOut(t, func(t *testing.T) {

					f, err := os.OpenFile(p, os.O_RDONLY, os.ModePerm)
					if err != nil {
						t.Errorf("read file error = %v", err)
						return
					}
					defer func() { _ = f.Close() }()

					r := scanner.NewScan(p, f)
					parser := NewParser(r)
					expr := parser.ParseExpr()

					expectExpr := p + ".json"
					expectErr := p + ".err.json"
					expectCode := p + "c"

					if !exists(expectExpr) || !exists(expectCode) || !exists(expectErr) {
						if err := saveExprJson(expectExpr, expr); err != nil {
							t.Errorf("SaveJson error = %v", err)
						}
						if err := parser.Error().SaveFile(expectErr); err != nil {
							t.Errorf("SaveError error = %v", err)
						}
						if err := ioutil.WriteFile(expectCode, []byte(ExprString(expr)), os.ModePerm); err != nil {
							t.Errorf("SaveResult error = %v", err)
						}
						return
					}

					if data, err := ioutil.ReadFile(expectCode); err != nil {
						t.Errorf("LoadResult error = %v", err)
					} else {
						got := ExprString(expr)
						if !bytes.Equal(data, []byte(got)) {
							t.Errorf("result error:want:\t%s\ngot:\t%s\n", string(data), got)
						}
					}

					if data, err := ioutil.ReadFile(expectExpr); err != nil {
						t.Errorf("LoadResult error = %v", err)
					} else {
						got, err := encodingExprJson(expr)
						if err != nil {
							t.Errorf("LoadResult error = %v", err)
							return
						}
						if !bytes.Equal(data, got) {
							t.Errorf("result error:want:\t%s\ngot:\t%s\n", string(data), got)
						}
					}

					var loadError ErrorList
					if err := loadError.LoadFromFile(expectErr); err != nil {
						t.Errorf("LoadFromFile error = %v", err)
					} else if !reflect.DeepEqual(loadError, parser.Error()) {
						t.Errorf("ScanFile() want = %v, got %v", loadError, parser.Error())
					}
				})
			})
		}
		return nil
	}); err != nil {
		t.Error(err)
	}
}
