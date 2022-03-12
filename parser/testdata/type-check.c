typedef struct tree {
    struct tree *left, right;
} tree;

int main() {
    int a;
    char b;
    tree * tree;
    struct tree *tree;
}

// ===========================
// TranslationUnit
//  `+Files = 
//   `-File
//    |+Name = testdata\type-check.c
//    |+Decl = 
//    | |-TypedefDecl
//    | | |+Typedef = testdata\type-check.c:1:1
//    | | |+Type =  struct tree
//    | | `+Name = tree
//    | `-FuncDecl
//    |  |+Name = main
//    |  |+Type =  int ()
//    |  |+Decl = 
//    |  `+Body = CompoundStmt
//    |   |+Lbrace = testdata\type-check.c:5:12
//    |   |+Stmts = 
//    |   | |-DeclStmt
//    |   | | `-VarDecl
//    |   | |  |+Type =  int
//    |   | |  |+Name = a
//    |   | |  `+Init = <nil>
//    |   | |-DeclStmt
//    |   | | `-VarDecl
//    |   | |  |+Type =  char
//    |   | |  |+Name = b
//    |   | |  `+Init = <nil>
//    |   | |-DeclStmt
//    |   | | `-VarDecl
//    |   | |  |+Type =  struct tree *
//    |   | |  |+Name = tree
//    |   | |  `+Init = <nil>
//    |   | `-DeclStmt
//    |   |  `-VarDecl
//    |   |   |+Type =  struct tree *
//    |   |   |+Name = tree
//    |   |   `+Init = <nil>
//    |   `+Rbrace = testdata\type-check.c:10:1
//    `+Unresolved = 
// ===========================
//
// |-Error
// | |+Pos = testdata\type-check.c:2:12
// | |+Typ = 0
// | `+Msg = 在 testdata\type-check.c 文件的第2行12列: 不完全的结构体类型 tree
// `-Error
//  |+Pos = testdata\type-check.c:9:18
//  |+Typ = 0
//  `+Msg = 在 testdata\type-check.c 文件的第9行18列: 重复声明的变量名 tree，上次声明的位置 testdata\type-check.c:8:12
// ===========================
