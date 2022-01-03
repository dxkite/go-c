package ast

import (
	"dxkite.cn/c/token"
	"fmt"
	"strings"
)

type Node interface {
	Pos() token.Position
	End() token.Position
}

type Expr interface {
	expr()
}

type (
	// 无法解析的表达式
	BadExpr struct {
		token.Token
	}

	// 标识符
	Ident struct {
		token.Token
		Type Typename
	}

	// 字面常量
	BasicLit struct {
		// token.INT token.FLOAT token.CHAR token.STRING
		token.Token
	}

	// 类型初始化表达式
	CompoundLit struct {
		Type     Typename
		InitList *InitializerExpr
	}

	InitializerExpr []Expr

	// 数组编号
	ArrayDesignatorExpr struct {
		Index Expr
		X     Expr
	}

	// struct/union
	RecordDesignatorExpr struct {
		Field *Ident
		X     Expr
	}

	// 函数调用
	CallExpr struct {
		Func Expr
		Args []Expr
	}

	// 数组下标表达式
	IndexExpr struct {
		Arr   Expr
		Index Expr
	}

	// 选择表达式 -> .
	SelectorExpr struct {
		X    Expr
		Op   token.Token
		Name *Ident
	}

	TypeCastExpr struct {
		X    Expr
		Type Typename
	}

	ParenExpr struct {
		X Expr
	}

	// "++", "--", "&", "*", "+", "-", "~", "!" sizeof
	UnaryExpr struct {
		Op token.Token
		X  Expr
	}

	BinaryExpr struct {
		X  Expr
		Op token.Token
		Y  Expr
	}

	// 条件表达式 logical-OR-expression ? expression : conditional-expression
	CondExpr struct {
		X    Expr
		Op   token.Token // 操作类型
		Then Expr        // 左值
		Else Expr        // 右值
	}

	// 常量表达式
	ConstantExpr struct {
		X Expr
	}

	// 赋值表达式
	AssignExpr struct {
		X  Expr
		Op token.Token
		Y  Expr
	}

	// 逗号表达式
	CommaExpr []Expr

	// sizeof 表达式
	SizeOfExpr struct {
		Type Typename
	}
)

func (*BadExpr) expr()                  {}
func (e *BadExpr) Pos() token.Position  { return e.Position() }
func (*Ident) expr()                    {}
func (e *Ident) Pos() token.Position    { return e.Position() }
func (*BasicLit) expr()                 {}
func (e *BasicLit) Pos() token.Position { return e.Position() }
func (*CompoundLit) expr()              {}

func (e *Ident) String() string    { return e.Literal() }
func (e *BasicLit) String() string { return e.Literal() }

//func (e *CompoundLit) Pos() token.Position { return e }
func (*InitializerExpr) expr()      {}
func (*RecordDesignatorExpr) expr() {}
func (*ArrayDesignatorExpr) expr()  {}
func (*IndexExpr) expr()            {}
func (*SelectorExpr) expr()         {}
func (*CallExpr) expr()             {}
func (*ParenExpr) expr()            {}
func (*UnaryExpr) expr()            {}
func (*TypeCastExpr) expr()         {}
func (*BinaryExpr) expr()           {}
func (*CondExpr) expr()             {}
func (*ConstantExpr) expr()         {}
func (*AssignExpr) expr()           {}
func (*CommaExpr) expr()            {}
func (*SizeOfExpr) expr()           {}

type Typename interface {
	typeName()
	Qualifier() *Qualifier
	String() string
}

