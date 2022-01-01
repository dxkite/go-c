package parser

import (
	"dxkite.cn/c/ast"
	"dxkite.cn/c/errors"
	"dxkite.cn/c/preprocess"
	"dxkite.cn/c/scanner"
	"dxkite.cn/c/token"
	stderr "errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type Error struct {
	Pos token.Position
	Typ errors.ErrorType
	Msg string
}

//func testParseType(name, code, json, err string) bool {
//	errList := []Error{}
//	errHandler := errors.ErrorHandler(func(pos token.Position, typ errors.ErrorType, code errors.ErrCode, params ...interface{}) {
//		errList = append(errList, Error{
//			Pos: pos,
//			Typ: typ,
//			Msg: errors.New(pos, code, params...).Error(),
//		})
//	})
//
//	parser := newParser(scanner.NewStringScan("", code, nil), errHandler)
//	t := parser.parseTypeName()
//	jsonCode, _ := ast.Json(t, true)
//	jsonErr, _ := ast.Json(errList, true)
//	if string(jsonCode) != json || string(jsonErr) != err {
//		code := fmt.Sprintf("{%s,%s,%s,%s},", strconv.QuoteToGraphic(name),
//			strconv.QuoteToGraphic(code),
//			strconv.QuoteToGraphic(string(jsonCode)),
//			strconv.QuoteToGraphic(string(jsonErr)))
//		fmt.Println("code:\n", code)
//		fmt.Println("want:", json)
//		fmt.Println("got", string(jsonCode))
//		return false
//	}
//	return true
//}
//
//func TestParser_parseTypeName(t *testing.T) {
//	testCase := [][4]string{
//		{"simple-type", "int", "{\"@type\":\"BuildInType\",\"Lit\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":1,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}]}\n", "[]\n"},
//		{"incomplete-arr-parse", "const int int long *const[]", "{\"@type\":\"IncompleteArrayType\",\"Return\":{\"@type\":\"PointerType\",\"Return\":{\"@type\":\"TypeQualifier\",\"Return\":{\"@type\":\"BuildInType\",\"Lit\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":7,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"},{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":11,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"},{\"@type\":\"Token\",\"Lit\":\"long\",\"Pos\":{\"@type\":\"Position\",\"Column\":15,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}]},\"Qualifier\":{\"const\":true}},\"Qualifier\":{\"const\":true}}}\n", "[]\n"},
//		{"const-arr", "const int int long *const[10?2:1]", "{\"@type\":\"ConstArrayType\",\"Return\":{\"@type\":\"PointerType\",\"Return\":{\"@type\":\"TypeQualifier\",\"Return\":{\"@type\":\"BuildInType\",\"Lit\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":7,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"},{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":11,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"},{\"@type\":\"Token\",\"Lit\":\"long\",\"Pos\":{\"@type\":\"Position\",\"Column\":15,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}]},\"Qualifier\":{\"const\":true}},\"Qualifier\":{\"const\":true}},\"Size\":{\"@type\":\"CondExpr\",\"Else\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"1\",\"Pos\":{\"@type\":\"Position\",\"Column\":32,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}},\"Op\":{\"@type\":\"Token\",\"Lit\":\"?\",\"Pos\":{\"@type\":\"Position\",\"Column\":29,\"Filename\":\"\",\"Line\":1},\"Decl\":\"PUNCTUATOR\"},\"Then\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"2\",\"Pos\":{\"@type\":\"Position\",\"Column\":30,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}},\"X\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"10\",\"Pos\":{\"@type\":\"Position\",\"Column\":27,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}}}}\n", "[]\n"},
//		{"arr-parse", "const int int long *const[m=10]", "{\"@type\":\"ArrayType\",\"Return\":{\"@type\":\"PointerType\",\"Return\":{\"@type\":\"TypeQualifier\",\"Return\":{\"@type\":\"BuildInType\",\"Lit\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":7,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"},{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":11,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"},{\"@type\":\"Token\",\"Lit\":\"long\",\"Pos\":{\"@type\":\"Position\",\"Column\":15,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}]},\"Qualifier\":{\"const\":true}},\"Qualifier\":{\"const\":true}},\"Qualifier\":{},\"Size\":{\"@type\":\"AssignExpr\",\"Op\":{\"@type\":\"Token\",\"Lit\":\"=\",\"Pos\":{\"@type\":\"Position\",\"Column\":28,\"Filename\":\"\",\"Line\":1},\"Decl\":\"PUNCTUATOR\"},\"X\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"m\",\"Pos\":{\"@type\":\"Position\",\"Column\":27,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Y\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"10\",\"Pos\":{\"@type\":\"Position\",\"Column\":29,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}}},\"Static\":false}\n", "[]\n"},
//		{"func-ptr", "int(*)(int,int)", "{\"@type\":\"FuncType\",\"Ellipsis\":false,\"Return\":{\"@type\":\"PointerType\",\"Return\":{\"@type\":\"ParenType\",\"Return\":{\"@type\":\"BuildInType\",\"Lit\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":1,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}]}},\"Qualifier\":{}},\"Params\":[{\"@type\":\"ParamVarDecl\",\"Name\":null,\"Lit\":{\"@type\":\"BuildInType\",\"Lit\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":8,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}]}},{\"@type\":\"ParamVarDecl\",\"Name\":null,\"Lit\":{\"@type\":\"BuildInType\",\"Lit\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":12,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}]}}]}\n", "[]\n"},
//		{"func-with-type", "int(*)(int a,int b)", "{\"@type\":\"FuncType\",\"Ellipsis\":false,\"Return\":{\"@type\":\"PointerType\",\"Return\":{\"@type\":\"ParenType\",\"Return\":{\"@type\":\"BuildInType\",\"Lit\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":1,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}]}},\"Qualifier\":{}},\"Params\":[{\"@type\":\"ParamVarDecl\",\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"a\",\"Pos\":{\"@type\":\"Position\",\"Column\":12,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Lit\":{\"@type\":\"BuildInType\",\"Lit\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":8,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}]}},{\"@type\":\"ParamVarDecl\",\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"b\",\"Pos\":{\"@type\":\"Position\",\"Column\":18,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Lit\":{\"@type\":\"BuildInType\",\"Lit\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":14,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}]}}]}\n", "[]\n"},
//		{"func-arr-const", "int(*)(int***const a,const int b[a][m])", "{\"@type\":\"FuncType\",\"Ellipsis\":false,\"Return\":{\"@type\":\"PointerType\",\"Return\":{\"@type\":\"ParenType\",\"Return\":{\"@type\":\"BuildInType\",\"Lit\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":1,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}]}},\"Qualifier\":{}},\"Params\":[{\"@type\":\"ParamVarDecl\",\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"a\",\"Pos\":{\"@type\":\"Position\",\"Column\":20,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Lit\":{\"@type\":\"PointerType\",\"Return\":{\"@type\":\"PointerType\",\"Return\":{\"@type\":\"PointerType\",\"Return\":{\"@type\":\"BuildInType\",\"Lit\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":8,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}]},\"Qualifier\":{}},\"Qualifier\":{}},\"Qualifier\":{\"const\":true}}},{\"@type\":\"ParamVarDecl\",\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"b\",\"Pos\":{\"@type\":\"Position\",\"Column\":32,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Lit\":{\"@type\":\"ConstArrayType\",\"Return\":{\"@type\":\"ConstArrayType\",\"Return\":{\"@type\":\"TypeQualifier\",\"Return\":{\"@type\":\"BuildInType\",\"Lit\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":28,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}]},\"Qualifier\":{\"const\":true}},\"Size\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"a\",\"Pos\":{\"@type\":\"Position\",\"Column\":34,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}}},\"Size\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"m\",\"Pos\":{\"@type\":\"Position\",\"Column\":37,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}}}}]}\n", "[]\n"},
//		{"struct-defined-noname-bit-field", "volatile struct {int a:1?2:3,:9; char c;};", "{\"@type\":\"TypeQualifier\",\"Return\":{\"@type\":\"RecordType\",\"Fields\":[{\"@type\":\"RecordField\",\"Bit\":{\"@type\":\"ConstantExpr\",\"X\":{\"@type\":\"CondExpr\",\"Else\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"3\",\"Pos\":{\"@type\":\"Position\",\"Column\":28,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}},\"Op\":{\"@type\":\"Token\",\"Lit\":\"?\",\"Pos\":{\"@type\":\"Position\",\"Column\":25,\"Filename\":\"\",\"Line\":1},\"Decl\":\"PUNCTUATOR\"},\"Then\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"2\",\"Pos\":{\"@type\":\"Position\",\"Column\":26,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}},\"X\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"1\",\"Pos\":{\"@type\":\"Position\",\"Column\":24,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}}}},\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"a\",\"Pos\":{\"@type\":\"Position\",\"Column\":22,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Lit\":{\"@type\":\"BuildInType\",\"Lit\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":18,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}]}},{\"@type\":\"RecordField\",\"Bit\":{\"@type\":\"ConstantExpr\",\"X\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"9\",\"Pos\":{\"@type\":\"Position\",\"Column\":31,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}}},\"Name\":null,\"Lit\":{\"@type\":\"BuildInType\",\"Lit\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":18,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}]}},{\"@type\":\"RecordField\",\"Bit\":null,\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"c\",\"Pos\":{\"@type\":\"Position\",\"Column\":39,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Lit\":{\"@type\":\"BuildInType\",\"Lit\":[{\"@type\":\"Token\",\"Lit\":\"char\",\"Pos\":{\"@type\":\"Position\",\"Column\":34,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}]}}],\"Name\":null,\"Lit\":{\"@type\":\"Token\",\"Lit\":\"struct\",\"Pos\":{\"@type\":\"Position\",\"Column\":10,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}},\"Qualifier\":{\"volatile\":true}}\n", "[]\n"},
//		{"struct-nested-union", "struct {union{ int a; int b;} const; int c; };", "{\"@type\":\"RecordType\",\"Fields\":[{\"@type\":\"RecordField\",\"Bit\":null,\"Name\":null,\"Lit\":{\"@type\":\"TypeQualifier\",\"Return\":{\"@type\":\"RecordType\",\"Fields\":[{\"@type\":\"RecordField\",\"Bit\":null,\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"a\",\"Pos\":{\"@type\":\"Position\",\"Column\":20,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Lit\":{\"@type\":\"BuildInType\",\"Lit\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":16,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}]}},{\"@type\":\"RecordField\",\"Bit\":null,\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"b\",\"Pos\":{\"@type\":\"Position\",\"Column\":27,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Lit\":{\"@type\":\"BuildInType\",\"Lit\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":23,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}]}}],\"Name\":null,\"Lit\":{\"@type\":\"Token\",\"Lit\":\"union\",\"Pos\":{\"@type\":\"Position\",\"Column\":9,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}},\"Qualifier\":{\"const\":true}}},{\"@type\":\"RecordField\",\"Bit\":null,\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"c\",\"Pos\":{\"@type\":\"Position\",\"Column\":42,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Lit\":{\"@type\":\"BuildInType\",\"Lit\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":38,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}]}}],\"Name\":null,\"Lit\":{\"@type\":\"Token\",\"Lit\":\"struct\",\"Pos\":{\"@type\":\"Position\",\"Column\":1,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}}\n", "[]\n"},
//		//添加报错
//		{"struct-nested-union", "struct test {struct test const; int c; };", "{\"@type\":\"RecordType\",\"Fields\":[{\"@type\":\"RecordField\",\"Bit\":null,\"Name\":null,\"Lit\":{\"@type\":\"TypeQualifier\",\"Return\":{\"@type\":\"RecordType\",\"Fields\":[],\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"test\",\"Pos\":{\"@type\":\"Position\",\"Column\":21,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Lit\":{\"@type\":\"Token\",\"Lit\":\"struct\",\"Pos\":{\"@type\":\"Position\",\"Column\":14,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}},\"Qualifier\":{\"const\":true}}},{\"@type\":\"RecordField\",\"Bit\":null,\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"c\",\"Pos\":{\"@type\":\"Position\",\"Column\":37,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Lit\":{\"@type\":\"BuildInType\",\"Lit\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":33,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}]}}],\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"test\",\"Pos\":{\"@type\":\"Position\",\"Column\":8,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Lit\":{\"@type\":\"Token\",\"Lit\":\"struct\",\"Pos\":{\"@type\":\"Position\",\"Column\":1,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}}\n", "[]\n"},
//		{"struct-nested-union-pointer", "struct test {struct test const * left; int c; };", "{\"@type\":\"RecordType\",\"Fields\":[{\"@type\":\"RecordField\",\"Bit\":null,\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"left\",\"Pos\":{\"@type\":\"Position\",\"Column\":34,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Lit\":{\"@type\":\"PointerType\",\"Return\":{\"@type\":\"TypeQualifier\",\"Return\":{\"@type\":\"RecordType\",\"Fields\":[],\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"test\",\"Pos\":{\"@type\":\"Position\",\"Column\":21,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Lit\":{\"@type\":\"Token\",\"Lit\":\"struct\",\"Pos\":{\"@type\":\"Position\",\"Column\":14,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}},\"Qualifier\":{\"const\":true}},\"Qualifier\":{}}},{\"@type\":\"RecordField\",\"Bit\":null,\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"c\",\"Pos\":{\"@type\":\"Position\",\"Column\":44,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Lit\":{\"@type\":\"BuildInType\",\"Lit\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":40,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}]}}],\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"test\",\"Pos\":{\"@type\":\"Position\",\"Column\":8,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Lit\":{\"@type\":\"Token\",\"Lit\":\"struct\",\"Pos\":{\"@type\":\"Position\",\"Column\":1,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}}\n", "[]\n"},
//		{"func", "int(a,b)", "{\"@type\":\"FuncType\",\"Ellipsis\":false,\"Return\":{\"@type\":\"BuildInType\",\"Lit\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":1,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}]},\"Params\":[{\"@type\":\"ParamVarDecl\",\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"a\",\"Pos\":{\"@type\":\"Position\",\"Column\":5,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Lit\":null},{\"@type\":\"ParamVarDecl\",\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"b\",\"Pos\":{\"@type\":\"Position\",\"Column\":7,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Lit\":null}]}\n", "[]\n"},
//		{"typedef", "struct test {struct test const * left; int c; }", "{\"@type\":\"RecordType\",\"Fields\":[{\"@type\":\"RecordField\",\"Bit\":null,\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"left\",\"Pos\":{\"@type\":\"Position\",\"Column\":34,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Lit\":{\"@type\":\"PointerType\",\"Return\":{\"@type\":\"TypeQualifier\",\"Return\":{\"@type\":\"RecordType\",\"Fields\":[],\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"test\",\"Pos\":{\"@type\":\"Position\",\"Column\":21,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Lit\":{\"@type\":\"Token\",\"Lit\":\"struct\",\"Pos\":{\"@type\":\"Position\",\"Column\":14,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}},\"Qualifier\":{\"const\":true}},\"Qualifier\":{}}},{\"@type\":\"RecordField\",\"Bit\":null,\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"c\",\"Pos\":{\"@type\":\"Position\",\"Column\":44,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Lit\":{\"@type\":\"BuildInType\",\"Lit\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":40,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}]}}],\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"test\",\"Pos\":{\"@type\":\"Position\",\"Column\":8,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Lit\":{\"@type\":\"Token\",\"Lit\":\"struct\",\"Pos\":{\"@type\":\"Position\",\"Column\":1,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}}\n", "[]\n"},
//	}
//
//	for _, cas := range testCase {
//		t.Run(cas[0], func(t *testing.T) {
//			if !testParseType(cas[0], cas[1], cas[2], cas[3]) {
//				t.Error("test error")
//			}
//		})
//	}
//}
//
//func testParseExpr(name, code, json, err string) bool {
//	errList := []Error{}
//	errHandler := errors.ErrorHandler(func(pos token.Position, typ errors.ErrorType, code errors.ErrCode, params ...interface{}) {
//		errList = append(errList, Error{
//			Pos: pos,
//			Typ: typ,
//			Msg: errors.New(pos, code, params...).Error(),
//		})
//	})
//	parser := newParser(scanner.NewStringScan("", code, nil), errHandler)
//	t := parser.parseExpr()
//	jsonCode, _ := ast.Json(t, true)
//	jsonErr, _ := ast.Json(errList, true)
//	if string(jsonCode) != json || string(jsonErr) != err {
//		code := fmt.Sprintf("{%s,%s,%s,%s},", strconv.QuoteToGraphic(name),
//			strconv.QuoteToGraphic(code),
//			strconv.QuoteToGraphic(string(jsonCode)),
//			strconv.QuoteToGraphic(string(jsonErr)))
//		fmt.Println("code:\n", code)
//		fmt.Println("want:", json)
//		fmt.Println("got", string(jsonCode))
//		return false
//	}
//	return true
//}
//
//func TestParser_parseAssignExpr(t *testing.T) {
//	testCase := [][4]string{
//		{"initializer-expr", "answer = (struct point){ .quot = 2, .rem = -1 }", "{\"@type\":\"AssignExpr\",\"Op\":{\"@type\":\"Token\",\"Lit\":\"=\",\"Pos\":{\"@type\":\"Position\",\"Column\":8,\"Filename\":\"\",\"Line\":1},\"Decl\":\"PUNCTUATOR\"},\"X\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"answer\",\"Pos\":{\"@type\":\"Position\",\"Column\":1,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Y\":{\"@type\":\"CompoundLit\",\"InitList\":[{\"@type\":\"RecordDesignatorExpr\",\"Field\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"quot\",\"Pos\":{\"@type\":\"Position\",\"Column\":27,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"X\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"2\",\"Pos\":{\"@type\":\"Position\",\"Column\":34,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}}},{\"@type\":\"RecordDesignatorExpr\",\"Field\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"rem\",\"Pos\":{\"@type\":\"Position\",\"Column\":38,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"X\":{\"@type\":\"UnaryExpr\",\"Op\":{\"@type\":\"Token\",\"Lit\":\"-\",\"Pos\":{\"@type\":\"Position\",\"Column\":44,\"Filename\":\"\",\"Line\":1},\"Decl\":\"PUNCTUATOR\"},\"X\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"1\",\"Pos\":{\"@type\":\"Position\",\"Column\":45,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}}}}],\"Lit\":{\"@type\":\"RecordType\",\"Fields\":[],\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"point\",\"Pos\":{\"@type\":\"Position\",\"Column\":18,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Lit\":{\"@type\":\"Token\",\"Lit\":\"struct\",\"Pos\":{\"@type\":\"Position\",\"Column\":11,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}}}}\n", "[]\n"},
//		{"initializer-expr-comma", "answer = (struct point){ .quot = 2, .rem = -1, }", "{\"@type\":\"AssignExpr\",\"Op\":{\"@type\":\"Token\",\"Lit\":\"=\",\"Pos\":{\"@type\":\"Position\",\"Column\":8,\"Filename\":\"\",\"Line\":1},\"Decl\":\"PUNCTUATOR\"},\"X\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"answer\",\"Pos\":{\"@type\":\"Position\",\"Column\":1,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Y\":{\"@type\":\"CompoundLit\",\"InitList\":[{\"@type\":\"RecordDesignatorExpr\",\"Field\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"quot\",\"Pos\":{\"@type\":\"Position\",\"Column\":27,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"X\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"2\",\"Pos\":{\"@type\":\"Position\",\"Column\":34,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}}},{\"@type\":\"RecordDesignatorExpr\",\"Field\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"rem\",\"Pos\":{\"@type\":\"Position\",\"Column\":38,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"X\":{\"@type\":\"UnaryExpr\",\"Op\":{\"@type\":\"Token\",\"Lit\":\"-\",\"Pos\":{\"@type\":\"Position\",\"Column\":44,\"Filename\":\"\",\"Line\":1},\"Decl\":\"PUNCTUATOR\"},\"X\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"1\",\"Pos\":{\"@type\":\"Position\",\"Column\":45,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}}}}],\"Lit\":{\"@type\":\"RecordType\",\"Fields\":[],\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"point\",\"Pos\":{\"@type\":\"Position\",\"Column\":18,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Lit\":{\"@type\":\"Token\",\"Lit\":\"struct\",\"Pos\":{\"@type\":\"Position\",\"Column\":11,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}}}}\n", "[]\n"},
//		{"simple-expr", "a+b+c*d/e", "{\"@type\":\"BinaryExpr\",\"Op\":{\"@type\":\"Token\",\"Lit\":\"+\",\"Pos\":{\"@type\":\"Position\",\"Column\":4,\"Filename\":\"\",\"Line\":1},\"Decl\":\"PUNCTUATOR\"},\"X\":{\"@type\":\"BinaryExpr\",\"Op\":{\"@type\":\"Token\",\"Lit\":\"+\",\"Pos\":{\"@type\":\"Position\",\"Column\":2,\"Filename\":\"\",\"Line\":1},\"Decl\":\"PUNCTUATOR\"},\"X\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"a\",\"Pos\":{\"@type\":\"Position\",\"Column\":1,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Y\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"b\",\"Pos\":{\"@type\":\"Position\",\"Column\":3,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}}},\"Y\":{\"@type\":\"BinaryExpr\",\"Op\":{\"@type\":\"Token\",\"Lit\":\"/\",\"Pos\":{\"@type\":\"Position\",\"Column\":8,\"Filename\":\"\",\"Line\":1},\"Decl\":\"PUNCTUATOR\"},\"X\":{\"@type\":\"BinaryExpr\",\"Op\":{\"@type\":\"Token\",\"Lit\":\"*\",\"Pos\":{\"@type\":\"Position\",\"Column\":6,\"Filename\":\"\",\"Line\":1},\"Decl\":\"PUNCTUATOR\"},\"X\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"c\",\"Pos\":{\"@type\":\"Position\",\"Column\":5,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Y\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"d\",\"Pos\":{\"@type\":\"Position\",\"Column\":7,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}}},\"Y\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"e\",\"Pos\":{\"@type\":\"Position\",\"Column\":9,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}}}}\n", "[]\n"},
//		{"simple-expr-postfix", "a+b+c*d/e++", "{\"@type\":\"BinaryExpr\",\"Op\":{\"@type\":\"Token\",\"Lit\":\"+\",\"Pos\":{\"@type\":\"Position\",\"Column\":4,\"Filename\":\"\",\"Line\":1},\"Decl\":\"PUNCTUATOR\"},\"X\":{\"@type\":\"BinaryExpr\",\"Op\":{\"@type\":\"Token\",\"Lit\":\"+\",\"Pos\":{\"@type\":\"Position\",\"Column\":2,\"Filename\":\"\",\"Line\":1},\"Decl\":\"PUNCTUATOR\"},\"X\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"a\",\"Pos\":{\"@type\":\"Position\",\"Column\":1,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Y\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"b\",\"Pos\":{\"@type\":\"Position\",\"Column\":3,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}}},\"Y\":{\"@type\":\"BinaryExpr\",\"Op\":{\"@type\":\"Token\",\"Lit\":\"/\",\"Pos\":{\"@type\":\"Position\",\"Column\":8,\"Filename\":\"\",\"Line\":1},\"Decl\":\"PUNCTUATOR\"},\"X\":{\"@type\":\"BinaryExpr\",\"Op\":{\"@type\":\"Token\",\"Lit\":\"*\",\"Pos\":{\"@type\":\"Position\",\"Column\":6,\"Filename\":\"\",\"Line\":1},\"Decl\":\"PUNCTUATOR\"},\"X\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"c\",\"Pos\":{\"@type\":\"Position\",\"Column\":5,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Y\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"d\",\"Pos\":{\"@type\":\"Position\",\"Column\":7,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}}},\"Y\":{\"@type\":\"UnaryExpr\",\"Op\":{\"@type\":\"Token\",\"Lit\":\"++\",\"Pos\":{\"@type\":\"Position\",\"Column\":10,\"Filename\":\"\",\"Line\":1},\"Decl\":\"PUNCTUATOR\"},\"X\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"e\",\"Pos\":{\"@type\":\"Position\",\"Column\":9,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}}}}}\n", "[]\n"},
//		{"condition-expr", "1?2:3?4:5", "{\"@type\":\"CondExpr\",\"Else\":{\"@type\":\"CondExpr\",\"Else\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"5\",\"Pos\":{\"@type\":\"Position\",\"Column\":9,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}},\"Op\":{\"@type\":\"Token\",\"Lit\":\"?\",\"Pos\":{\"@type\":\"Position\",\"Column\":6,\"Filename\":\"\",\"Line\":1},\"Decl\":\"PUNCTUATOR\"},\"Then\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"4\",\"Pos\":{\"@type\":\"Position\",\"Column\":7,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}},\"X\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"3\",\"Pos\":{\"@type\":\"Position\",\"Column\":5,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}}},\"Op\":{\"@type\":\"Token\",\"Lit\":\"?\",\"Pos\":{\"@type\":\"Position\",\"Column\":2,\"Filename\":\"\",\"Line\":1},\"Decl\":\"PUNCTUATOR\"},\"Then\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"2\",\"Pos\":{\"@type\":\"Position\",\"Column\":3,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}},\"X\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"1\",\"Pos\":{\"@type\":\"Position\",\"Column\":1,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}}}\n", "[]\n"},
//		{"arr-expr", "(1?2:3?4:5)[12][23]", "{\"@type\":\"IndexExpr\",\"Arr\":{\"@type\":\"IndexExpr\",\"Arr\":{\"@type\":\"ParenExpr\",\"X\":{\"@type\":\"CondExpr\",\"Else\":{\"@type\":\"CondExpr\",\"Else\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"5\",\"Pos\":{\"@type\":\"Position\",\"Column\":10,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}},\"Op\":{\"@type\":\"Token\",\"Lit\":\"?\",\"Pos\":{\"@type\":\"Position\",\"Column\":7,\"Filename\":\"\",\"Line\":1},\"Decl\":\"PUNCTUATOR\"},\"Then\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"4\",\"Pos\":{\"@type\":\"Position\",\"Column\":8,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}},\"X\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"3\",\"Pos\":{\"@type\":\"Position\",\"Column\":6,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}}},\"Op\":{\"@type\":\"Token\",\"Lit\":\"?\",\"Pos\":{\"@type\":\"Position\",\"Column\":3,\"Filename\":\"\",\"Line\":1},\"Decl\":\"PUNCTUATOR\"},\"Then\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"2\",\"Pos\":{\"@type\":\"Position\",\"Column\":4,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}},\"X\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"1\",\"Pos\":{\"@type\":\"Position\",\"Column\":2,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}}}},\"Index\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"12\",\"Pos\":{\"@type\":\"Position\",\"Column\":13,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}}},\"Index\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"23\",\"Pos\":{\"@type\":\"Position\",\"Column\":17,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}}}\n", "[]\n"},
//		{"call-expr", "(1?2:3?4:5)[12][23](1,2,3,4,b(a,c,v))", "{\"@type\":\"CallExpr\",\"Args\":[{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"1\",\"Pos\":{\"@type\":\"Position\",\"Column\":21,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}},{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"2\",\"Pos\":{\"@type\":\"Position\",\"Column\":23,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}},{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"3\",\"Pos\":{\"@type\":\"Position\",\"Column\":25,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}},{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"4\",\"Pos\":{\"@type\":\"Position\",\"Column\":27,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}},{\"@type\":\"CallExpr\",\"Args\":[{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"a\",\"Pos\":{\"@type\":\"Position\",\"Column\":31,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"c\",\"Pos\":{\"@type\":\"Position\",\"Column\":33,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"v\",\"Pos\":{\"@type\":\"Position\",\"Column\":35,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}}],\"Func\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"b\",\"Pos\":{\"@type\":\"Position\",\"Column\":29,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}}}],\"Func\":{\"@type\":\"IndexExpr\",\"Arr\":{\"@type\":\"IndexExpr\",\"Arr\":{\"@type\":\"ParenExpr\",\"X\":{\"@type\":\"CondExpr\",\"Else\":{\"@type\":\"CondExpr\",\"Else\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"5\",\"Pos\":{\"@type\":\"Position\",\"Column\":10,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}},\"Op\":{\"@type\":\"Token\",\"Lit\":\"?\",\"Pos\":{\"@type\":\"Position\",\"Column\":7,\"Filename\":\"\",\"Line\":1},\"Decl\":\"PUNCTUATOR\"},\"Then\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"4\",\"Pos\":{\"@type\":\"Position\",\"Column\":8,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}},\"X\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"3\",\"Pos\":{\"@type\":\"Position\",\"Column\":6,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}}},\"Op\":{\"@type\":\"Token\",\"Lit\":\"?\",\"Pos\":{\"@type\":\"Position\",\"Column\":3,\"Filename\":\"\",\"Line\":1},\"Decl\":\"PUNCTUATOR\"},\"Then\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"2\",\"Pos\":{\"@type\":\"Position\",\"Column\":4,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}},\"X\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"1\",\"Pos\":{\"@type\":\"Position\",\"Column\":2,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}}}},\"Index\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"12\",\"Pos\":{\"@type\":\"Position\",\"Column\":13,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}}},\"Index\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"23\",\"Pos\":{\"@type\":\"Position\",\"Column\":17,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}}}}\n", "[]\n"},
//	}
//
//	for _, cas := range testCase {
//		t.Run(cas[0], func(t *testing.T) {
//			if !testParseExpr(cas[0], cas[1], cas[2], cas[3]) {
//				t.Error("test error")
//			}
//		})
//	}
//}
//
//func testParseDecl(name, code, json, err string) bool {
//	errList := []Error{}
//	errHandler := errors.ErrorHandler(func(pos token.Position, typ errors.ErrorType, code errors.ErrCode, params ...interface{}) {
//		errList = append(errList, Error{
//			Pos: pos,
//			Typ: typ,
//			Msg: errors.New(pos, code, params...).Error(),
//		})
//	})
//	parser := newParser(scanner.NewStringScan("", code, nil), errHandler)
//	t := parser.parseUnit()
//	jsonCode, _ := ast.Json(t, true)
//	jsonErr, _ := ast.Json(errList, true)
//	if string(jsonCode) != json || string(jsonErr) != err {
//		code := fmt.Sprintf("{%s,%s,%s,%s},", strconv.QuoteToGraphic(name),
//			strconv.QuoteToGraphic(code),
//			strconv.QuoteToGraphic(string(jsonCode)),
//			strconv.QuoteToGraphic(string(jsonErr)))
//		fmt.Println("code:\n", code)
//		fmt.Println("want:", json)
//		fmt.Println("got", string(jsonCode))
//		return false
//	}
//	return true
//}
//func TestParser_ParseDecl(t *testing.T) {
//	testCase := [][4]string{
//		{"int-main", "int main(){}", "{\"@type\":\"TranslationUnitDecl\",\"Decl\":[{\"@type\":\"FuncDecl\",\"Body\":[],\"Decl\":[],\"Ellipsis\":false,\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"main\",\"Pos\":{\"@type\":\"Position\",\"Column\":5,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Params\":[],\"Return\":{\"@type\":\"BuildInType\",\"Lit\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":1,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}]}}]}\n", "[]\n"},
//		{"max", "int max(int a, int b) { if (a>b) return a; return b; } int main(){ return max(10,20); }", "{\"@type\":\"TranslationUnitDecl\",\"Decl\":[{\"@type\":\"FuncDecl\",\"Body\":[{\"@type\":\"IfStmt\",\"Else\":null,\"Then\":{\"@type\":\"ReturnStmt\",\"X\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"a\",\"Pos\":{\"@type\":\"Position\",\"Column\":41,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}}},\"X\":{\"@type\":\"BinaryExpr\",\"Op\":{\"@type\":\"Token\",\"Lit\":\"\\u003e\",\"Pos\":{\"@type\":\"Position\",\"Column\":30,\"Filename\":\"\",\"Line\":1},\"Decl\":\"PUNCTUATOR\"},\"X\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"a\",\"Pos\":{\"@type\":\"Position\",\"Column\":29,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Y\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"b\",\"Pos\":{\"@type\":\"Position\",\"Column\":31,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}}}},{\"@type\":\"ReturnStmt\",\"X\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"b\",\"Pos\":{\"@type\":\"Position\",\"Column\":51,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}}}],\"Decl\":[],\"Ellipsis\":false,\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"max\",\"Pos\":{\"@type\":\"Position\",\"Column\":5,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Params\":[{\"@type\":\"ParamVarDecl\",\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"a\",\"Pos\":{\"@type\":\"Position\",\"Column\":13,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Lit\":{\"@type\":\"BuildInType\",\"Lit\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":9,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}]}},{\"@type\":\"ParamVarDecl\",\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"b\",\"Pos\":{\"@type\":\"Position\",\"Column\":20,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Lit\":{\"@type\":\"BuildInType\",\"Lit\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":16,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}]}}],\"Return\":{\"@type\":\"BuildInType\",\"Lit\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":1,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}]}},{\"@type\":\"FuncDecl\",\"Body\":[{\"@type\":\"ReturnStmt\",\"X\":{\"@type\":\"CallExpr\",\"Args\":[{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"10\",\"Pos\":{\"@type\":\"Position\",\"Column\":79,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}},{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"20\",\"Pos\":{\"@type\":\"Position\",\"Column\":82,\"Filename\":\"\",\"Line\":1},\"Decl\":\"INT\"}}],\"Func\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"max\",\"Pos\":{\"@type\":\"Position\",\"Column\":75,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}}}}],\"Decl\":[],\"Ellipsis\":false,\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"main\",\"Pos\":{\"@type\":\"Position\",\"Column\":60,\"Filename\":\"\",\"Line\":1},\"Decl\":\"IDENT\"}},\"Params\":[],\"Return\":{\"@type\":\"BuildInType\",\"Lit\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":56,\"Filename\":\"\",\"Line\":1},\"Decl\":\"KEYWORD\"}]}}]}\n", "[]\n"},
//	}
//
//	for _, cas := range testCase {
//		t.Run(cas[0], func(t *testing.T) {
//			if !testParseDecl(cas[0], cas[1], cas[2], cas[3]) {
//				t.Error("test error")
//			}
//		})
//	}
//}

