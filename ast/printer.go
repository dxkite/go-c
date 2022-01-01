package ast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

func Json(x interface{}, minify bool) ([]byte, error) {
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

func String(x interface{}, prefix, ident string) string {
	buf := &bytes.Buffer{}
	printString(buf, x, prefix, "", ident)
	return buf.String()
}

func printString(w io.Writer, x interface{}, prefix, name, ident string) {
	t := reflect.TypeOf(x)
	rv := reflect.ValueOf(x)

	if x == nil || (t.Kind() == reflect.Ptr && rv.IsNil()) {
		_, _ = fmt.Fprint(w, prefix+name+"<nil>\n")
		return
	}

	if v, ok := x.(fmt.Stringer); ok {
		_, _ = fmt.Fprint(w, prefix+name+v.String()+"\n")
		return
	}

	if t.Kind() == reflect.Ptr {
		printString(w, rv.Elem().Interface(), prefix, name, ident)
		return
	}

	if t.Kind() == reflect.Struct {
		_, _ = fmt.Fprint(w, prefix+name+t.Name()+"\n")
		n := t.NumField()
		for i := 0; i < n; i++ {
			f := t.Field(i)
			if i+1 == n {
				printString(w, rv.Field(i).Interface(), prefix+ident, "`+"+f.Name+" = ", ident)
			} else {
				printString(w, rv.Field(i).Interface(), prefix+ident+"|", "+"+f.Name+" = ", ident)
			}
		}
		return
	}

	if t.Kind() == reflect.Array || t.Kind() == reflect.Slice {
		n := rv.Len()
		_, _ = fmt.Fprint(w, prefix+name+t.Name()+"\n")
		for i := 0; i < n; i++ {
			if i+1 == n {
				printString(w, rv.Index(i).Interface(), prefix+ident, "`-", ident)
			} else {
				printString(w, rv.Index(i).Interface(), prefix+ident+"|", "-", ident)
			}
		}
		return
	}

	_, _ = fmt.Fprintf(w, prefix+name+"%v\n", x)
}