type (
	Qualifier map[string]bool

	// 结构体/联合体
	RecordType struct {
		Qua    *Qualifier
		Type   token.Token
		Name   *Ident
		Fields []*RecordField
	}

	RecordField struct {
		Type Typename
		Name *Ident
		Bit  Expr
	}

	// 枚举类型
	EnumFieldDecl struct {
		Name *Ident
		Val  Expr
	}
	EnumFieldList []*EnumFieldDecl
	EnumType      struct {
		Qua  *Qualifier
		Name *Ident
		List EnumFieldList
	}

	// 内置类型
	BuildInType struct {
		Qua  *Qualifier
		Type BasicType
	}

	Specifier map[string]bool

	// 指针类型
	PointerType struct {
		Qua   *Qualifier
		Inner Typename
	}

	// 数组类型
	ArrayType struct {
		Qua    *Qualifier
		Inner  Typename
		Static bool
		Size   Expr
	}

	// 常量数组
	ConstArrayType struct {
		Qua   *Qualifier
		Inner Typename
		Size  Expr
	}

	// T [*]
	// T []
	IncompleteArrayType struct {
		Qua   *Qualifier
		Inner Typename
	}

	// 函数类型
	ParamVarDecl struct {
		Qua  *StorageSpecifier
		Name *Ident
		Type Typename
	}

	ParamList []*ParamVarDecl
	FuncType  struct {
		Qua      *Qualifier
		Params   ParamList
		Ellipsis bool // ...
		Return   Typename
	}

	// 括号 (abstract)
	ParenType struct {
		Inner Typename
	}
)

func (*RecordType) typeName()               {}
func (t *RecordType) Qualifier() *Qualifier { return t.Qua }

func (*EnumType) typeName() {}
func (t *EnumType) Qualifier() *Qualifier {
	if t.Qua == nil {
		t.Qua = &Qualifier{}
	}
	return t.Qua
}
func (*BuildInType) typeName() {}
func (t *BuildInType) Qualifier() *Qualifier {
	if t.Qua == nil {
		t.Qua = &Qualifier{}
	}
	return t.Qua
}
func (*PointerType) typeName() {}
func (t *PointerType) Qualifier() *Qualifier {
	if t.Qua == nil {
		t.Qua = &Qualifier{}
	}
	return t.Qua
}
func (*ArrayType) typeName() {}
func (t *ArrayType) Qualifier() *Qualifier {
	if t.Qua == nil {
		t.Qua = &Qualifier{}
	}
	return t.Qua
}
func (*IncompleteArrayType) typeName() {}
func (t *IncompleteArrayType) Qualifier() *Qualifier {
	if t.Qua == nil {
		t.Qua = &Qualifier{}
	}
	return t.Qua
}
func (*ConstArrayType) typeName() {}
func (t *ConstArrayType) Qualifier() *Qualifier {
	if t.Qua == nil {
		t.Qua = &Qualifier{}
	}
	return t.Qua
}
func (*FuncType) typeName() {}
func (t *FuncType) Qualifier() *Qualifier {
	if t.Qua == nil {
		t.Qua = &Qualifier{}
	}
	return t.Qua
}
func (*ParenType) typeName()               {}
func (t *ParenType) Qualifier() *Qualifier { return t.Inner.Qualifier() }

func (q *Qualifier) String() string {
	if q == nil {
		return ""
	}
	var qs []string
	for name := range *q {
		qs = append(qs, name)
	}
	return strings.Join(qs, " ")
}

func (t *RecordType) String() string {
	name := ""
	if t.Name != nil {
		name = t.Name.Literal()
	}
	return fmt.Sprintf("%s struct %s", t.Qualifier().String(), name)
}

func (t *EnumType) String() string {
	name := ""
	if t.Name != nil {
		name = t.Name.Literal()
	}
	return fmt.Sprintf("%s enum %s", t.Qualifier().String(), name)
}

func (t *BuildInType) String() string {
	name := []string{t.Qualifier().String()}
	name = append(name, t.Type.String())
	return strings.Join(name, " ")
}

func (t *PointerType) String() string {
	return fmt.Sprintf("%s *%s", t.Inner.String(), t.Qua.String())
}

