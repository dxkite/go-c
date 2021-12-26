package ast

import (
	"bytes"
	"encoding/json"
	"reflect"
)

func DumpJson(x interface{}, minify bool) ([]byte, error) {
	buf := &bytes.Buffer{}
	je := json.NewEncoder(buf)
	if !minify {
		je.SetEscapeHTML(false)
		je.SetIndent("", "  ")
	}
	if err := je.Encode(dumpInterface(x)); err != nil {
		return nil, err
	} else {
		return buf.Bytes(), nil
	}
}

func dumpInterface(x interface{}) interface{} {
	t := reflect.TypeOf(x)
	rv := reflect.ValueOf(x)
	if t == nil {
		return nil
	}
	if t.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}
		return dumpInterface(rv.Elem().Interface())
	}
	if t.Kind() == reflect.Struct {
		v := map[string]interface{}{}
		v["@type"] = t.Name()
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			v[f.Name] = dumpInterface(rv.Field(i).Interface())
		}
		return v
	}
	if t.Kind() == reflect.Array || t.Kind() == reflect.Slice {
		n := rv.Len()
		v := make([]interface{}, n)
		for i := 0; i < n; i++ {
			v[i] = dumpInterface(rv.Index(i).Interface())
		}
		return v
	}
	return x
}
