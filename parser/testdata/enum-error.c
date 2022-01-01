typedef enum Color {
    BLUE = 1,
    GREEN,
    YELLOW,
} color_t;

int main() {
    color_t color1;
    int color1;
    char color2;
    enum long Color color2;
    long enum Color color2;
    const long enum Color color3;
}

// ===========================
// TranslationUnit
//  `+Files = 
//   `-File
//    |+Name = testdata\enum-error.c
//    |+Decl = 
//    | |-TypedefDecl
//    | | |+Name = Ident
//    | | | `+Token = "color_t"<IDENT@testdata\enum-error.c:5:3>
//    | | `+Type = EnumType
//    | |  |+Name = Ident
//    | |  | `+Token = "Color"<IDENT@testdata\enum-error.c:1:14>
//    | |  `+List = EnumFieldList
//    | |   |-EnumFieldDecl
//    | |   | |+Name = Ident
//    | |   | | `+Token = "BLUE"<IDENT@testdata\enum-error.c:2:5>
//    | |   | `+Val = ConstantExpr
//    | |   |  `+X = BasicLit
//    | |   |   `+Token = "1"<INT@testdata\enum-error.c:2:12>
//    | |   |-EnumFieldDecl
//    | |   | |+Name = Ident
//    | |   | | `+Token = "GREEN"<IDENT@testdata\enum-error.c:3:5>
//    | |   | `+Val = <nil>
//    | |   `-EnumFieldDecl
//    | |    |+Name = Ident
//    | |    | `+Token = "YELLOW"<IDENT@testdata\enum-error.c:4:5>
//    | |    `+Val = <nil>
//    | `-FuncDecl
//    |  |+Name = Ident
//    |  | `+Token = "main"<IDENT@testdata\enum-error.c:7:5>
//    |  |+Params = ParamList
//    |  |+Ellipsis = false
//    |  |+Return = BuildInType
//    |  | |+Lit = 
//    |  | | `-"int"<KEYWORD@testdata\enum-error.c:7:1>
//    |  | `+Type = int
//    |  |+Decl = 
//    |  `+Body = CompoundStmt
//    |   |-DeclStmt
//    |   | `-VarDecl
//    |   |  |+Name = Ident
//    |   |  | `+Token = "color1"<IDENT@testdata\enum-error.c:8:13>
//    |   |  |+Type = EnumType
//    |   |  | |+Name = Ident
//    |   |  | | `+Token = "Color"<IDENT@testdata\enum-error.c:1:14>
//    |   |  | `+List = EnumFieldList
//    |   |  |  |-EnumFieldDecl
//    |   |  |  | |+Name = Ident
//    |   |  |  | | `+Token = "BLUE"<IDENT@testdata\enum-error.c:2:5>
//    |   |  |  | `+Val = ConstantExpr
//    |   |  |  |  `+X = BasicLit
//    |   |  |  |   `+Token = "1"<INT@testdata\enum-error.c:2:12>
//    |   |  |  |-EnumFieldDecl
//    |   |  |  | |+Name = Ident
//    |   |  |  | | `+Token = "GREEN"<IDENT@testdata\enum-error.c:3:5>
//    |   |  |  | `+Val = <nil>
//    |   |  |  `-EnumFieldDecl
//    |   |  |   |+Name = Ident
//    |   |  |   | `+Token = "YELLOW"<IDENT@testdata\enum-error.c:4:5>
//    |   |  |   `+Val = <nil>
//    |   |  `+Init = <nil>
//    |   |-DeclStmt
//    |   | `-VarDecl
//    |   |  |+Name = Ident
//    |   |  | `+Token = "color1"<IDENT@testdata\enum-error.c:9:9>
//    |   |  |+Type = BuildInType
//    |   |  | |+Lit = 
//    |   |  | | `-"int"<KEYWORD@testdata\enum-error.c:9:5>
//    |   |  | `+Type = int
//    |   |  `+Init = <nil>
//    |   |-DeclStmt
//    |   | `-VarDecl
//    |   |  |+Name = Ident
//    |   |  | `+Token = "color2"<IDENT@testdata\enum-error.c:10:10>
//    |   |  |+Type = BuildInType
//    |   |  | |+Lit = 
//    |   |  | | `-"char"<KEYWORD@testdata\enum-error.c:10:5>
//    |   |  | `+Type = char
//    |   |  `+Init = <nil>
//    |   |-DeclStmt
//    |   | `-VarDecl
//    |   |  |+Name = Ident
//    |   |  | `+Token = "Color"<IDENT@testdata\enum-error.c:11:15>
//    |   |  |+Type = BuildInType
//    |   |  | |+Lit = 
//    |   |  | | `-"long"<KEYWORD@testdata\enum-error.c:11:10>
//    |   |  | `+Type = CUnknownType
//    |   |  `+Init = <nil>
//    |   |-ExprStmt
//    |   | `+Expr = Ident
//    |   |  `+Token = "color2"<IDENT@testdata\enum-error.c:11:21>
//    |   |-DeclStmt
//    |   | `-VarDecl
//    |   |  |+Name = Ident
//    |   |  | `+Token = "Color"<IDENT@testdata\enum-error.c:12:15>
//    |   |  |+Type = BuildInType
//    |   |  | |+Lit = 
//    |   |  | | |-"long"<KEYWORD@testdata\enum-error.c:12:5>
//    |   |  | | `-"enum"<KEYWORD@testdata\enum-error.c:12:10>
//    |   |  | `+Type = int
//    |   |  `+Init = <nil>
//    |   |-ExprStmt
//    |   | `+Expr = Ident
//    |   |  `+Token = "color2"<IDENT@testdata\enum-error.c:12:21>
//    |   |-DeclStmt
//    |   | `-VarDecl
//    |   |  |+Name = Ident
//    |   |  | `+Token = "Color"<IDENT@testdata\enum-error.c:13:21>
//    |   |  |+Type = TypeQualifier
//    |   |  | |+Qualifier = map[const:true]
//    |   |  | `+Inner = BuildInType
//    |   |  |  |+Lit = 
//    |   |  |  | |-"long"<KEYWORD@testdata\enum-error.c:13:11>
//    |   |  |  | `-"enum"<KEYWORD@testdata\enum-error.c:13:16>
//    |   |  |  `+Type = int
//    |   |  `+Init = <nil>
//    |   `-ExprStmt
//    |    `+Expr = Ident
//    |     `+Token = "color3"<IDENT@testdata\enum-error.c:13:27>
//    `+Unresolved = 
// ===========================
//
// |-Error
// | |+Pos = testdata\enum-error.c:9:9
// | |+Typ = 0
// | `+Msg = 在 testdata\enum-error.c 文件的第9行9列: 重复的标识符 color1，上次声明的位置 testdata\enum-error.c:8:13
// |-Error
// | |+Pos = testdata\enum-error.c:11:21
// | |+Typ = 0
// | `+Msg = 在 testdata\enum-error.c 文件的第11行21列: 这里应该是一个 ; ，不应该出现 color2
// |-Error
// | |+Pos = testdata\enum-error.c:12:10
// | |+Typ = 0
// | `+Msg = 在 testdata\enum-error.c 文件的第12行10列: 非预期的类型定义符号 enum
// |-Error
// | |+Pos = testdata\enum-error.c:12:15
// | |+Typ = 0
// | `+Msg = 在 testdata\enum-error.c 文件的第12行15列: 重复的标识符 Color，上次声明的位置 testdata\enum-error.c:11:15
// |-Error
// | |+Pos = testdata\enum-error.c:12:21
// | |+Typ = 0
// | `+Msg = 在 testdata\enum-error.c 文件的第12行21列: 这里应该是一个 ; ，不应该出现 color2
// |-Error
// | |+Pos = testdata\enum-error.c:13:16
// | |+Typ = 0
// | `+Msg = 在 testdata\enum-error.c 文件的第13行16列: 非预期的类型定义符号 enum
// |-Error
// | |+Pos = testdata\enum-error.c:13:21
// | |+Typ = 0
// | `+Msg = 在 testdata\enum-error.c 文件的第13行21列: 重复的标识符 Color，上次声明的位置 testdata\enum-error.c:11:15
// `-Error
//  |+Pos = testdata\enum-error.c:13:27
//  |+Typ = 0
//  `+Msg = 在 testdata\enum-error.c 文件的第13行27列: 这里应该是一个 ; ，不应该出现 color3
// ===========================
