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
//    | | |+Typedef = testdata\enum.c:1:1
//    | | |+Type =  enum Color
//    | | `+Name = color_t
//    | `-FuncDecl
//    |  |+Name = main
//    |  |+Type =  int ()
//    |  |+Decl = 
//    |  `+Body = CompoundStmt
//    |   |+Lbrace = testdata\enum.c:7:12
//    |   |+Stmts = 
//    |   | |-DeclStmt
//    |   | | `-VarDecl
//    |   | |  |+Type =  enum Color
//    |   | |  |+Name = color1
//    |   | |  `+Init = <nil>
//    |   | `-DeclStmt
//    |   |  `-VarDecl
//    |   |   |+Type =  enum Color
//    |   |   |+Name = color2
//    |   |   `+Init = <nil>
//    |   `+Rbrace = testdata\enum.c:10:1
//    `+Unresolved = 
// ===========================
//
// ===========================
