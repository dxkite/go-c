int max(int a, int b) {
    if (a>b) {
        return a;
    }
    return b;
}

int main() {
    return max(10,20);
}

// ===========================
// TranslationUnit
//  `+Files = 
//   `-File
//    |+Name = testdata\fundef.c
//    |+Decl = 
//    | |-FuncDecl
//    | | |+Name = max
//    | | |+Params = ParamList
//    | | | |-ParamVarDecl
//    | | | | |+Qua = map[]
//    | | | | |+Name = a
//    | | | | `+Type =  int
//    | | | `-ParamVarDecl
//    | | |  |+Qua = map[]
//    | | |  |+Name = b
//    | | |  `+Type =  int
//    | | |+Ellipsis = false
//    | | |+Return =  int
//    | | |+Decl = 
//    | | `+Body = CompoundStmt
//    | |  |-IfStmt
//    | |  | |+X = BinaryExpr
//    | |  | | |+X = a
//    | |  | | |+Op = ">"<PUNCTUATOR@testdata\fundef.c:2:10>
//    | |  | | `+Y = b
//    | |  | |+Then = CompoundStmt
//    | |  | | `-ReturnStmt
//    | |  | |  `+X = a
//    | |  | `+Else = <nil>
//    | |  `-ReturnStmt
//    | |   `+X = b
//    | `-FuncDecl
//    |  |+Name = main
//    |  |+Params = ParamList
//    |  |+Ellipsis = false
//    |  |+Return =  int
//    |  |+Decl = 
//    |  `+Body = CompoundStmt
//    |   `-ReturnStmt
//    |    `+X = CallExpr
//    |     |+Func = max
//    |     `+Args = 
//    |      |-10
//    |      `-20
//    `+Unresolved = 
// ===========================
//
// ===========================
