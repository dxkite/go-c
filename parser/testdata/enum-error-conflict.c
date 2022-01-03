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
//    | | |+Name = color_t
//    | | `+Type =  enum Color
//    | `-FuncDecl
//    |  |+Name = main
//    |  |+Type =  int ()
//    |  |+Decl = 
//    |  `+Body = CompoundStmt
//    |   |-DeclStmt
//    |   | `-VarDecl
//    |   |  |+Name = color1
//    |   |  |+Type =  enum Color
//    |   |  `+Init = <nil>
//    |   |-DeclStmt
//    |   | `-VarDecl
//    |   |  |+Name = a
//    |   |  |+Type =  int
//    |   |  `+Init = YELLOW
//    |   `-DeclStmt
//    |    `-VarDecl
//    |     |+Name = YELLOW
//    |     |+Type =  int
//    |     `+Init = 10
//    `+Unresolved = 
// ===========================
//
// `-Error
//  |+Pos = testdata\enum-error-conflict.c:5:5
//  |+Typ = 0
//  `+Msg = 在 testdata\enum-error-conflict.c 文件的第5行5列: 重复的标识符 YELLOW，上次声明的位置 testdata\enum-error-conflict.c:4:5
// ===========================
