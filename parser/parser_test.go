package parser

import (
	"dxkite.cn/c/ast"
	"dxkite.cn/c/scanner"
	"dxkite.cn/c/token"
	"fmt"
	"strconv"
	"testing"
)

type Error struct {
	Tok token.Token
	Typ ErrorType
	Msg string
}

func testParseType(name, code, json, err string) bool {
	errList := []Error{}
	errHandler := ErrorHandler(func(tok token.Token, typ ErrorType, msg string) {
		errList = append(errList, Error{
			Tok: tok,
			Typ: typ,
			Msg: msg,
		})
	})
	p := NewParser(scanner.NewStringScan("", code), errHandler)
	t := p.parseTypeName()
	jsonCode, _ := ast.DumpJson(t, true)
	jsonErr, _ := ast.DumpJson(errList, true)
	if string(jsonCode) != json || string(jsonErr) != err {
		code := fmt.Sprintf("{%s,%s,%s,%s},", strconv.QuoteToGraphic(name),
			strconv.QuoteToGraphic(code),
			strconv.QuoteToGraphic(string(jsonCode)),
			strconv.QuoteToGraphic(string(jsonErr)))
		fmt.Println("code:\n", code)
		fmt.Println("want:", json)
		fmt.Println("got", string(jsonCode))
		return false
	}
	return true
}

