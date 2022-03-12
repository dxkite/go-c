int main() {
    int a = {.xx = 10, [21] = 13, .x.y.z=10, };
    return 0;
}

// ===========================
// TranslationUnit
//  `+Files = 
//   `-File
//    |+Name = testdata\designator.c
//    |+Decl = 
//    | `-FuncDecl
//    |  |+Name = main
//    |  |+Type =  int ()
//    |  |+Decl = 
//    |  `+Body = CompoundStmt
//    |   |+Lbrace = testdata\designator.c:1:12
//    |   |+Stmts = 
//    |   | |-DeclStmt
//    |   | | `-VarDecl
//    |   | |  |+Type =  int
//    |   | |  |+Name = a
//    |   | |  `+Init = InitializerExpr
//    |   | |   |+Lbrace = testdata\designator.c:2:13
//    |   | |   |+List = 
//    |   | |   | |-RecordDesignatorExpr
//    |   | |   | | |+Dot = testdata\designator.c:2:14
//    |   | |   | | |+Field = xx
//    |   | |   | | `+X = 10
//    |   | |   | |-ArrayDesignatorExpr
//    |   | |   | | |+Lbrack = testdata\designator.c:2:24
//    |   | |   | | |+Index = ConstantExpr
//    |   | |   | | | `+X = 21
//    |   | |   | | |+Rbrack = testdata\designator.c:2:27
//    |   | |   | | `+X = 13
//    |   | |   | `-RecordDesignatorExpr
//    |   | |   |  |+Dot = testdata\designator.c:2:35
//    |   | |   |  |+Field = x
//    |   | |   |  `+X = RecordDesignatorExpr
//    |   | |   |   |+Dot = testdata\designator.c:2:37
//    |   | |   |   |+Field = y
//    |   | |   |   `+X = RecordDesignatorExpr
//    |   | |   |    |+Dot = testdata\designator.c:2:39
//    |   | |   |    |+Field = z
//    |   | |   |    `+X = 10
//    |   | |   `+Rbrace = testdata\designator.c:2:46
//    |   | `-ReturnStmt
//    |   |  |+Return = testdata\designator.c:3:5
//    |   |  |+X = 0
//    |   |  `+Semicolon = testdata\designator.c:3:13
//    |   `+Rbrace = testdata\designator.c:4:1
//    `+Unresolved = 
// ===========================
//
// ===========================
