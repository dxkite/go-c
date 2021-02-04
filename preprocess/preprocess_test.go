package preprocess

import (
	"dxkite.cn/go-c11/scanner"
	"dxkite.cn/go-c11/token"
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
				ctx := New(scanner.NewScan(p, f))

				if !exists(p + ".json") {
					tks := []*token.Token{}
					for {
						t := ctx.Scan()
						if t == nil {
							break
						}
						tks = append(tks, t)
					}
					if err := token.SaveJson(p+".json", tks); err != nil {
						t.Errorf("SaveJson error = %v", err)
					}
					return
				}

				wantList, err := token.LoadJson(p + ".json")

				if err != nil {
					t.Errorf("LoadJson() error = %v", err)
					return
				}

				for _, want := range wantList {
					got := ctx.Scan()
					if !reflect.DeepEqual(got, want) {
						t.Errorf("ScanFile() got = %+v, want %+v", *got, *want)
					}
				}
			})
		}
		return nil
	}); err != nil {
		t.Error(err)
	}
}