func TestParser_parseTypeName(t *testing.T) {
	testCase := [][4]string{
		{"simple-type", "int", "{\"@type\":\"BuildInType\",\"Type\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":1,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"}]}\n", "[]\n"},
		{"incomplete-arr-parse", "const int int long *const[]", "{\"@type\":\"IncompleteArrayType\",\"Inner\":{\"@type\":\"PointerType\",\"Inner\":{\"@type\":\"TypeQualifier\",\"Inner\":{\"@type\":\"BuildInType\",\"Type\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":7,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"},{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":11,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"},{\"@type\":\"Token\",\"Lit\":\"long\",\"Pos\":{\"@type\":\"Position\",\"Column\":15,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"}]},\"Qualifier\":{\"const\":true}},\"Qualifier\":{\"const\":true}}}\n", "[]\n"},
		{"const-arr", "const int int long *const[10?2:1]", "{\"@type\":\"ConstArrayType\",\"Inner\":{\"@type\":\"PointerType\",\"Inner\":{\"@type\":\"TypeQualifier\",\"Inner\":{\"@type\":\"BuildInType\",\"Type\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":7,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"},{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":11,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"},{\"@type\":\"Token\",\"Lit\":\"long\",\"Pos\":{\"@type\":\"Position\",\"Column\":15,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"}]},\"Qualifier\":{\"const\":true}},\"Qualifier\":{\"const\":true}},\"Size\":{\"@type\":\"CondExpr\",\"Else\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"1\",\"Pos\":{\"@type\":\"Position\",\"Column\":32,\"Filename\":\"\",\"Line\":1},\"Typ\":\"INT\"}},\"Op\":{\"@type\":\"Token\",\"Lit\":\"?\",\"Pos\":{\"@type\":\"Position\",\"Column\":29,\"Filename\":\"\",\"Line\":1},\"Typ\":\"PUNCTUATOR\"},\"Then\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"2\",\"Pos\":{\"@type\":\"Position\",\"Column\":30,\"Filename\":\"\",\"Line\":1},\"Typ\":\"INT\"}},\"X\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"10\",\"Pos\":{\"@type\":\"Position\",\"Column\":27,\"Filename\":\"\",\"Line\":1},\"Typ\":\"INT\"}}}}\n", "[]\n"},
		{"arr-parse", "const int int long *const[m=10]", "{\"@type\":\"ArrayType\",\"Inner\":{\"@type\":\"PointerType\",\"Inner\":{\"@type\":\"TypeQualifier\",\"Inner\":{\"@type\":\"BuildInType\",\"Type\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":7,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"},{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":11,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"},{\"@type\":\"Token\",\"Lit\":\"long\",\"Pos\":{\"@type\":\"Position\",\"Column\":15,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"}]},\"Qualifier\":{\"const\":true}},\"Qualifier\":{\"const\":true}},\"Qualifier\":{},\"Size\":{\"@type\":\"AssignExpr\",\"Op\":{\"@type\":\"Token\",\"Lit\":\"=\",\"Pos\":{\"@type\":\"Position\",\"Column\":28,\"Filename\":\"\",\"Line\":1},\"Typ\":\"PUNCTUATOR\"},\"X\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"m\",\"Pos\":{\"@type\":\"Position\",\"Column\":27,\"Filename\":\"\",\"Line\":1},\"Typ\":\"IDENT\"}},\"Y\":null},\"Static\":false}\n", "[{\"@type\":\"Error\",\"Msg\":\"expect ) got ]\",\"Tok\":{\"@type\":\"Token\",\"Lit\":\"]\",\"Pos\":{\"@type\":\"Position\",\"Column\":31,\"Filename\":\"\",\"Line\":1},\"Typ\":\"PUNCTUATOR\"},\"Typ\":0}]\n"},
		{"func", "int(*)(int,int)", "{\"@type\":\"FuncType\",\"Ellipsis\":false,\"Inner\":{\"@type\":\"PointerType\",\"Inner\":{\"@type\":\"ParenType\",\"Inner\":{\"@type\":\"BuildInType\",\"Type\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":1,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"}]}},\"Qualifier\":{}},\"Params\":[{\"@type\":\"ParamVarDecl\",\"Name\":null,\"Type\":{\"@type\":\"BuildInType\",\"Type\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":8,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"}]}},{\"@type\":\"ParamVarDecl\",\"Name\":null,\"Type\":{\"@type\":\"BuildInType\",\"Type\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":12,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"}]}}]}\n", "[{\"@type\":\"Error\",\"Msg\":\"expect , got )\",\"Tok\":{\"@type\":\"Token\",\"Lit\":\")\",\"Pos\":{\"@type\":\"Position\",\"Column\":15,\"Filename\":\"\",\"Line\":1},\"Typ\":\"PUNCTUATOR\"},\"Typ\":0}]\n"},
		{"func-with-type", "int(*)(int a,int b)", "{\"@type\":\"FuncType\",\"Ellipsis\":false,\"Inner\":{\"@type\":\"PointerType\",\"Inner\":{\"@type\":\"ParenType\",\"Inner\":{\"@type\":\"BuildInType\",\"Type\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":1,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"}]}},\"Qualifier\":{}},\"Params\":[{\"@type\":\"ParamVarDecl\",\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"a\",\"Pos\":{\"@type\":\"Position\",\"Column\":12,\"Filename\":\"\",\"Line\":1},\"Typ\":\"IDENT\"}},\"Type\":{\"@type\":\"BuildInType\",\"Type\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":8,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"}]}},{\"@type\":\"ParamVarDecl\",\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"b\",\"Pos\":{\"@type\":\"Position\",\"Column\":18,\"Filename\":\"\",\"Line\":1},\"Typ\":\"IDENT\"}},\"Type\":{\"@type\":\"BuildInType\",\"Type\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":14,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"}]}}]}\n", "[{\"@type\":\"Error\",\"Msg\":\"expect , got )\",\"Tok\":{\"@type\":\"Token\",\"Lit\":\")\",\"Pos\":{\"@type\":\"Position\",\"Column\":19,\"Filename\":\"\",\"Line\":1},\"Typ\":\"PUNCTUATOR\"},\"Typ\":0}]\n"},
		{"func-arr-const", "int(*)(int***const a,const int b[a][m])", "{\"@type\":\"FuncType\",\"Ellipsis\":false,\"Inner\":{\"@type\":\"PointerType\",\"Inner\":{\"@type\":\"ParenType\",\"Inner\":{\"@type\":\"BuildInType\",\"Type\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":1,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"}]}},\"Qualifier\":{}},\"Params\":[{\"@type\":\"ParamVarDecl\",\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"a\",\"Pos\":{\"@type\":\"Position\",\"Column\":20,\"Filename\":\"\",\"Line\":1},\"Typ\":\"IDENT\"}},\"Type\":{\"@type\":\"PointerType\",\"Inner\":{\"@type\":\"PointerType\",\"Inner\":{\"@type\":\"PointerType\",\"Inner\":{\"@type\":\"BuildInType\",\"Type\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":8,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"}]},\"Qualifier\":{}},\"Qualifier\":{}},\"Qualifier\":{\"const\":true}}},{\"@type\":\"ParamVarDecl\",\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"b\",\"Pos\":{\"@type\":\"Position\",\"Column\":32,\"Filename\":\"\",\"Line\":1},\"Typ\":\"IDENT\"}},\"Type\":{\"@type\":\"ConstArrayType\",\"Inner\":{\"@type\":\"ConstArrayType\",\"Inner\":{\"@type\":\"TypeQualifier\",\"Inner\":{\"@type\":\"BuildInType\",\"Type\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":28,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"}]},\"Qualifier\":{\"const\":true}},\"Size\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"a\",\"Pos\":{\"@type\":\"Position\",\"Column\":34,\"Filename\":\"\",\"Line\":1},\"Typ\":\"IDENT\"}}},\"Size\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"m\",\"Pos\":{\"@type\":\"Position\",\"Column\":37,\"Filename\":\"\",\"Line\":1},\"Typ\":\"IDENT\"}}}}]}\n", "[{\"@type\":\"Error\",\"Msg\":\"expect , got )\",\"Tok\":{\"@type\":\"Token\",\"Lit\":\")\",\"Pos\":{\"@type\":\"Position\",\"Column\":39,\"Filename\":\"\",\"Line\":1},\"Typ\":\"PUNCTUATOR\"},\"Typ\":0}]\n"},
		{"struct-defined-noname-bit-field", "volatile struct {int a:1?2:3,:9; char c;};", "{\"@type\":\"TypeQualifier\",\"Inner\":{\"@type\":\"RecordType\",\"Fields\":[{\"@type\":\"RecordField\",\"Bit\":{\"@type\":\"ConstantExpr\",\"X\":{\"@type\":\"CondExpr\",\"Else\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"3\",\"Pos\":{\"@type\":\"Position\",\"Column\":28,\"Filename\":\"\",\"Line\":1},\"Typ\":\"INT\"}},\"Op\":{\"@type\":\"Token\",\"Lit\":\"?\",\"Pos\":{\"@type\":\"Position\",\"Column\":25,\"Filename\":\"\",\"Line\":1},\"Typ\":\"PUNCTUATOR\"},\"Then\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"2\",\"Pos\":{\"@type\":\"Position\",\"Column\":26,\"Filename\":\"\",\"Line\":1},\"Typ\":\"INT\"}},\"X\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"1\",\"Pos\":{\"@type\":\"Position\",\"Column\":24,\"Filename\":\"\",\"Line\":1},\"Typ\":\"INT\"}}}},\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"a\",\"Pos\":{\"@type\":\"Position\",\"Column\":22,\"Filename\":\"\",\"Line\":1},\"Typ\":\"IDENT\"}},\"Type\":{\"@type\":\"BuildInType\",\"Type\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":18,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"}]}},{\"@type\":\"RecordField\",\"Bit\":{\"@type\":\"ConstantExpr\",\"X\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"9\",\"Pos\":{\"@type\":\"Position\",\"Column\":31,\"Filename\":\"\",\"Line\":1},\"Typ\":\"INT\"}}},\"Name\":null,\"Type\":{\"@type\":\"BuildInType\",\"Type\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":18,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"}]}},{\"@type\":\"RecordField\",\"Bit\":null,\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"c\",\"Pos\":{\"@type\":\"Position\",\"Column\":39,\"Filename\":\"\",\"Line\":1},\"Typ\":\"IDENT\"}},\"Type\":{\"@type\":\"BuildInType\",\"Type\":[{\"@type\":\"Token\",\"Lit\":\"char\",\"Pos\":{\"@type\":\"Position\",\"Column\":34,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"}]}}],\"Name\":null,\"Type\":{\"@type\":\"Token\",\"Lit\":\"struct\",\"Pos\":{\"@type\":\"Position\",\"Column\":10,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"}},\"Qualifier\":{\"volatile\":true}}\n", "[]\n"},
		{"struct-nested-union", "struct {union{ int a; int b;} const; int c; };", "{\"@type\":\"RecordType\",\"Fields\":[{\"@type\":\"RecordField\",\"Bit\":null,\"Name\":null,\"Type\":{\"@type\":\"TypeQualifier\",\"Inner\":{\"@type\":\"RecordType\",\"Fields\":[{\"@type\":\"RecordField\",\"Bit\":null,\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"a\",\"Pos\":{\"@type\":\"Position\",\"Column\":20,\"Filename\":\"\",\"Line\":1},\"Typ\":\"IDENT\"}},\"Type\":{\"@type\":\"BuildInType\",\"Type\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":16,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"}]}},{\"@type\":\"RecordField\",\"Bit\":null,\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"b\",\"Pos\":{\"@type\":\"Position\",\"Column\":27,\"Filename\":\"\",\"Line\":1},\"Typ\":\"IDENT\"}},\"Type\":{\"@type\":\"BuildInType\",\"Type\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":23,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"}]}}],\"Name\":null,\"Type\":{\"@type\":\"Token\",\"Lit\":\"union\",\"Pos\":{\"@type\":\"Position\",\"Column\":9,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"}},\"Qualifier\":{\"const\":true}}},{\"@type\":\"RecordField\",\"Bit\":null,\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"c\",\"Pos\":{\"@type\":\"Position\",\"Column\":42,\"Filename\":\"\",\"Line\":1},\"Typ\":\"IDENT\"}},\"Type\":{\"@type\":\"BuildInType\",\"Type\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":38,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"}]}}],\"Name\":null,\"Type\":{\"@type\":\"Token\",\"Lit\":\"struct\",\"Pos\":{\"@type\":\"Position\",\"Column\":1,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"}}\n", "[]\n"},
		//添加报错
		{"struct-nested-union", "struct test {struct test const; int c; };", "{\"@type\":\"RecordType\",\"Fields\":[{\"@type\":\"RecordField\",\"Bit\":null,\"Name\":null,\"Type\":{\"@type\":\"TypeQualifier\",\"Inner\":{\"@type\":\"RecordType\",\"Fields\":[],\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"test\",\"Pos\":{\"@type\":\"Position\",\"Column\":21,\"Filename\":\"\",\"Line\":1},\"Typ\":\"IDENT\"}},\"Type\":{\"@type\":\"Token\",\"Lit\":\"struct\",\"Pos\":{\"@type\":\"Position\",\"Column\":14,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"}},\"Qualifier\":{\"const\":true}}},{\"@type\":\"RecordField\",\"Bit\":null,\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"c\",\"Pos\":{\"@type\":\"Position\",\"Column\":37,\"Filename\":\"\",\"Line\":1},\"Typ\":\"IDENT\"}},\"Type\":{\"@type\":\"BuildInType\",\"Type\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":33,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"}]}}],\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"test\",\"Pos\":{\"@type\":\"Position\",\"Column\":8,\"Filename\":\"\",\"Line\":1},\"Typ\":\"IDENT\"}},\"Type\":{\"@type\":\"Token\",\"Lit\":\"struct\",\"Pos\":{\"@type\":\"Position\",\"Column\":1,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"}}\n", "[]\n"},
		{"struct-nested-union-pointer", "struct test {struct test const * left; int c; };", "{\"@type\":\"RecordType\",\"Fields\":[{\"@type\":\"RecordField\",\"Bit\":null,\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"left\",\"Pos\":{\"@type\":\"Position\",\"Column\":34,\"Filename\":\"\",\"Line\":1},\"Typ\":\"IDENT\"}},\"Type\":{\"@type\":\"PointerType\",\"Inner\":{\"@type\":\"TypeQualifier\",\"Inner\":{\"@type\":\"RecordType\",\"Fields\":[],\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"test\",\"Pos\":{\"@type\":\"Position\",\"Column\":21,\"Filename\":\"\",\"Line\":1},\"Typ\":\"IDENT\"}},\"Type\":{\"@type\":\"Token\",\"Lit\":\"struct\",\"Pos\":{\"@type\":\"Position\",\"Column\":14,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"}},\"Qualifier\":{\"const\":true}},\"Qualifier\":{}}},{\"@type\":\"RecordField\",\"Bit\":null,\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"c\",\"Pos\":{\"@type\":\"Position\",\"Column\":44,\"Filename\":\"\",\"Line\":1},\"Typ\":\"IDENT\"}},\"Type\":{\"@type\":\"BuildInType\",\"Type\":[{\"@type\":\"Token\",\"Lit\":\"int\",\"Pos\":{\"@type\":\"Position\",\"Column\":40,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"}]}}],\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"test\",\"Pos\":{\"@type\":\"Position\",\"Column\":8,\"Filename\":\"\",\"Line\":1},\"Typ\":\"IDENT\"}},\"Type\":{\"@type\":\"Token\",\"Lit\":\"struct\",\"Pos\":{\"@type\":\"Position\",\"Column\":1,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"}}\n", "[]\n"},
	}

	for _, cas := range testCase {
		t.Run(cas[0], func(t *testing.T) {
			if !testParseType(cas[0], cas[1], cas[2], cas[3]) {
				t.Error("test error")
			}
		})
	}
}

