#include "printf.h"

int main() {
    goto test;
    print("demo");
    goto test1;
    test:
        printf("test");
    test2:
        printf("test2");
}

// ===========================
// TranslationUnit
//  `+Files = 
//   |-File
//   | |+Name = testdata/printf.h
//   | |+Decl = 
//   | | `-FuncDecl
//   | |  |+Name = Ident
//   | |  | `+Token = "printf"<IDENT@testdata/printf.h:1:5>
//   | |  |+Params = ParamList
//   | |  | `-ParamVarDecl
//   | |  |  |+Name = Ident
//   | |  |  | `+Token = "fmt"<IDENT@testdata/printf.h:1:24>
//   | |  |  `+Type = PointerType
//   | |  |   `+Inner = TypeQualifier
//   | |  |    |+Qualifier = map[const:true]
//   | |  |    `+Inner = BuildInType
//   | |  |     `+Lit = 
//   | |  |      `-"char"<KEYWORD@testdata/printf.h:1:18>
//   | |  |+Ellipsis = true
//   | |  |+Return = BuildInType
//   | |  | `+Lit = 
//   | |  |  `-"int"<KEYWORD@testdata/printf.h:1:1>
//   | |  |+Decl = 
//   | |  `+Body = <nil>
//   | `+Unresolved = 
//   `-File
//    |+Name = testdata\label.c
//    |+Decl = 
//    | `-FuncDecl
//    |  |+Name = Ident
//    |  | `+Token = "main"<IDENT@testdata\label.c:3:5>
//    |  |+Params = ParamList
//    |  |+Ellipsis = false
//    |  |+Return = BuildInType
//    |  | `+Lit = 
//    |  |  `-"int"<KEYWORD@testdata\label.c:3:1>
//    |  |+Decl = 
//    |  `+Body = CompoundStmt
//    |   |-GotoStmt
//    |   | `+Id = Ident
//    |   |  `+Token = "test"<IDENT@testdata\label.c:4:10>
//    |   |-ExprStmt
//    |   | `+Expr = CallExpr
//    |   |  |+Func = Ident
//    |   |  | `+Token = "print"<IDENT@testdata\label.c:5:5>
//    |   |  `+Args = 
//    |   |   `-BasicLit
//    |   |    `+Token = "\"demo\""<STRING@testdata\label.c:5:11>
//    |   |-GotoStmt
//    |   | `+Id = Ident
//    |   |  `+Token = "test1"<IDENT@testdata\label.c:6:10>
//    |   |-LabelStmt
//    |   | |+Id = Ident
//    |   | | `+Token = "test"<IDENT@testdata\label.c:7:5>
//    |   | `+Stmt = ExprStmt
//    |   |  `+Expr = CallExpr
//    |   |   |+Func = Ident
//    |   |   | `+Token = "printf"<IDENT@testdata\label.c:8:9>
//    |   |   `+Args = 
//    |   |    `-BasicLit
//    |   |     `+Token = "\"test\""<STRING@testdata\label.c:8:16>
//    |   `-LabelStmt
//    |    |+Id = Ident
//    |    | `+Token = "test2"<IDENT@testdata\label.c:9:5>
//    |    `+Stmt = ExprStmt
//    |     `+Expr = CallExpr
//    |      |+Func = Ident
//    |      | `+Token = "printf"<IDENT@testdata\label.c:10:9>
//    |      `+Args = 
//    |       `-BasicLit
//    |        `+Token = "\"test2\""<STRING@testdata\label.c:10:16>
//    `+Unresolved = 
// ===========================
//
// `-Error
//  |+Pos = testdata\label.c:6:10
//  |+Typ = 0
//  `+Msg = 在 testdata\label.c 文件的第6行10列: 未定义的标签 test1
// ===========================
