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
//    | | |+Typedef = testdata\enum-error-conflict.c:1:1
//    | | |+Type =  enum Color
//    | | `+Name = color_t
//    | `-FuncDecl
//    |  |+Name = main
//    |  |+Type =  int ()
//    |  |+Decl = 
//    |  `+Body = CompoundStmt
//    |   |+Lbrace = testdata\enum-error-conflict.c:8:12
//    |   |+Stmts = 
//    |   | |-DeclStmt
//    |   | | `-VarDecl
//    |   | |  |+Type =  enum Color
//    |   | |  |+Name = color1
//    |   | |  `+Init = <nil>
//    |   | |-DeclStmt
//    |   | | `-VarDecl
//    |   | |  |+Type =  int
//    |   | |  |+Name = a
//    |   | |  `+Init = YELLOW
//    |   | `-DeclStmt
//    |   |  `-VarDecl
//    |   |   |+Type =  int
//    |   |   |+Name = YELLOW
//    |   |   `+Init = 10
//    |   `+Rbrace = testdata\enum-error-conflict.c:12:1
//    `+Unresolved = 
// ===========================
//
// `-Error
//  |+Pos = testdata\enum-error-conflict.c:5:5
//  |+Typ = 0
//  `+Msg = 在 testdata\enum-error-conflict.c 文件的第5行5列: 重复的标识符 YELLOW，上次声明的位置 testdata\enum-error-conflict.c:4:5
// ===========================