func (t *ArrayType) String() string {
	return fmt.Sprintf("%s %s[]", t.Qualifier().String(), t.Inner.String())
}
func (t *ConstArrayType) String() string {
	return fmt.Sprintf("%s %s[]", t.Qualifier().String(), t.Inner.String())
}
func (t *IncompleteArrayType) String() string {
	return fmt.Sprintf("%s %s[]", t.Qualifier().String(), t.Inner.String())
}

func (t *FuncType) String() string {
	var params []string
	for _, v := range t.Params {
		params = append(params, v.Type.String())
	}
	if t.Ellipsis {
		params = append(params, "...")
	}
	return fmt.Sprintf("%s (%s)", t.Return.String(), strings.Join(params, ","))
}

func (t *ParenType) String() string {
	return fmt.Sprintf("(%s)", t.Inner.String())
}

type Stmt interface {
	stmt()
}

type (
	// 定义语句
	DeclStmt []Decl

	// typedef /extern /static /auto /register
	StorageSpecifier map[string]bool

	LabelStmt struct {
		Id   *Ident
		Stmt Stmt
	}

	CaseStmt struct {
		Expr Expr
		Stmt Stmt
	}

	DefaultStmt struct {
		Stmt Stmt
	}

	CompoundStmt []Stmt

	ExprStmt struct {
		Expr Expr
	}

	IfStmt struct {
		X    Expr
		Then Stmt
		Else Stmt
	}

	SwitchStmt struct {
		X    Expr
		Stmt Stmt
	}

	WhileStmt struct {
		X    Expr
		Stmt Stmt
	}

	DoWhileStmt struct {
		Stmt Stmt
		X    Expr
	}

	ForStmt struct {
		Init Expr
		Decl Stmt
		Cond Expr
		Post Expr
		Stmt Stmt
	}

	GotoStmt struct {
		Id *Ident
	}

	ContinueStmt struct{}
	BreakStmt    struct{}
	ReturnStmt   struct {
		X Expr
	}
)

func (*DeclStmt) stmt()     {}
func (*LabelStmt) stmt()    {}
func (*CaseStmt) stmt()     {}
func (*DefaultStmt) stmt()  {}
func (*CompoundStmt) stmt() {}
func (*ExprStmt) stmt()     {}
func (*IfStmt) stmt()       {}
func (*SwitchStmt) stmt()   {}
func (*WhileStmt) stmt()    {}
func (*DoWhileStmt) stmt()  {}
func (*ForStmt) stmt()      {}
func (*GotoStmt) stmt()     {}
func (*ContinueStmt) stmt() {}
func (*BreakStmt) stmt()    {}
func (*ReturnStmt) stmt()   {}

type Decl interface {
	decl()
	Ident() *Ident
}

type (
	// 编译单元
	TranslationUnit struct {
		Files []*File
	}

	// 按照文件为单位分割代码
	File struct {
		Name       string
		Decl       []Decl   // 文件内部的定义
		Unresolved []*Ident // 未解析的标识符
	}

	// 函数定义
	FuncDecl struct {
		Name     *Ident
		Params   ParamList
		Ellipsis bool // ...
		Return   Typename
		Decl     []Decl
		Body     *CompoundStmt
	}

	// 变量定义
	VarDecl struct {
		Name *Ident
		Type Typename
		Init Expr
	}

	// 类型定义
	TypedefDecl struct {
		Name *Ident
		Type Typename
	}
)

func (*FuncDecl) decl() {}
func (f *FuncDecl) Ident() *Ident {
	return f.Name
}

func (*VarDecl) decl() {}
func (f *VarDecl) Ident() *Ident {
	return f.Name
}

func (*TypedefDecl) decl() {}
func (f *TypedefDecl) Ident() *Ident {
	return f.Name
}

func (*ParamVarDecl) decl() {}
func (f *ParamVarDecl) Ident() *Ident {
	return f.Name
}

func (*EnumFieldDecl) decl() {}
func (f *EnumFieldDecl) Ident() *Ident {
	return f.Name
}