func testParseExpr(name, code, json, err string) bool {
	errList := []Error{}
	errHandler := ErrorHandler(func(tok token.Token, typ ErrorType, msg string) {
		errList = append(errList, Error{
			Tok: tok,
			Typ: typ,
			Msg: msg,
		})
	})
	p := NewParser(scanner.NewStringScan("", code), errHandler)
	t := p.parseAssignExpr()
	jsonCode, _ := ast.DumpJson(t, true)
	jsonErr, _ := ast.DumpJson(errList, true)
	if string(jsonCode) != json || string(jsonErr) != err {
		code := fmt.Sprintf("{%s,%s,%s,%s},", strconv.QuoteToGraphic(name),
			strconv.QuoteToGraphic(code),
			strconv.QuoteToGraphic(string(jsonCode)),
			strconv.QuoteToGraphic(string(jsonErr)))
		fmt.Println("code:\n", code)
		fmt.Println("want:", json)
		fmt.Println("got", string(jsonCode))
		return false
	}
	return true
}

func TestParser_parseAssignExpr(t *testing.T) {
	testCase := [][4]string{
		{"initializer-expr", "answer = (struct point){ .quot = 2, .rem = -1 }", "{\"@type\":\"AssignExpr\",\"Op\":{\"@type\":\"Token\",\"Lit\":\"=\",\"Pos\":{\"@type\":\"Position\",\"Column\":8,\"Filename\":\"\",\"Line\":1},\"Typ\":\"PUNCTUATOR\"},\"X\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"answer\",\"Pos\":{\"@type\":\"Position\",\"Column\":1,\"Filename\":\"\",\"Line\":1},\"Typ\":\"IDENT\"}},\"Y\":{\"@type\":\"CompoundLit\",\"InitList\":[{\"@type\":\"RecordDesignatorExpr\",\"Field\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"quot\",\"Pos\":{\"@type\":\"Position\",\"Column\":27,\"Filename\":\"\",\"Line\":1},\"Typ\":\"IDENT\"}},\"X\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"2\",\"Pos\":{\"@type\":\"Position\",\"Column\":34,\"Filename\":\"\",\"Line\":1},\"Typ\":\"INT\"}}},{\"@type\":\"RecordDesignatorExpr\",\"Field\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"rem\",\"Pos\":{\"@type\":\"Position\",\"Column\":38,\"Filename\":\"\",\"Line\":1},\"Typ\":\"IDENT\"}},\"X\":{\"@type\":\"UnaryExpr\",\"Op\":{\"@type\":\"Token\",\"Lit\":\"-\",\"Pos\":{\"@type\":\"Position\",\"Column\":44,\"Filename\":\"\",\"Line\":1},\"Typ\":\"PUNCTUATOR\"},\"X\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"1\",\"Pos\":{\"@type\":\"Position\",\"Column\":45,\"Filename\":\"\",\"Line\":1},\"Typ\":\"INT\"}}}}],\"Type\":{\"@type\":\"RecordType\",\"Fields\":[],\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"point\",\"Pos\":{\"@type\":\"Position\",\"Column\":18,\"Filename\":\"\",\"Line\":1},\"Typ\":\"IDENT\"}},\"Type\":{\"@type\":\"Token\",\"Lit\":\"struct\",\"Pos\":{\"@type\":\"Position\",\"Column\":11,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"}}}}\n", "[]\n"},
		{"initializer-expr-comma", "answer = (struct point){ .quot = 2, .rem = -1, }", "{\"@type\":\"AssignExpr\",\"Op\":{\"@type\":\"Token\",\"Lit\":\"=\",\"Pos\":{\"@type\":\"Position\",\"Column\":8,\"Filename\":\"\",\"Line\":1},\"Typ\":\"PUNCTUATOR\"},\"X\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"answer\",\"Pos\":{\"@type\":\"Position\",\"Column\":1,\"Filename\":\"\",\"Line\":1},\"Typ\":\"IDENT\"}},\"Y\":{\"@type\":\"CompoundLit\",\"InitList\":[{\"@type\":\"RecordDesignatorExpr\",\"Field\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"quot\",\"Pos\":{\"@type\":\"Position\",\"Column\":27,\"Filename\":\"\",\"Line\":1},\"Typ\":\"IDENT\"}},\"X\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"2\",\"Pos\":{\"@type\":\"Position\",\"Column\":34,\"Filename\":\"\",\"Line\":1},\"Typ\":\"INT\"}}},{\"@type\":\"RecordDesignatorExpr\",\"Field\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"rem\",\"Pos\":{\"@type\":\"Position\",\"Column\":38,\"Filename\":\"\",\"Line\":1},\"Typ\":\"IDENT\"}},\"X\":{\"@type\":\"UnaryExpr\",\"Op\":{\"@type\":\"Token\",\"Lit\":\"-\",\"Pos\":{\"@type\":\"Position\",\"Column\":44,\"Filename\":\"\",\"Line\":1},\"Typ\":\"PUNCTUATOR\"},\"X\":{\"@type\":\"BasicLit\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"1\",\"Pos\":{\"@type\":\"Position\",\"Column\":45,\"Filename\":\"\",\"Line\":1},\"Typ\":\"INT\"}}}}],\"Type\":{\"@type\":\"RecordType\",\"Fields\":[],\"Name\":{\"@type\":\"Ident\",\"Token\":{\"@type\":\"Token\",\"Lit\":\"point\",\"Pos\":{\"@type\":\"Position\",\"Column\":18,\"Filename\":\"\",\"Line\":1},\"Typ\":\"IDENT\"}},\"Type\":{\"@type\":\"Token\",\"Lit\":\"struct\",\"Pos\":{\"@type\":\"Position\",\"Column\":11,\"Filename\":\"\",\"Line\":1},\"Typ\":\"KEYWORD\"}}}}\n", "[]\n"},
	}

	for _, cas := range testCase {
		t.Run(cas[0], func(t *testing.T) {
			if !testParseExpr(cas[0], cas[1], cas[2], cas[3]) {
				t.Error("test error")
			}
		})
	}
}
