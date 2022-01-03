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
//    |   `-DeclStmt
//    |    `-VarDecl
//    |     |+Name = color2
//    |     |+Type =  enum Color
//    |     `+Init = <nil>
//    `+Unresolved = 
// ===========================
//
// ===========================
