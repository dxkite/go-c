package preprocess

import (
	"bytes"
	"dxkite.cn/c/scanner"
	"dxkite.cn/c/token"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
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

func saveJson(filename string, tks []token.Token) error {
	if b, err := encodingJson(tks); err != nil {
		return err
	} else {
		return ioutil.WriteFile(filename, b, os.ModePerm)
	}
}

func timeOut(t *testing.T, f func(t *testing.T)) {
	c := make(chan struct{})
	go func() {
		f(t)
		c <- struct{}{}
	}()
	select {
	case <-c:
	case <-time.After(time.Second * 3):
		t.Errorf("timeout")
	}
}

func TestScanFile(t *testing.T) {
	if err := filepath.Walk("testdata/macro", func(p string, info os.FileInfo, err error) error {
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
					ctx := NewContext()
					ctx.Init()

					exp := newProcessor(ctx, scanner.NewScan(p, f), true)

					tks, _ := scanner.ScanToken(exp)

					expectToken := p + ".json"
					expectCode := p + ".txt"
					expectErr := p + ".err.json"

					if !exists(expectToken) || !exists(expectCode) || !exists(expectErr) {

						if err := saveJson(expectToken, tks); err != nil {
							t.Errorf("SaveJson error = %v", err)
						}

						if err := ioutil.WriteFile(expectCode, []byte(tokenString(tks)), os.ModePerm); err != nil {
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
						got := tokenString(tks)
						if !bytes.Equal(data, []byte(got)) {
							t.Errorf("result error:want:\t%s\ngot:\t%s\n", string(data), got)
						}
					}

					var loadError ErrorList
					if err := loadError.LoadFromFile(expectErr); err != nil {
						t.Errorf("LoadFromFile error = %v", err)
					} else if !reflect.DeepEqual(loadError, exp.Error()) {
						t.Errorf("ScanFile() want = %v, got %v", loadError, exp.Error())
					}
				})
			})
		}
		return nil
	}); err != nil {
		t.Error(err)
	}
}

func Test_columnDelta(t *testing.T) {
	tks, _ := scanner.ScanString("", "abc + 1234")
	columnDelta(tks[2:], 10)
	if relativeTokenString(tks) != "abc           + 1234" {
		t.Error("delta error")
	}
}
