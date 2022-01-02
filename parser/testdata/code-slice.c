int main() {
#include "code-slice.h"
#include "code-slice.h"
#include "code-slice.h"
}

// ===========================
// TranslationUnit
//  `+Files = 
//   `-File
//    |+Name = testdata\code-slice.c
//    |+Decl = 
//    | `-FuncDecl
//    |  |+Name = main
//    |  |+Params = ParamList
//    |  |+Ellipsis = false
//    |  |+Return =  int
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