func testParseFile(filename string) error {
	errList := []Error{}
	errHandler := errors.ErrorHandler(func(pos token.Position, typ errors.ErrorType, code errors.ErrCode, params ...interface{}) {
		errList = append(errList, Error{
			Pos: pos,
			Typ: typ,
			Msg: errors.New(pos, code, params...).Error(),
		})
	})

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	sep := "// ==========================="
	code := strings.Split(string(file), sep)
	ctx := preprocess.NewContext()
	r := preprocess.New(ctx, scanner.NewStringScan(filename, code[0], nil), nil)
	p := newMultiparser(r, errHandler)
	t := p.parseUnit()

	dump := ast.String(t, "// ", " ")
	errStr := ast.String(errList, "//", " ")

	if len(code) == 1 {
		return ioutil.WriteFile(filename, []byte(fmt.Sprintf("%s\n%s\n%s%s\n%s%s\n", code[0], sep, dump, sep, errStr, sep)), os.ModePerm)
	}

	if len(code) >= 2 {
		code[1] = strings.ReplaceAll(code[1], "\r\n", "\n")
		if strings.TrimSpace(dump) != strings.TrimSpace(code[1]) {
			return stderr.New(fmt.Sprintf("want ast: code:\n%s\nwant:\n%s\ngot:\n%s\n", code[0], code[1], dump))
		}
	}

	if len(code) >= 3 {
		code[2] = strings.ReplaceAll(code[2], "\r\n", "\n")
		if strings.TrimSpace(errStr) != strings.TrimSpace(code[2]) {
			return stderr.New(fmt.Sprintf("want error: code:\n%s\nwant:\n%s\ngot:\n%s\n", code[0], code[2], errStr))
		}
	}
	return nil
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

func TestParser_ParseUnit(t *testing.T) {
	if err := filepath.Walk("./testdata", func(p string, info os.FileInfo, err error) error {
		ext := filepath.Ext(p)
		name := filepath.Base(p)
		if ext == ".c" {
			t.Run(name, func(t *testing.T) {
				timeOut(t, func(t *testing.T) {
					if err := testParseFile(p); err != nil {
						t.Error(err)
					}
				})
			})
		}
		return nil
	}); err != nil {
		t.Error(err)
	}
}
