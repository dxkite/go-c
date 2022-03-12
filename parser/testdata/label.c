#include "printf.h"

int main() {
    goto test;
    printf("demo");
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
//   | |  |+Name = printf
//   | |  |+Type =  int (const char *,...)
//   | |  |+Decl = 
//   | |  `+Body = <nil>
//   | `+Unresolved = 
//   `-File
//    |+Name = testdata\label.c
//    |+Decl = 
//    | `-FuncDecl
//    |  |+Name = main
//    |  |+Type =  int ()
//    |  |+Decl = 
//    |  `+Body = CompoundStmt
//    |   |+Lbrace = testdata\label.c:3:12
//    |   |+Stmts = 
//    |   | |-GotoStmt
//    |   | | |+Goto = testdata\label.c:4:5
//    |   | | |+Id = test
//    |   | | `+Semicolon = testdata\label.c:4:14
//    |   | |-ExprStmt
//    |   | | |+Expr = CallExpr
//    |   | | | |+Func = printf
//    |   | | | |+Lparen = testdata\label.c:5:11
//    |   | | | |+Args = 
//    |   | | | | `-"demo"
//    |   | | | `+Rparen = testdata\label.c:5:18
//    |   | | `+Semicolon = testdata\label.c:5:19
//    |   | |-GotoStmt
//    |   | | |+Goto = testdata\label.c:6:5
//    |   | | |+Id = test1
//    |   | | `+Semicolon = testdata\label.c:6:15
//    |   | |-LabelStmt
//    |   | | |+Id = test
//    |   | | `+Stmt = ExprStmt
//    |   | |  |+Expr = CallExpr
//    |   | |  | |+Func = printf
//    |   | |  | |+Lparen = testdata\label.c:8:15
//    |   | |  | |+Args = 
//    |   | |  | | `-"test"
//    |   | |  | `+Rparen = testdata\label.c:8:22
//    |   | |  `+Semicolon = testdata\label.c:8:23
//    |   | `-LabelStmt
//    |   |  |+Id = test2
//    |   |  `+Stmt = ExprStmt
//    |   |   |+Expr = CallExpr
//    |   |   | |+Func = printf
//    |   |   | |+Lparen = testdata\label.c:10:15
//    |   |   | |+Args = 
//    |   |   | | `-"test2"
//    |   |   | `+Rparen = testdata\label.c:10:23
//    |   |   `+Semicolon = testdata\label.c:10:24
//    |   `+Rbrace = testdata\label.c:11:1
//    `+Unresolved = 
// ===========================
//
// `-Error
//  |+Pos = testdata\label.c:6:10
//  |+Typ = 0
//  `+Msg = 在 testdata\label.c 文件的第6行10列: 未定义的标签 test1
// ===========================
