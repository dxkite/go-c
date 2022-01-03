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
//    |   |-ExprStmt
//    |   | `+Expr = CallExpr
//    |   |  |+Func = printf
//    |   |  `+Args = 
//    |   |   `-"code"
//    |   |-ExprStmt
//    |   | `+Expr = CallExpr
//    |   |  |+Func = printf
//    |   |  `+Args = 
//    |   |   `-"code"
//    |   `-ExprStmt
//    |    `+Expr = CallExpr
//    |     |+Func = printf
//    |     `+Args = 
//    |      `-"code"
//    `+Unresolved = 
// ===========================
//
// ===========================
