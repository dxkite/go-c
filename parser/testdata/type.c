int i;
const int ci;
const long long long int clli = 10;
const double cd;
const long double cld[] = {1,2,3,4};
const float cf;
const int short is_err;
extern int *ei_v;
static int *const si;

struct tree {
    struct tree *left, right;
} abc;

int main() {
    int a;
    char b;
    return 0;
}

// ===========================
// TranslationUnit
//  `+Files = 
//   `-File
//    |+Name = testdata\type.c
//    |+Decl = 
//    | |-VarDecl
//    | | |+Name = i
//    | | |+Type =  int
//    | | `+Init = <nil>
//    | |-VarDecl
//    | | |+Name = ci
//    | | |+Type = const int
//    | | `+Init = <nil>
//    | |-VarDecl
//    | | |+Name = clli
//    | | |+Type = const long long long int
//    | | `+Init = 10
//    | |-VarDecl
//    | | |+Name = cd
//    | | |+Type = const double
//    | | `+Init = <nil>
//    | |-VarDecl
//    | | |+Name = cld
//    | | |+Type =  const long double[]
//    | | `+Init = InitializerExpr
//    | |  |-1
//    | |  |-2
//    | |  |-3
//    | |  `-4
//    | |-VarDecl
//    | | |+Name = cf
//    | | |+Type = const float
//    | | `+Init = <nil>
//    | |-VarDecl
//    | | |+Name = is_err
//    | | |+Type = const int short
//    | | `+Init = <nil>
//    | |-VarDecl
//    | | |+Name = ei_v
//    | | |+Type =  int *
//    | | `+Init = <nil>
//    | |-VarDecl
//    | | |+Name = si
//    | | |+Type =  int *const
//    | | `+Init = <nil>
//    | |-VarDecl
//    | | |+Name = abc
//    | | |+Type =  struct tree
//    | | `+Init = <nil>
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
//    |   `-ReturnStmt
//    |    `+X = 0
//    `+Unresolved = 
// ===========================
//
// `-Error
//  |+Pos = testdata\type.c:12:12
//  |+Typ = 0
//  `+Msg = 在 testdata\type.c 文件的第12行12列: 不完全的结构体类型 tree
// ===========================
