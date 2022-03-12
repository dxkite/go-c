int max(int a, int b) {
    if (a>b) {
        return a;
    }
    return b;
}

int main() {
    return max(10,20);
}

// ===========================
// TranslationUnit
//  `+Files = 
//   `-File
//    |+Name = testdata\fundef.c
//    |+Decl = 
//    | |-FuncDecl
//    | | |+Name = max
//    | | |+Type =  int ( int, int)
//    | | |+Decl = 
//    | | `+Body = CompoundStmt
//    | |  |+Lbrace = testdata\fundef.c:1:23
//    | |  |+Stmts = 
//    | |  | |-IfStmt
//    | |  | | |+If = testdata\fundef.c:2:5
//    | |  | | |+X = BinaryExpr
//    | |  | | | |+X = a
//    | |  | | | |+Op = ">"<PUNCTUATOR@testdata\fundef.c:2:10>
//    | |  | | | `+Y = b
//    | |  | | |+Then = CompoundStmt
//    | |  | | | |+Lbrace = testdata\fundef.c:2:14
//    | |  | | | |+Stmts = 
//    | |  | | | | `-ReturnStmt
//    | |  | | | |  |+Return = testdata\fundef.c:3:9
//    | |  | | | |  |+X = a
//    | |  | | | |  `+Semicolon = testdata\fundef.c:3:17
//    | |  | | | `+Rbrace = testdata\fundef.c:4:5
//    | |  | | `+Else = <nil>
//    | |  | `-ReturnStmt
//    | |  |  |+Return = testdata\fundef.c:5:5
//    | |  |  |+X = b
//    | |  |  `+Semicolon = testdata\fundef.c:5:13
//    | |  `+Rbrace = testdata\fundef.c:6:1
//    | `-FuncDecl
//    |  |+Name = main
//    |  |+Type =  int ()
//    |  |+Decl = 
//    |  `+Body = CompoundStmt
//    |   |+Lbrace = testdata\fundef.c:8:12
//    |   |+Stmts = 
//    |   | `-ReturnStmt
//    |   |  |+Return = testdata\fundef.c:9:5
//    |   |  |+X = CallExpr
//    |   |  | |+Func = max
//    |   |  | |+Lparen = testdata\fundef.c:9:15
//    |   |  | |+Args = 
//    |   |  | | |-10
//    |   |  | | `-20
//    |   |  | `+Rparen = testdata\fundef.c:9:21
//    |   |  `+Semicolon = testdata\fundef.c:9:22
//    |   `+Rbrace = testdata\fundef.c:10:1
//    `+Unresolved = 
// ===========================
//
// ===========================
