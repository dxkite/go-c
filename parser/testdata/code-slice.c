#include "printf.h"

int main() {
#include "code-slice.h"
#include "code-slice.h"
#include "code-slice.h"
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
//    |+Name = testdata\code-slice.c
//    |+Decl = 
//    | `-FuncDecl
//    |  |+Name = main
//    |  |+Type =  int ()
//    |  |+Decl = 
//    |  `+Body = CompoundStmt
//    |   |+Lbrace = testdata\code-slice.c:3:12
//    |   |+Stmts = 
//    |   | |-ExprStmt
//    |   | | |+Expr = CallExpr
//    |   | | | |+Func = printf
//    |   | | | |+Lparen = testdata/code-slice.h:1:7
//    |   | | | |+Args = 
//    |   | | | | `-"code"
//    |   | | | `+Rparen = testdata/code-slice.h:1:14
//    |   | | `+Semicolon = testdata/code-slice.h:1:15
//    |   | |-ExprStmt
//    |   | | |+Expr = CallExpr
//    |   | | | |+Func = printf
//    |   | | | |+Lparen = testdata/code-slice.h:1:7
//    |   | | | |+Args = 
//    |   | | | | `-"code"
//    |   | | | `+Rparen = testdata/code-slice.h:1:14
//    |   | | `+Semicolon = testdata/code-slice.h:1:15
//    |   | `-ExprStmt
//    |   |  |+Expr = CallExpr
//    |   |  | |+Func = printf
//    |   |  | |+Lparen = testdata/code-slice.h:1:7
//    |   |  | |+Args = 
//    |   |  | | `-"code"
//    |   |  | `+Rparen = testdata/code-slice.h:1:14
//    |   |  `+Semicolon = testdata/code-slice.h:1:15
//    |   `+Rbrace = testdata\code-slice.c:7:1
//    `+Unresolved = 
// ===========================
//
// ===========================
