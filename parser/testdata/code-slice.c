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
//    |  |+Name = Ident
//    |  | `+Token = "main"<IDENT@testdata\code-slice.c:1:5>
//    |  |+Params = ParamList
//    |  |+Ellipsis = false
//    |  |+Return = BuildInType
//    |  | `+Lit = 
//    |  |  `-"int"<KEYWORD@testdata\code-slice.c:1:1>
//    |  |+Decl = 
//    |  `+Body = CompoundStmt
//    |   |-ExprStmt
//    |   | `+Expr = CallExpr
//    |   |  |+Func = Ident
//    |   |  | `+Token = "printf"<IDENT@testdata/code-slice.h:1:1>
//    |   |  `+Args = 
//    |   |   `-BasicLit
//    |   |    `+Token = "\"code\""<STRING@testdata/code-slice.h:1:8>
//    |   |-ExprStmt
//    |   | `+Expr = CallExpr
//    |   |  |+Func = Ident
//    |   |  | `+Token = "printf"<IDENT@testdata/code-slice.h:1:1>
//    |   |  `+Args = 
//    |   |   `-BasicLit
//    |   |    `+Token = "\"code\""<STRING@testdata/code-slice.h:1:8>
//    |   `-ExprStmt
//    |    `+Expr = CallExpr
//    |     |+Func = Ident
//    |     | `+Token = "printf"<IDENT@testdata/code-slice.h:1:1>
//    |     `+Args = 
//    |      `-BasicLit
//    |       `+Token = "\"code\""<STRING@testdata/code-slice.h:1:8>
//    `+Unresolved = 
// ===========================
//
// ===========================
