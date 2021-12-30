package ast

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
)

func DumpString(x interface{}, prefix, ident string) string {
	buf := &bytes.Buffer{}
	printString(buf, x, prefix, "", ident)
	return buf.String()
}

func printString(w io.Writer, x interface{}, prefix, name, ident string) {
	t := reflect.TypeOf(x)
	rv := reflect.ValueOf(x)

	if t == nil {
		_, _ = fmt.Fprint(w, prefix+name+"<nil>\n")
		return
	}

	if v, ok := x.(fmt.Stringer); ok {
		_, _ = fmt.Fprint(w, prefix+name+v.String()+"\n")
		return
	}

	if t.Kind() == reflect.Ptr {
		if rv.IsNil() {
			_, _ = fmt.Fprint(w, prefix+name+"<nil>\n")
			return
		}
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
