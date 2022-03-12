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
//    | | |+Typedef = testdata\enum-error.c:1:1
//    | | |+Type =  enum Color
//    | | `+Name = color_t
//    | `-FuncDecl
//    |  |+Name = main
//    |  |+Type =  int ()
//    |  |+Decl = 
//    |  `+Body = CompoundStmt
//    |   |+Lbrace = testdata\enum-error.c:7:12
//    |   |+Stmts = 
//    |   | |-DeclStmt
//    |   | | `-VarDecl
//    |   | |  |+Type =  enum Color
//    |   | |  |+Name = color1
//    |   | |  `+Init = <nil>
//    |   | |-DeclStmt
//    |   | | `-VarDecl
//    |   | |  |+Type =  int
//    |   | |  |+Name = color1
//    |   | |  `+Init = <nil>
//    |   | |-DeclStmt
//    |   | | `-VarDecl
//    |   | |  |+Type =  char
//    |   | |  |+Name = color2
//    |   | |  `+Init = <nil>
//    |   | |-DeclStmt
//    |   | | `-VarDecl
//    |   | |  |+Type =  UnknownType
//    |   | |  |+Name = Color
//    |   | |  `+Init = <nil>
//    |   | |-ExprStmt
//    |   | | |+Expr = color2
//    |   | | `+Semicolon = testdata\enum-error.c:11:27
//    |   | |-DeclStmt
//    |   | | `-VarDecl
//    |   | |  |+Type =  int
//    |   | |  |+Name = Color
//    |   | |  `+Init = <nil>
//    |   | |-ExprStmt
//    |   | | |+Expr = color2
//    |   | | `+Semicolon = testdata\enum-error.c:12:27
//    |   | |-DeclStmt
//    |   | | `-VarDecl
//    |   | |  |+Type = const int
//    |   | |  |+Name = Color
//    |   | |  `+Init = <nil>
//    |   | `-ExprStmt
//    |   |  |+Expr = color3
//    |   |  `+Semicolon = testdata\enum-error.c:13:33
//    |   `+Rbrace = testdata\enum-error.c:14:1
//    `+Unresolved = 
// ===========================
//
// |-Error
// | |+Pos = testdata\enum-error.c:9:9
// | |+Typ = 0
// | `+Msg = 在 testdata\enum-error.c 文件的第9行9列: 重复声明的变量名 color1，上次声明的位置 testdata\enum-error.c:8:13
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
// | `+Msg = 在 testdata\enum-error.c 文件的第12行15列: 重复声明的变量名 Color，上次声明的位置 testdata\enum-error.c:11:15
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
// | `+Msg = 在 testdata\enum-error.c 文件的第13行21列: 重复声明的变量名 Color，上次声明的位置 testdata\enum-error.c:11:15
// |-Error
// | |+Pos = testdata\enum-error.c:13:27
// | |+Typ = 0
// | `+Msg = 在 testdata\enum-error.c 文件的第13行27列: 这里应该是一个 ; ，不应该出现 color3
// `-Error
//  |+Pos = testdata\enum-error.c:13:27
//  |+Typ = 0
//  `+Msg = 在 testdata\enum-error.c 文件的第13行27列: 未定义的标识符 color3
// ===========================
