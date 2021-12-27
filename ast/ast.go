package ast

import "dxkite.cn/c/token"

type Expr interface {
	expr()
}

type (
	// 无法解析的表达式
	BadExpr []token.Token

	// 标识符
	Ident struct {
		token.Token
	}

	// 字面常量
	BasicLit struct {
		// token.INT token.FLOAT token.CHAR token.STRING
		token.Token
	}

	// 匿名字面量
	CompoundLit struct {
		Type     TypeName
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
		Fun  Expr
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
		Type TypeName
	}

	ParenExpr struct {
		X Expr
	}

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
		Type TypeName
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

type TypeName interface {
	typeName()
}

type (
	BadType struct {
		token.Token
	}

	// 结构体/联合体
	RecordType struct {
		Type   token.Token
		Name   *Ident
		Fields []*RecordField
	}

	RecordField struct {
		Type TypeName
		Name *Ident
		Bit  Expr
	}

	// 用户定义类型
	UserType struct {
		Name *Ident
		Type TypeName
	}

	// 枚举类型
	EnumField struct {
		Name *Ident
		Val  Expr
	}
	EnumFieldList []*EnumField
	EnumType      struct {
		Name *Ident
		List EnumFieldList
	}

	// 内置类型
	BuildInType struct {
		Type []token.Token
	}

	Qualifier map[string]bool
	Specifier map[string]bool

	// const/restrict/volatile
	TypeQualifier struct {
		Qualifier *Qualifier
		Inner     TypeName
	}

	// extern/static/auto/register
	TypeStorageSpecifier struct {
		Specifier *Specifier
		Inner     TypeName
	}

	// 指针类型
	PointerType struct {
		Qualifier *Qualifier
		Inner     TypeName
	}

	// 数组类型
	ArrayType struct {
		Inner     TypeName
		Qualifier *Qualifier
		Static    bool
		Size      Expr
	}

	// 常量数组
	ConstArrayType struct {
		Inner TypeName
		Size  Expr
	}

	// T [*]
	// T []
	IncompleteArrayType struct {
		Inner TypeName
	}

	// 函数类型
	ParamVarDecl struct {
		Type TypeName
		Name *Ident
	}
	ParamList []*ParamVarDecl
	FuncType  struct {
		Inner    TypeName
		Params   ParamList
		Ellipsis bool // ...
	}

	// 括号 (abstract)
	ParenType struct {
		Inner TypeName
	}
)

func (*BadType) typeName()              {}
func (*RecordType) typeName()           {}
func (*UserType) typeName()             {}
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
}

type (
	TranslationUnitDecl struct {
		Decl []Decl
	}

	// 函数定义
	FuncDecl struct {
		Name     *Ident
		Return   TypeName
		Params   ParamList
		Ellipsis bool // ...
		Decl     []Decl
		Body     *CompoundStmt
	}

	// 变量定义
	VarDecl struct {
		Type TypeName
		Name *Ident
		Init Expr
	}

	// 类型定义
	TypedefDecl struct {
		Type TypeName
		Name *Ident
	}
)

func (*FuncDecl) decl()            {}
func (*VarDecl) decl()             {}
func (*TypedefDecl) decl()         {}
func (*TranslationUnitDecl) decl() {}
