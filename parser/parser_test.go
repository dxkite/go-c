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
	}

	for _, cas := range testCase {
		t.Run(cas[0], func(t *testing.T) {
			if !testParseType(cas[0], cas[1], cas[2], cas[3]) {
				t.Error("test error")
			}
		})
	}
}
