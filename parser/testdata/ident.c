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
//    | |  |+Lbrace = testdata\ident.c:1:23
//    | |  |+Stmts = 
//    | |  | |-IfStmt
//    | |  | | |+If = testdata\ident.c:2:5
//    | |  | | |+X = BinaryExpr
//    | |  | | | |+X = a
//    | |  | | | |+Op = ">"<PUNCTUATOR@testdata\ident.c:2:10>
//    | |  | | | `+Y = BinaryExpr
//    | |  | | |  |+X = b
//    | |  | | |  |+Op = "+"<PUNCTUATOR@testdata\ident.c:2:14>
//    | |  | | |  `+Y = c
//    | |  | | |+Then = CompoundStmt
//    | |  | | | |+Lbrace = testdata\ident.c:2:19
//    | |  | | | |+Stmts = 
//    | |  | | | | `-ReturnStmt
//    | |  | | | |  |+Return = testdata\ident.c:3:9
//    | |  | | | |  |+X = a
//    | |  | | | |  `+Semicolon = testdata\ident.c:3:17
//    | |  | | | `+Rbrace = testdata\ident.c:4:5
//    | |  | | `+Else = <nil>
//    | |  | `-ReturnStmt
//    | |  |  |+Return = testdata\ident.c:5:5
//    | |  |  |+X = b
//    | |  |  `+Semicolon = testdata\ident.c:5:13
//    | |  `+Rbrace = testdata\ident.c:6:1
//    | `-FuncDecl
//    |  |+Name = main
//    |  |+Type =  int ()
//    |  |+Decl = 
//    |  `+Body = CompoundStmt
//    |   |+Lbrace = testdata\ident.c:8:12
//    |   |+Stmts = 
//    |   | |-DeclStmt
//    |   | | `-VarDecl
//    |   | |  |+Type =  int
//    |   | |  |+Name = a
//    |   | |  `+Init = 10
//    |   | `-ReturnStmt
//    |   |  |+Return = testdata\ident.c:10:5
//    |   |  |+X = BinaryExpr
//    |   |  | |+X = CallExpr
//    |   |  | | |+Func = max
//    |   |  | | |+Lparen = testdata\ident.c:10:15
//    |   |  | | |+Args = 
//    |   |  | | | |-10
//    |   |  | | | `-20
//    |   |  | | `+Rparen = testdata\ident.c:10:21
//    |   |  | |+Op = "+"<PUNCTUATOR@testdata\ident.c:10:23>
//    |   |  | `+Y = CallExpr
//    |   |  |  |+Func = min
//    |   |  |  |+Lparen = testdata\ident.c:10:28
//    |   |  |  |+Args = 
//    |   |  |  | |-10
//    |   |  |  | `-CallExpr
//    |   |  |  |  |+Func = max
//    |   |  |  |  |+Lparen = testdata\ident.c:10:36
//    |   |  |  |  |+Args = 
//    |   |  |  |  | |-a
//    |   |  |  |  | `-b
//    |   |  |  |  `+Rparen = testdata\ident.c:10:41
//    |   |  |  `+Rparen = testdata\ident.c:10:42
//    |   |  `+Semicolon = testdata\ident.c:10:43
//    |   `+Rbrace = testdata\ident.c:11:1
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
