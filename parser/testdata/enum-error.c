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
//    | | |+Name = color_t
//    | | `+Type =  enum Color
//    | `-FuncDecl
//    |  |+Name = main
//    |  |+Params = ParamList
//    |  |+Ellipsis = false
//    |  |+Return =  int
//    |  |+Decl = 
//    |  `+Body = CompoundStmt
//    |   |-DeclStmt
//    |   | `-VarDecl
//    |   |  |+Name = color1
//    |   |  |+Type =  enum Color
//    |   |  `+Init = <nil>
//    |   |-DeclStmt
//    |   | `-VarDecl
//    |   |  |+Name = color1
//    |   |  |+Type =  int
//    |   |  `+Init = <nil>
//    |   |-DeclStmt
//    |   | `-VarDecl
//    |   |  |+Name = color2
//    |   |  |+Type =  char
//    |   |  `+Init = <nil>
//    |   |-DeclStmt
//    |   | `-VarDecl
//    |   |  |+Name = Color
//    |   |  |+Type =  long
//    |   |  `+Init = <nil>
//    |   |-ExprStmt
//    |   | `+Expr = color2
//    |   |-DeclStmt
//    |   | `-VarDecl
//    |   |  |+Name = Color
//    |   |  |+Type =  long enum
//    |   |  `+Init = <nil>
//    |   |-ExprStmt
//    |   | `+Expr = color2
//    |   |-DeclStmt
//    |   | `-VarDecl
//    |   |  |+Name = Color
//    |   |  |+Type = const long enum
//    |   |  `+Init = <nil>
//    |   `-ExprStmt
//    |    `+Expr = color3
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
// | |+Pos = testdata\enum-error.c:12:15
// | |+Typ = 0
// | `+Msg = 在 testdata\enum-error.c 文件的第12行15列: 重复声明的变量名 Color，上次声明的位置 testdata\enum-error.c:11:15
// |-Error
// | |+Pos = testdata\enum-error.c:12:21
// | |+Typ = 0
// | `+Msg = 在 testdata\enum-error.c 文件的第12行21列: 这里应该是一个 ; ，不应该出现 color2
// |-Error
// | |+Pos = testdata\enum-error.c:13:21
// | |+Typ = 0
// | `+Msg = 在 testdata\enum-error.c 文件的第13行21列: 重复声明的变量名 Color，上次声明的位置 testdata\enum-error.c:11:15
// `-Error
//  |+Pos = testdata\enum-error.c:13:27
//  |+Typ = 0
//  `+Msg = 在 testdata\enum-error.c 文件的第13行27列: 这里应该是一个 ; ，不应该出现 color3
// ===========================
