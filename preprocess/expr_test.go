package preprocess

import (
	"bytes"
	"dxkite.cn/c/errors"
	"dxkite.cn/c/scanner"
	"encoding/json"
	"fmt"
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

func jsonEqual(a, b []byte) bool {
	aj := map[string]interface{}{}
	bj := map[string]interface{}{}
	if err := json.Unmarshal(a, &aj); err != nil {
		fmt.Println(err)
	}
	if err := json.Unmarshal(b, &bj); err != nil {
		fmt.Println(err)
	}
	return reflect.DeepEqual(aj, bj)
}

func TestParser_ParseExpr(t *testing.T) {
	_ = os.MkdirAll(result+"/expr", os.ModePerm)
	if err := filepath.Walk(source+"/expr", func(p string, info os.FileInfo, err error) error {
		ext := filepath.Ext(p)
		if ext == ".c" {
			name := filepath.Base(p)
			t.Run(name, func(t *testing.T) {
				timeOut(t, func(t *testing.T) {

					f, err := os.OpenFile(p, os.O_RDONLY, os.ModePerm)
					if err != nil {
						t.Errorf("read file error = %v", err)
						return
					}
					defer func() { _ = f.Close() }()
					ctx := NewContext()
					r := scanner.NewScan(p, f, nil)
					parser := NewParser(ctx, r)
					expr := parser.ParseExpr()

					expectExpr := fmt.Sprintf("%s/expr/%s.json", result, name)
					expectCode := fmt.Sprintf("%s/expr/%s", result, name)
					expectErr := fmt.Sprintf("%s/expr/%s.err.json", result, name)

					if !exists(expectExpr) || !exists(expectCode) || !exists(expectErr) {
						if err := saveExprJson(expectExpr, expr); err != nil {
							t.Errorf("SaveJson error = %v", err)
						}
						if err := ctx.Error().SaveFile(expectErr); err != nil {
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
						if !jsonEqual(data, got) {
							t.Errorf("result error:want:\t%s\ngot:\t%s\n", string(data), got)
						}
					}

					var loadError errors.ErrorList
					if err := loadError.LoadFromFile(expectErr); err != nil {
						t.Errorf("LoadFromFile error = %v", err)
					} else if !reflect.DeepEqual(loadError, ctx.Error()) {
						t.Errorf("ScanFile() want = %v, got %v", loadError, ctx.Error())
					}
				})
			})
		}
		return nil
	}); err != nil {
		t.Error(err)
	}
}
