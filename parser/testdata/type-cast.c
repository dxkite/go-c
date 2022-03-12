int main() {
    int a;
    char b = (char)a;
    return 0;
}

// ===========================
// TranslationUnit
//  `+Files = 
//   `-File
//    |+Name = testdata\type-cast.c
//    |+Decl = 
//    | `-FuncDecl
//    |  |+Name = main
//    |  |+Type =  int ()
//    |  |+Decl = 
//    |  `+Body = CompoundStmt
//    |   |+Lbrace = testdata\type-cast.c:1:12
//    |   |+Stmts = 
//    |   | |-DeclStmt
//    |   | | `-VarDecl
//    |   | |  |+Type =  int
//    |   | |  |+Name = a
//    |   | |  `+Init = <nil>
//    |   | |-DeclStmt
//    |   | | `-VarDecl
//    |   | |  |+Type =  char
//    |   | |  |+Name = b
//    |   | |  `+Init = TypeCastExpr
//    |   | |   |+Lparen = testdata\type-cast.c:3:14
//    |   | |   |+Type =  char
//    |   | |   |+Rparen = testdata\type-cast.c:3:19
//    |   | |   `+X = a
//    |   | `-ReturnStmt
//    |   |  |+Return = testdata\type-cast.c:4:5
//    |   |  |+X = 0
//    |   |  `+Semicolon = testdata\type-cast.c:4:13
//    |   `+Rbrace = testdata\type-cast.c:5:1
//    `+Unresolved = 
// ===========================
//
// ===========================
