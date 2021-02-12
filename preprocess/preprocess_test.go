package preprocess

import (
	"bytes"
	go_c11 "dxkite.cn/c11"
	"dxkite.cn/c11/scanner"
	"dxkite.cn/c11/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

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

func TestScanFile(t *testing.T) {
	if err := filepath.Walk("testdata/", func(p string, info os.FileInfo, err error) error {
		ext := filepath.Ext(p)
		if ext == ".c" {

			t.Run(p, func(t *testing.T) {
				f, err := os.OpenFile(p, os.O_RDONLY, os.ModePerm)
				if err != nil {
					t.Errorf("read file error = %v", err)
					return
				}
				defer func() { _ = f.Close() }()
				ctx := NewContext()
				ctx.Init()

				exp := NewExpander(ctx, scanner.NewScan(p, f))

				tks := scanner.ScanToken(exp)

				expectToken := p + ".json"
				expectCode := p + "c"
				expectErr := p + ".err.json"

				if !exists(expectToken) || !exists(expectCode) || !exists(expectErr) {

					if err := scanner.SaveJson(expectToken, tks); err != nil {
						t.Errorf("SaveJson error = %v", err)
					}

					if err := ioutil.WriteFile(expectCode, []byte(token.String(tks)), os.ModePerm); err != nil {
						t.Errorf("SaveResult error = %v", err)
					}

					if err := exp.Error().SaveFile(expectErr); err != nil {
						t.Errorf("SaveError error = %v", err)
					}
					return
				}

				if data, err := ioutil.ReadFile(expectCode); err != nil {
					t.Errorf("LoadResult error = %v", err)
				} else {
					got := token.String(tks)
					if !bytes.Equal(data, []byte(got)) {
						t.Errorf("result error:want:\t%s\ngot:\t%s\n", string(data), got)
					}
				}

				var loadError c11.ErrorList
				if err := loadError.LoadFromFile(expectErr); err != nil {
					t.Errorf("LoadFromFile error = %v", err)
				} else if !reflect.DeepEqual(loadError, *exp.Error()) {
					t.Errorf("ScanFile() got = %v, want %v", loadError, exp.Error())
				}
			})
		}
		return nil
	}); err != nil {
		t.Error(err)
	}
}
