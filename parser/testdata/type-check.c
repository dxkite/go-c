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
//    | | |+Name = tree
//    | | `+Type =  struct tree
//    | `-FuncDecl
//    |  |+Name = main
//    |  |+Params = ParamList
//    |  |+Ellipsis = false
//    |  |+Return =  int
//    |  |+Decl = 
//    |  `+Body = CompoundStmt
//    |   |-DeclStmt
//    |   | `-VarDecl
//    |   |  |+Name = a
//    |   |  |+Type =  int
//    |   |  `+Init = <nil>
//    |   |-DeclStmt
//    |   | `-VarDecl
//    |   |  |+Name = b
//    |   |  |+Type =  char
//    |   |  `+Init = <nil>
//    |   |-DeclStmt
//    |   | `-VarDecl
//    |   |  |+Name = tree
//    |   |  |+Type =  struct tree *
//    |   |  `+Init = <nil>
//    |   `-DeclStmt
//    |    `-VarDecl
//    |     |+Name = tree
//    |     |+Type =  struct tree *
//    |     `+Init = <nil>
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
