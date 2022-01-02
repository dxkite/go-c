typedef enum Color {
    BLUE = 1,
    GREEN,
    YELLOW,
} color_t;

int main() {
    color_t color1;
    enum Color color2;
}

// ===========================
// TranslationUnit
//  `+Files = 
//   `-File
//    |+Name = testdata\enum.c
//    |+Decl = 
//    | |-TypedefDecl
//    | | |+Name = Ident
//    | | | `+Token = "color_t"<IDENT@testdata\enum.c:5:3>
//    | | `+Type = EnumType
//    | |  |+Name = Ident
//    | |  | `+Token = "Color"<IDENT@testdata\enum.c:1:14>
//    | |  `+List = EnumFieldList
//    | |   |-EnumFieldDecl
//    | |   | |+Name = Ident
//    | |   | | `+Token = "BLUE"<IDENT@testdata\enum.c:2:5>
//    | |   | `+Val = ConstantExpr
//    | |   |  `+X = BasicLit
//    | |   |   `+Token = "1"<INT@testdata\enum.c:2:12>
//    | |   |-EnumFieldDecl
//    | |   | |+Name = Ident
//    | |   | | `+Token = "GREEN"<IDENT@testdata\enum.c:3:5>
//    | |   | `+Val = <nil>
//    | |   `-EnumFieldDecl
//    | |    |+Name = Ident
//    | |    | `+Token = "YELLOW"<IDENT@testdata\enum.c:4:5>
//    | |    `+Val = <nil>
//    | `-FuncDecl
//    |  |+Name = Ident
//    |  | `+Token = "main"<IDENT@testdata\enum.c:7:5>
//    |  |+Params = ParamList
//    |  |+Ellipsis = false
//    |  |+Return = BuildInType
//    |  | `+Lit = 
//    |  |  `-"int"<KEYWORD@testdata\enum.c:7:1>
//    |  |+Decl = 
//    |  `+Body = CompoundStmt
//    |   |-DeclStmt
//    |   | `-VarDecl
//    |   |  |+Name = Ident
//    |   |  | `+Token = "color1"<IDENT@testdata\enum.c:8:13>
//    |   |  |+Type = EnumType
//    |   |  | |+Name = Ident
//    |   |  | | `+Token = "Color"<IDENT@testdata\enum.c:1:14>
//    |   |  | `+List = EnumFieldList
//    |   |  |  |-EnumFieldDecl
//    |   |  |  | |+Name = Ident
//    |   |  |  | | `+Token = "BLUE"<IDENT@testdata\enum.c:2:5>
//    |   |  |  | `+Val = ConstantExpr
//    |   |  |  |  `+X = BasicLit
//    |   |  |  |   `+Token = "1"<INT@testdata\enum.c:2:12>
//    |   |  |  |-EnumFieldDecl
//    |   |  |  | |+Name = Ident
//    |   |  |  | | `+Token = "GREEN"<IDENT@testdata\enum.c:3:5>
//    |   |  |  | `+Val = <nil>
//    |   |  |  `-EnumFieldDecl
//    |   |  |   |+Name = Ident
//    |   |  |   | `+Token = "YELLOW"<IDENT@testdata\enum.c:4:5>
//    |   |  |   `+Val = <nil>
//    |   |  `+Init = <nil>
//    |   `-DeclStmt
//    |    `-VarDecl
//    |     |+Name = Ident
//    |     | `+Token = "color2"<IDENT@testdata\enum.c:9:16>
//    |     |+Type = EnumType
//    |     | |+Name = Ident
//    |     | | `+Token = "Color"<IDENT@testdata\enum.c:9:10>
//    |     | `+List = EnumFieldList
//    |     `+Init = <nil>
//    `+Unresolved = 
// ===========================
//
// ===========================
