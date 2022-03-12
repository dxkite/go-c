int main() {
    return 0;
}

// ===========================
// TranslationUnit
//  `+Files = 
//   `-File
//    |+Name = testdata\simple.c
//    |+Decl = 
//    | `-FuncDecl
//    |  |+Name = main
//    |  |+Type =  int ()
//    |  |+Decl = 
//    |  `+Body = CompoundStmt
//    |   |+Lbrace = testdata\simple.c:1:12
//    |   |+Stmts = 
//    |   | `-ReturnStmt
//    |   |  |+Return = testdata\simple.c:2:5
//    |   |  |+X = 0
//    |   |  `+Semicolon = testdata\simple.c:2:13
//    |   `+Rbrace = testdata\simple.c:3:1
//    `+Unresolved = 
// ===========================
//
// ===========================
