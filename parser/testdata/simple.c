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
//    |  |+Name = Ident
//    |  | `+Token = "main"<IDENT@testdata\simple.c:1:5>
//    |  |+Params = ParamList
//    |  |+Ellipsis = false
//    |  |+Return = BuildInType
//    |  | |+Lit = 
//    |  | | `-"int"<KEYWORD@testdata\simple.c:1:1>
//    |  | `+Type = int
//    |  |+Decl = 
//    |  `+Body = CompoundStmt
//    |   `-ReturnStmt
//    |    `+X = BasicLit
//    |     `+Token = "0"<INT@testdata\simple.c:2:12>
//    `+Unresolved = 
// ===========================
//
// ===========================
