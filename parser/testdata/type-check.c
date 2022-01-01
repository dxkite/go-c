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
//    | | |+Name = Ident
//    | | | `+Token = "tree"<IDENT@testdata\type-check.c:3:3>
//    | | `+Type = RecordType
//    | |  |+Type = "struct"<KEYWORD@testdata\type-check.c:1:9>
//    | |  |+Name = Ident
//    | |  | `+Token = "tree"<IDENT@testdata\type-check.c:1:16>
//    | |  `+Fields = 
//    | |   |-RecordField
//    | |   | |+Type = PointerType
//    | |   | | `+Inner = RecordType
//    | |   | |  |+Type = "struct"<KEYWORD@testdata\type-check.c:2:5>
//    | |   | |  |+Name = Ident
//    | |   | |  | `+Token = "tree"<IDENT@testdata\type-check.c:2:12>
//    | |   | |  `+Fields = 
//    | |   | |+Name = Ident
//    | |   | | `+Token = "left"<IDENT@testdata\type-check.c:2:18>
//    | |   | `+Bit = <nil>
//    | |   `-RecordField
//    | |    |+Type = RecordType
//    | |    | |+Type = "struct"<KEYWORD@testdata\type-check.c:2:5>
//    | |    | |+Name = Ident
//    | |    | | `+Token = "tree"<IDENT@testdata\type-check.c:2:12>
//    | |    | `+Fields = 
//    | |    |+Name = Ident
//    | |    | `+Token = "right"<IDENT@testdata\type-check.c:2:24>
//    | |    `+Bit = <nil>
//    | `-FuncDecl
//    |  |+Name = Ident
//    |  | `+Token = "main"<IDENT@testdata\type-check.c:5:5>
//    |  |+Params = ParamList
//    |  |+Ellipsis = false
//    |  |+Return = BuildInType
//    |  | |+Lit = 
//    |  | | `-"int"<KEYWORD@testdata\type-check.c:5:1>
//    |  | `+Type = int
//    |  |+Decl = 
//    |  `+Body = CompoundStmt
//    |   |-DeclStmt
//    |   | `-VarDecl
//    |   |  |+Name = Ident
//    |   |  | `+Token = "a"<IDENT@testdata\type-check.c:6:9>
//    |   |  |+Type = BuildInType
//    |   |  | |+Lit = 
//    |   |  | | `-"int"<KEYWORD@testdata\type-check.c:6:5>
//    |   |  | `+Type = int
//    |   |  `+Init = <nil>
//    |   |-DeclStmt
//    |   | `-VarDecl
//    |   |  |+Name = Ident
//    |   |  | `+Token = "b"<IDENT@testdata\type-check.c:7:10>
//    |   |  |+Type = BuildInType
//    |   |  | |+Lit = 
//    |   |  | | `-"char"<KEYWORD@testdata\type-check.c:7:5>
//    |   |  | `+Type = char
//    |   |  `+Init = <nil>
//    |   |-DeclStmt
//    |   | `-VarDecl
//    |   |  |+Name = Ident
//    |   |  | `+Token = "tree"<IDENT@testdata\type-check.c:8:12>
//    |   |  |+Type = PointerType
//    |   |  | `+Inner = RecordType
//    |   |  |  |+Type = "struct"<KEYWORD@testdata\type-check.c:1:9>
//    |   |  |  |+Name = Ident
//    |   |  |  | `+Token = "tree"<IDENT@testdata\type-check.c:1:16>
//    |   |  |  `+Fields = 
//    |   |  |   |-RecordField
//    |   |  |   | |+Type = PointerType
//    |   |  |   | | `+Inner = RecordType
//    |   |  |   | |  |+Type = "struct"<KEYWORD@testdata\type-check.c:2:5>
//    |   |  |   | |  |+Name = Ident
//    |   |  |   | |  | `+Token = "tree"<IDENT@testdata\type-check.c:2:12>
//    |   |  |   | |  `+Fields = 
//    |   |  |   | |+Name = Ident
//    |   |  |   | | `+Token = "left"<IDENT@testdata\type-check.c:2:18>
//    |   |  |   | `+Bit = <nil>
//    |   |  |   `-RecordField
//    |   |  |    |+Type = RecordType
//    |   |  |    | |+Type = "struct"<KEYWORD@testdata\type-check.c:2:5>
//    |   |  |    | |+Name = Ident
//    |   |  |    | | `+Token = "tree"<IDENT@testdata\type-check.c:2:12>
//    |   |  |    | `+Fields = 
//    |   |  |    |+Name = Ident
//    |   |  |    | `+Token = "right"<IDENT@testdata\type-check.c:2:24>
//    |   |  |    `+Bit = <nil>
//    |   |  `+Init = <nil>
//    |   `-DeclStmt
//    |    `-VarDecl
//    |     |+Name = Ident
//    |     | `+Token = "tree"<IDENT@testdata\type-check.c:9:18>
//    |     |+Type = PointerType
//    |     | `+Inner = RecordType
//    |     |  |+Type = "struct"<KEYWORD@testdata\type-check.c:9:5>
//    |     |  |+Name = Ident
//    |     |  | `+Token = "tree"<IDENT@testdata\type-check.c:9:12>
//    |     |  `+Fields = 
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
//  `+Msg = 在 testdata\type-check.c 文件的第9行18列: 重复的标识符 tree，上次声明的位置 testdata\type-check.c:8:12
// ===========================
