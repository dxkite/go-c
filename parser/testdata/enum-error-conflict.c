typedef enum Color {
    BLUE = 1,
    GREEN,
    YELLOW,
    YELLOW,
} color_t;

int main() {
    color_t color1;
    int a = YELLOW;
    int YELLOW = 10;
}

// ===========================
// TranslationUnit
//  `+Files = 
//   `-File
//    |+Name = testdata\enum-error-conflict.c
//    |+Decl = 
//    | |-TypedefDecl
//    | | |+Name = Ident
//    | | | `+Token = "color_t"<IDENT@testdata\enum-error-conflict.c:6:3>
//    | | `+Type = EnumType
//    | |  |+Name = Ident
//    | |  | `+Token = "Color"<IDENT@testdata\enum-error-conflict.c:1:14>
//    | |  `+List = EnumFieldList
//    | |   |-EnumFieldDecl
//    | |   | |+Name = Ident
//    | |   | | `+Token = "BLUE"<IDENT@testdata\enum-error-conflict.c:2:5>
//    | |   | `+Val = ConstantExpr
//    | |   |  `+X = BasicLit
//    | |   |   `+Token = "1"<INT@testdata\enum-error-conflict.c:2:12>
//    | |   |-EnumFieldDecl
//    | |   | |+Name = Ident
//    | |   | | `+Token = "GREEN"<IDENT@testdata\enum-error-conflict.c:3:5>
//    | |   | `+Val = <nil>
//    | |   |-EnumFieldDecl
//    | |   | |+Name = Ident
//    | |   | | `+Token = "YELLOW"<IDENT@testdata\enum-error-conflict.c:4:5>
//    | |   | `+Val = <nil>
//    | |   `-EnumFieldDecl
//    | |    |+Name = Ident
//    | |    | `+Token = "YELLOW"<IDENT@testdata\enum-error-conflict.c:5:5>
//    | |    `+Val = <nil>
//    | `-FuncDecl
//    |  |+Name = Ident
//    |  | `+Token = "main"<IDENT@testdata\enum-error-conflict.c:8:5>
//    |  |+Params = ParamList
//    |  |+Ellipsis = false
//    |  |+Return = BuildInType
//    |  | |+Lit = 
//    |  | | `-"int"<KEYWORD@testdata\enum-error-conflict.c:8:1>
//    |  | `+Type = int
//    |  |+Decl = 
//    |  `+Body = CompoundStmt
//    |   |-DeclStmt
//    |   | `-VarDecl
//    |   |  |+Name = Ident
//    |   |  | `+Token = "color1"<IDENT@testdata\enum-error-conflict.c:9:13>
//    |   |  |+Type = EnumType
//    |   |  | |+Name = Ident
//    |   |  | | `+Token = "Color"<IDENT@testdata\enum-error-conflict.c:1:14>
//    |   |  | `+List = EnumFieldList
//    |   |  |  |-EnumFieldDecl
//    |   |  |  | |+Name = Ident
//    |   |  |  | | `+Token = "BLUE"<IDENT@testdata\enum-error-conflict.c:2:5>
//    |   |  |  | `+Val = ConstantExpr
//    |   |  |  |  `+X = BasicLit
//    |   |  |  |   `+Token = "1"<INT@testdata\enum-error-conflict.c:2:12>
//    |   |  |  |-EnumFieldDecl
//    |   |  |  | |+Name = Ident
//    |   |  |  | | `+Token = "GREEN"<IDENT@testdata\enum-error-conflict.c:3:5>
//    |   |  |  | `+Val = <nil>
//    |   |  |  |-EnumFieldDecl
//    |   |  |  | |+Name = Ident
//    |   |  |  | | `+Token = "YELLOW"<IDENT@testdata\enum-error-conflict.c:4:5>
//    |   |  |  | `+Val = <nil>
//    |   |  |  `-EnumFieldDecl
//    |   |  |   |+Name = Ident
//    |   |  |   | `+Token = "YELLOW"<IDENT@testdata\enum-error-conflict.c:5:5>
//    |   |  |   `+Val = <nil>
//    |   |  `+Init = <nil>
//    |   |-DeclStmt
//    |   | `-VarDecl
//    |   |  |+Name = Ident
//    |   |  | `+Token = "a"<IDENT@testdata\enum-error-conflict.c:10:9>
//    |   |  |+Type = BuildInType
//    |   |  | |+Lit = 
//    |   |  | | `-"int"<KEYWORD@testdata\enum-error-conflict.c:10:5>
//    |   |  | `+Type = int
//    |   |  `+Init = Ident
//    |   |   `+Token = "YELLOW"<IDENT@testdata\enum-error-conflict.c:10:13>
//    |   `-DeclStmt
//    |    `-VarDecl
//    |     |+Name = Ident
//    |     | `+Token = "YELLOW"<IDENT@testdata\enum-error-conflict.c:11:9>
//    |     |+Type = BuildInType
//    |     | |+Lit = 
//    |     | | `-"int"<KEYWORD@testdata\enum-error-conflict.c:11:5>
//    |     | `+Type = int
//    |     `+Init = BasicLit
//    |      `+Token = "10"<INT@testdata\enum-error-conflict.c:11:18>
//    `+Unresolved = 
// ===========================
//
// `-Error
//  |+Pos = testdata\enum-error-conflict.c:5:5
//  |+Typ = 0
//  `+Msg = 在 testdata\enum-error-conflict.c 文件的第5行5列: 重复的标识符 YELLOW，上次声明的位置 testdata\enum-error-conflict.c:4:5
// ===========================
