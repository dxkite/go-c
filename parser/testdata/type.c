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
//    | | |+Name = Ident
//    | | | `+Token = "i"<IDENT@testdata\type.c:1:5>
//    | | |+Type = BuildInType
//    | | | `+Lit = 
//    | | |  `-"int"<KEYWORD@testdata\type.c:1:1>
//    | | `+Init = <nil>
//    | |-VarDecl
//    | | |+Name = Ident
//    | | | `+Token = "ci"<IDENT@testdata\type.c:2:11>
//    | | |+Type = TypeQualifier
//    | | | |+Qualifier = map[const:true]
//    | | | `+Inner = BuildInType
//    | | |  `+Lit = 
//    | | |   `-"int"<KEYWORD@testdata\type.c:2:7>
//    | | `+Init = <nil>
//    | |-VarDecl
//    | | |+Name = Ident
//    | | | `+Token = "clli"<IDENT@testdata\type.c:3:26>
//    | | |+Type = TypeQualifier
//    | | | |+Qualifier = map[const:true]
//    | | | `+Inner = BuildInType
//    | | |  `+Lit = 
//    | | |   |-"long"<KEYWORD@testdata\type.c:3:7>
//    | | |   |-"long"<KEYWORD@testdata\type.c:3:12>
//    | | |   |-"long"<KEYWORD@testdata\type.c:3:17>
//    | | |   `-"int"<KEYWORD@testdata\type.c:3:22>
//    | | `+Init = BasicLit
//    | |  `+Token = "10"<INT@testdata\type.c:3:33>
//    | |-VarDecl
//    | | |+Name = Ident
//    | | | `+Token = "cd"<IDENT@testdata\type.c:4:14>
//    | | |+Type = TypeQualifier
//    | | | |+Qualifier = map[const:true]
//    | | | `+Inner = BuildInType
//    | | |  `+Lit = 
//    | | |   `-"double"<KEYWORD@testdata\type.c:4:7>
//    | | `+Init = <nil>
//    | |-VarDecl
//    | | |+Name = Ident
//    | | | `+Token = "cld"<IDENT@testdata\type.c:5:19>
//    | | |+Type = IncompleteArrayType
//    | | | `+Inner = TypeQualifier
//    | | |  |+Qualifier = map[const:true]
//    | | |  `+Inner = BuildInType
//    | | |   `+Lit = 
//    | | |    |-"long"<KEYWORD@testdata\type.c:5:7>
//    | | |    `-"double"<KEYWORD@testdata\type.c:5:12>
//    | | `+Init = InitializerExpr
//    | |  |-BasicLit
//    | |  | `+Token = "1"<INT@testdata\type.c:5:28>
//    | |  |-BasicLit
//    | |  | `+Token = "2"<INT@testdata\type.c:5:30>
//    | |  |-BasicLit
//    | |  | `+Token = "3"<INT@testdata\type.c:5:32>
//    | |  `-BasicLit
//    | |   `+Token = "4"<INT@testdata\type.c:5:34>
//    | |-VarDecl
//    | | |+Name = Ident
//    | | | `+Token = "cf"<IDENT@testdata\type.c:6:13>
//    | | |+Type = TypeQualifier
//    | | | |+Qualifier = map[const:true]
//    | | | `+Inner = BuildInType
//    | | |  `+Lit = 
//    | | |   `-"float"<KEYWORD@testdata\type.c:6:7>
//    | | `+Init = <nil>
//    | |-VarDecl
//    | | |+Name = Ident
//    | | | `+Token = "is_err"<IDENT@testdata\type.c:7:17>
//    | | |+Type = TypeQualifier
//    | | | |+Qualifier = map[const:true]
//    | | | `+Inner = BuildInType
//    | | |  `+Lit = 
//    | | |   |-"int"<KEYWORD@testdata\type.c:7:7>
//    | | |   `-"short"<KEYWORD@testdata\type.c:7:11>
//    | | `+Init = <nil>
//    | |-VarDecl
//    | | |+Name = Ident
//    | | | `+Token = "ei_v"<IDENT@testdata\type.c:8:13>
//    | | |+Type = PointerType
//    | | | `+Inner = TypeStorageSpecifier
//    | | |  |+Specifier = map[extern:true]
//    | | |  `+Inner = BuildInType
//    | | |   `+Lit = 
//    | | |    `-"int"<KEYWORD@testdata\type.c:8:8>
//    | | `+Init = <nil>
//    | |-VarDecl
//    | | |+Name = Ident
//    | | | `+Token = "si"<IDENT@testdata\type.c:9:19>
//    | | |+Type = TypeQualifier
//    | | | |+Qualifier = map[const:true]
//    | | | `+Inner = PointerType
//    | | |  `+Inner = TypeStorageSpecifier
//    | | |   |+Specifier = map[static:true]
//    | | |   `+Inner = BuildInType
//    | | |    `+Lit = 
//    | | |     `-"int"<KEYWORD@testdata\type.c:9:8>
//    | | `+Init = <nil>
//    | |-VarDecl
//    | | |+Name = Ident
//    | | | `+Token = "abc"<IDENT@testdata\type.c:13:3>
//    | | |+Type = RecordType
//    | | | |+Type = "struct"<KEYWORD@testdata\type.c:11:1>
//    | | | |+Name = Ident
//    | | | | `+Token = "tree"<IDENT@testdata\type.c:11:8>
//    | | | `+Fields = 
//    | | |  |-RecordField
//    | | |  | |+Type = PointerType
//    | | |  | | `+Inner = RecordType
//    | | |  | |  |+Type = "struct"<KEYWORD@testdata\type.c:12:5>
//    | | |  | |  |+Name = Ident
//    | | |  | |  | `+Token = "tree"<IDENT@testdata\type.c:12:12>
//    | | |  | |  `+Fields = 
//    | | |  | |+Name = Ident
//    | | |  | | `+Token = "left"<IDENT@testdata\type.c:12:18>
//    | | |  | `+Bit = <nil>
//    | | |  `-RecordField
//    | | |   |+Type = RecordType
//    | | |   | |+Type = "struct"<KEYWORD@testdata\type.c:12:5>
//    | | |   | |+Name = Ident
//    | | |   | | `+Token = "tree"<IDENT@testdata\type.c:12:12>
//    | | |   | `+Fields = 
//    | | |   |+Name = Ident
//    | | |   | `+Token = "right"<IDENT@testdata\type.c:12:24>
//    | | |   `+Bit = <nil>
//    | | `+Init = <nil>
//    | `-FuncDecl
//    |  |+Name = Ident
//    |  | `+Token = "main"<IDENT@testdata\type.c:15:5>
//    |  |+Params = ParamList
//    |  |+Ellipsis = false
//    |  |+Return = BuildInType
//    |  | `+Lit = 
//    |  |  `-"int"<KEYWORD@testdata\type.c:15:1>
//    |  |+Decl = 
//    |  `+Body = CompoundStmt
//    |   |-DeclStmt
//    |   | `-VarDecl
//    |   |  |+Name = Ident
//    |   |  | `+Token = "a"<IDENT@testdata\type.c:16:9>
//    |   |  |+Type = BuildInType
//    |   |  | `+Lit = 
//    |   |  |  `-"int"<KEYWORD@testdata\type.c:16:5>
//    |   |  `+Init = <nil>
//    |   |-DeclStmt
//    |   | `-VarDecl
//    |   |  |+Name = Ident
//    |   |  | `+Token = "b"<IDENT@testdata\type.c:17:10>
//    |   |  |+Type = BuildInType
//    |   |  | `+Lit = 
//    |   |  |  `-"char"<KEYWORD@testdata\type.c:17:5>
//    |   |  `+Init = <nil>
//    |   `-ReturnStmt
//    |    `+X = BasicLit
//    |     `+Token = "0"<INT@testdata\type.c:18:12>
//    `+Unresolved = 
// ===========================
//
// `-Error
//  |+Pos = testdata\type.c:12:12
//  |+Typ = 0
//  `+Msg = 在 testdata\type.c 文件的第12行12列: 不完全的结构体类型 tree
// ===========================
