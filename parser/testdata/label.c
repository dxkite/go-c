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
//   | |  |+Params = ParamList
//   | |  | `-ParamVarDecl
//   | |  |  |+Qua = map[]
//   | |  |  |+Name = fmt
//   | |  |  `+Type = const char *
//   | |  |+Ellipsis = true
//   | |  |+Return =  int
//   | |  |+Decl = 
//   | |  `+Body = <nil>
//   | `+Unresolved = 
//   `-File
//    |+Name = testdata\label.c
//    |+Decl = 
//    | `-FuncDecl
//    |  |+Name = main
//    |  |+Params = ParamList
//    |  |+Ellipsis = false
//    |  |+Return =  int
//    |  |+Decl = 
//    |  `+Body = CompoundStmt
//    |   |-GotoStmt
//    |   | `+Id = test
//    |   |-ExprStmt
//    |   | `+Expr = CallExpr
//    |   |  |+Func = printf
//    |   |  `+Args = 
//    |   |   `-"demo"
//    |   |-GotoStmt
//    |   | `+Id = test1
//    |   |-LabelStmt
//    |   | |+Id = test
//    |   | `+Stmt = ExprStmt
//    |   |  `+Expr = CallExpr
//    |   |   |+Func = printf
//    |   |   `+Args = 
//    |   |    `-"test"
//    |   `-LabelStmt
//    |    |+Id = test2
//    |    `+Stmt = ExprStmt
//    |     `+Expr = CallExpr
//    |      |+Func = printf
//    |      `+Args = 
//    |       `-"test2"
//    `+Unresolved = 
// ===========================
//
// `-Error
//  |+Pos = testdata\label.c:6:10
//  |+Typ = 0
//  `+Msg = 在 testdata\label.c 文件的第6行10列: 未定义的标签 test1
// ===========================
