int main() {
    return 0;
}
// ===========================
// TranslationUnitDecl
//  `+Decl = 
//   `-FuncDecl
//    |+Name = Ident
//    | `+Token = "main"<IDENT@testdata\simple.c:1:5>
//    |+Return = BuildInType
//    | |+Lit = 
//    | | `-"int"<KEYWORD@testdata\simple.c:1:1>
//    | `+Type = int
//    |+Params = ParamList
//    |+Ellipsis = false
//    |+Decl = 
//    `+Body = CompoundStmt
//     `-ReturnStmt
//      `+X = BasicLit
//       `+Token = "0"<INT@testdata\simple.c:2:12>
// ===========================
