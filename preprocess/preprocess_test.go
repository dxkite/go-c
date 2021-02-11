package preprocess

import (
	"bytes"
	"dxkite.cn/go-c11/scanner"
	"dxkite.cn/go-c11/token"
	"io/ioutil"
	"os"
	"path/filepath"
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

				if !exists(p+".json") || !exists(p+"c") {

					if err := scanner.SaveJson(p+".json", tks); err != nil {
						t.Errorf("SaveJson error = %v", err)
					}

					if err := ioutil.WriteFile(p+"c", []byte(token.String(tks)), os.ModePerm); err != nil {
						t.Errorf("SaveResult error = %v", err)
					}
					return
				}

				if data, err := ioutil.ReadFile(p + "c"); err != nil {
					t.Errorf("LoadResult error = %v", err)
				} else {
					got := token.String(tks)
					if !bytes.Equal(data, []byte(got)) {
						t.Errorf("result error:want:\t%s\ngot:\t%s\n", string(data), got)
					}
				}
			})
		}
		return nil
	}); err != nil {
		t.Error(err)
	}
}
