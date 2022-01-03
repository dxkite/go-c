int max(int a, int b) {
    if (a> b + c) {
        return a;
    }
    return b;
}

int main() {
    int a = 10;
    return max(10,20) + min(10, max(a, b));
}

// ===========================
// TranslationUnit
//  `+Files = 
//   `-File
//    |+Name = testdata\ident.c
//    |+Decl = 
//    | |-FuncDecl
//    | | |+Name = max
//    | | |+Type =  int ( int, int)
//    | | |+Decl = 
//    | | `+Body = CompoundStmt
//    | |  |-IfStmt
//    | |  | |+X = BinaryExpr
//    | |  | | |+X = a
//    | |  | | |+Op = ">"<PUNCTUATOR@testdata\ident.c:2:10>
//    | |  | | `+Y = BinaryExpr
//    | |  | |  |+X = b
//    | |  | |  |+Op = "+"<PUNCTUATOR@testdata\ident.c:2:14>
//    | |  | |  `+Y = c
//    | |  | |+Then = CompoundStmt
//    | |  | | `-ReturnStmt
//    | |  | |  `+X = a
//    | |  | `+Else = <nil>
//    | |  `-ReturnStmt
//    | |   `+X = b
//    | `-FuncDecl
//    |  |+Name = main
//    |  |+Type =  int ()
//    |  |+Decl = 
//    |  `+Body = CompoundStmt
//    |   |-DeclStmt
//    |   | `-VarDecl
//    |   |  |+Name = a
//    |   |  |+Type =  int
//    |   |  `+Init = 10
//    |   `-ReturnStmt
//    |    `+X = BinaryExpr
//    |     |+X = CallExpr
//    |     | |+Func = max
//    |     | `+Args = 
//    |     |  |-10
//    |     |  `-20
//    |     |+Op = "+"<PUNCTUATOR@testdata\ident.c:10:23>
//    |     `+Y = CallExpr
//    |      |+Func = min
//    |      `+Args = 
//    |       |-10
//    |       `-CallExpr
//    |        |+Func = max
//    |        `+Args = 
//    |         |-a
//    |         `-b
//    `+Unresolved = 
// ===========================
//
// |-Error
// | |+Pos = testdata\ident.c:2:16
// | |+Typ = 0
// | `+Msg = 在 testdata\ident.c 文件的第2行16列: 未定义的标识符 c
// |-Error
// | |+Pos = testdata\ident.c:10:25
// | |+Typ = 0
// | `+Msg = 在 testdata\ident.c 文件的第10行25列: 未定义的标识符 min
// `-Error
//  |+Pos = testdata\ident.c:10:40
//  |+Typ = 0
//  `+Msg = 在 testdata\ident.c 文件的第10行40列: 未定义的标识符 b
// ===========================
