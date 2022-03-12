int i;
const int ci;
const long long long int clli = 10;
const double cd;
const long double cld[] = {1,2,3,4};
const float cf;
const int short is_err;
extern int *ei_v;
static int *const si;
int (a*)();

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
//    | | |+Type =  int
//    | | |+Name = i
//    | | `+Init = <nil>
//    | |-VarDecl
//    | | |+Type = const int
//    | | |+Name = ci
//    | | `+Init = <nil>
//    | |-VarDecl
//    | | |+Type = const int
//    | | |+Name = clli
//    | | `+Init = 10
//    | |-VarDecl
//    | | |+Type = const double
//    | | |+Name = cd
//    | | `+Init = <nil>
//    | |-VarDecl
//    | | |+Type =  const double[]
//    | | |+Name = cld
//    | | `+Init = InitializerExpr
//    | |  |+Lbrace = testdata\type.c:5:27
//    | |  |+List = 
//    | |  | |-1
//    | |  | |-2
//    | |  | |-3
//    | |  | `-4
//    | |  `+Rbrace = testdata\type.c:5:35
//    | |-VarDecl
//    | | |+Type = const float
//    | | |+Name = cf
//    | | `+Init = <nil>
//    | |-VarDecl
//    | | |+Type = const int
//    | | |+Name = is_err
//    | | `+Init = <nil>
//    | |-VarDecl
//    | | |+Type =  int *
//    | | |+Name = ei_v
//    | | `+Init = <nil>
//    | |-VarDecl
//    | | |+Type =  int *const
//    | | |+Name = si
//    | | `+Init = <nil>
//    | |-VarDecl
//    | | |+Type = ( int)
//    | | |+Name = a
//    | | `+Init = <nil>
//    | |-VarDecl
//    | | |+Type =  struct tree
//    | | |+Name = abc
//    | | `+Init = <nil>
//    | `-FuncDecl
//    |  |+Name = main
//    |  |+Type =  int ()
//    |  |+Decl = 
//    |  `+Body = CompoundStmt
//    |   |+Lbrace = testdata\type.c:16:12
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
//    |   | `-ReturnStmt
//    |   |  |+Return = testdata\type.c:19:5
//    |   |  |+X = 0
//    |   |  `+Semicolon = testdata\type.c:19:13
//    |   `+Rbrace = testdata\type.c:20:1
//    `+Unresolved = 
// ===========================
//
// |-Error
// | |+Pos = testdata\type.c:3:17
// | |+Typ = 0
// | `+Msg = 在 testdata\type.c 文件的第3行17列: 非预期的类型定义符号 long
// |-Error
// | |+Pos = testdata\type.c:7:11
// | |+Typ = 0
// | `+Msg = 在 testdata\type.c 文件的第7行11列: 非预期的类型定义符号 short
// |-Error
// | |+Pos = testdata\type.c:10:7
// | |+Typ = 0
// | `+Msg = 在 testdata\type.c 文件的第10行7列: 这里应该是一个 ) ，不应该出现 *
// |-Error
// | |+Pos = testdata\type.c:10:7
// | |+Typ = 0
// | `+Msg = 在 testdata\type.c 文件的第10行7列: 这里应该是一个 ; ，不应该出现 *
// |-Error
// | |+Pos = testdata\type.c:10:7
// | |+Typ = 0
// | `+Msg = 在 testdata\type.c 文件的第10行7列: 这里应该是一个 定义语句 ，不应该出现 *
// |-Error
// | |+Pos = testdata\type.c:10:8
// | |+Typ = 0
// | `+Msg = 在 testdata\type.c 文件的第10行8列: 这里应该是一个 定义语句 ，不应该出现 )
// |-Error
// | |+Pos = testdata\type.c:10:9
// | |+Typ = 0
// | `+Msg = 在 testdata\type.c 文件的第10行9列: 这里应该是一个 定义语句 ，不应该出现 (
// |-Error
// | |+Pos = testdata\type.c:10:10
// | |+Typ = 0
// | `+Msg = 在 testdata\type.c 文件的第10行10列: 这里应该是一个 定义语句 ，不应该出现 )
// |-Error
// | |+Pos = testdata\type.c:10:11
// | |+Typ = 0
// | `+Msg = 在 testdata\type.c 文件的第10行11列: 这里应该是一个 定义语句 ，不应该出现 ;
// `-Error
//  |+Pos = testdata\type.c:13:12
//  |+Typ = 0
//  `+Msg = 在 testdata\type.c 文件的第13行12列: 不完全的结构体类型 tree
// ===========================
