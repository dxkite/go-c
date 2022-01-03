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
//    |+Name = testdata\type.cc
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
//    | | |+Type = const int
//    | | `+Init = 10
//    | |-VarDecl
//    | | |+Name = cd
//    | | |+Type = const double
//    | | `+Init = <nil>
//    | |-VarDecl
//    | | |+Name = cld
//    | | |+Type =  const double[]
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
//    | | |+Type = const int
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
//    | | |+Name = a
//    | | |+Type = ( int)
//    | | `+Init = <nil>
//    | |-VarDecl
//    | | |+Name = abc
//    | | |+Type =  struct tree
//    | | `+Init = <nil>
//    | `-FuncDecl
//    |  |+Name = main
//    |  |+Type =  int ()
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
// |-Error
// | |+Pos = testdata\type.cc:3:17
// | |+Typ = 0
// | `+Msg = 在 testdata\type.cc 文件的第3行17列: 非预期的类型定义符号 long
// |-Error
// | |+Pos = testdata\type.cc:7:11
// | |+Typ = 0
// | `+Msg = 在 testdata\type.cc 文件的第7行11列: 非预期的类型定义符号 short
// |-Error
// | |+Pos = testdata\type.cc:10:7
// | |+Typ = 0
// | `+Msg = 在 testdata\type.cc 文件的第10行7列: 这里应该是一个 ) ，不应该出现 *
// |-Error
// | |+Pos = testdata\type.cc:10:7
// | |+Typ = 0
// | `+Msg = 在 testdata\type.cc 文件的第10行7列: 这里应该是一个 ; ，不应该出现 *
// |-Error
// | |+Pos = testdata\type.cc:10:7
// | |+Typ = 0
// | `+Msg = 在 testdata\type.cc 文件的第10行7列: 这里应该是一个 定义语句 ，不应该出现 *
// |-Error
// | |+Pos = testdata\type.cc:10:8
// | |+Typ = 0
// | `+Msg = 在 testdata\type.cc 文件的第10行8列: 这里应该是一个 定义语句 ，不应该出现 )
// |-Error
// | |+Pos = testdata\type.cc:10:9
// | |+Typ = 0
// | `+Msg = 在 testdata\type.cc 文件的第10行9列: 这里应该是一个 定义语句 ，不应该出现 (
// |-Error
// | |+Pos = testdata\type.cc:10:10
// | |+Typ = 0
// | `+Msg = 在 testdata\type.cc 文件的第10行10列: 这里应该是一个 定义语句 ，不应该出现 )
// |-Error
// | |+Pos = testdata\type.cc:10:11
// | |+Typ = 0
// | `+Msg = 在 testdata\type.cc 文件的第10行11列: 这里应该是一个 定义语句 ，不应该出现 ;
// `-Error
//  |+Pos = testdata\type.cc:13:12
//  |+Typ = 0
//  `+Msg = 在 testdata\type.cc 文件的第13行12列: 不完全的结构体类型 tree
// ===========================
