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
// TranslationUnitDecl
//  `+Decl = 
//   |-FuncDecl
//   | |+Name = Ident
//   | | `+Token = "max"<IDENT@testdata\fundef.c:1:5>
//   | |+Return = BuildInType
//   | | |+Lit = 
//   | | | `-"int"<KEYWORD@testdata\fundef.c:1:1>
//   | | `+Type = int
//   | |+Params = ParamList
//   | | |-ParamVarDecl
//   | | | |+Type = BuildInType
//   | | | | |+Lit = 
//   | | | | | `-"int"<KEYWORD@testdata\fundef.c:1:9>
//   | | | | `+Type = int
//   | | | `+Name = Ident
//   | | |  `+Token = "a"<IDENT@testdata\fundef.c:1:13>
//   | | `-ParamVarDecl
//   | |  |+Type = BuildInType
//   | |  | |+Lit = 
//   | |  | | `-"int"<KEYWORD@testdata\fundef.c:1:16>
//   | |  | `+Type = int
//   | |  `+Name = Ident
//   | |   `+Token = "b"<IDENT@testdata\fundef.c:1:20>
//   | |+Ellipsis = false
//   | |+Decl = 
//   | `+Body = CompoundStmt
//   |  |-IfStmt
//   |  | |+X = BinaryExpr
//   |  | | |+X = Ident
//   |  | | | `+Token = "a"<IDENT@testdata\fundef.c:2:9>
//   |  | | |+Op = ">"<PUNCTUATOR@testdata\fundef.c:2:10>
//   |  | | `+Y = Ident
//   |  | |  `+Token = "b"<IDENT@testdata\fundef.c:2:11>
//   |  | |+Then = CompoundStmt
//   |  | | `-ReturnStmt
//   |  | |  `+X = Ident
//   |  | |   `+Token = "a"<IDENT@testdata\fundef.c:3:16>
//   |  | `+Else = <nil>
//   |  `-ReturnStmt
//   |   `+X = Ident
//   |    `+Token = "b"<IDENT@testdata\fundef.c:5:12>
//   `-FuncDecl
//    |+Name = Ident
//    | `+Token = "main"<IDENT@testdata\fundef.c:8:5>
//    |+Return = BuildInType
//    | |+Lit = 
//    | | `-"int"<KEYWORD@testdata\fundef.c:8:1>
//    | `+Type = int
//    |+Params = ParamList
//    |+Ellipsis = false
//    |+Decl = 
//    `+Body = CompoundStmt
//     `-ReturnStmt
//      `+X = CallExpr
//       |+Fun = Ident
//       | `+Token = "max"<IDENT@testdata\fundef.c:9:12>
//       `+Args = 
//        |-BasicLit
//        | `+Token = "10"<INT@testdata\fundef.c:9:16>
//        `-BasicLit
//         `+Token = "20"<INT@testdata\fundef.c:9:19>
// ===========================
