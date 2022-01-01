package ast

import (
	"dxkite.cn/c/token"
)

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

func (*BadExpr) expr()              {}
func (*Ident) expr()                {}
func (*BasicLit) expr()             {}
func (*CompoundLit) expr()          {}
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
}

type (
	// 结构体/联合体
	RecordType struct {
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
		Name *Ident
		List EnumFieldList
	}

	// 内置类型
	BuildInType struct {
		Lit  []token.Token
		Type CBuildInType
	}

	Qualifier map[string]bool
	Specifier map[string]bool

	// const/restrict/volatile
	TypeQualifier struct {
		Qualifier *Qualifier
		Inner     Typename
	}

	// extern/static/auto/register
	TypeStorageSpecifier struct {
		Specifier *Specifier
		Inner     Typename
	}

	// 指针类型
	PointerType struct {
		Inner Typename
	}

	// 数组类型
	ArrayType struct {
		Inner  Typename
		Static bool
		Size   Expr
	}

	// 常量数组
	ConstArrayType struct {
		Inner Typename
		Size  Expr
	}

	// T [*]
	// T []
	IncompleteArrayType struct {
		Inner Typename
	}

	// 函数类型
	ParamVarDecl struct {
		Name *Ident
		Type Typename
	}
	ParamList []*ParamVarDecl
	FuncType  struct {
		Params   ParamList
		Ellipsis bool // ...
		Return   Typename
	}

	// 括号 (abstract)
	ParenType struct {
		Inner Typename
	}
)

func (*RecordType) typeName()           {}
func (*EnumType) typeName()             {}
func (*BuildInType) typeName()          {}
func (*TypeQualifier) typeName()        {}
func (*TypeStorageSpecifier) typeName() {}
func (*PointerType) typeName()          {}
func (*ArrayType) typeName()            {}
func (*IncompleteArrayType) typeName()  {}
func (*ConstArrayType) typeName()       {}
func (*FuncType) typeName()             {}
func (*ParenType) typeName()            {}

type Stmt interface {
	stmt()
}

type (
	// 定义语句
	DeclStmt []Decl

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
