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

type Declarable interface {
	declarable()
}

type TypeName interface {
	Declarable
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

	// type/extern/static/auto/register
	TypeSpecifier struct {
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
		Inner    Declarable
		Params   ParamList
		Ellipsis bool // ...
	}

	// 括号 (abstract)
	ParenType struct {
		Inner Declarable
	}
)

func (*BadType) typeName()               {}
func (*BadType) declarable()             {}
func (*RecordType) typeName()            {}
func (*RecordType) declarable()          {}
func (*UserType) typeName()              {}
func (*UserType) declarable()            {}
func (*EnumType) typeName()              {}
func (*EnumType) declarable()            {}
func (*BuildInType) typeName()           {}
func (*BuildInType) declarable()         {}
func (*TypeQualifier) typeName()         {}
func (*TypeQualifier) declarable()       {}
func (*TypeSpecifier) typeName()         {}
func (*TypeSpecifier) declarable()       {}
func (*PointerType) typeName()           {}
func (*PointerType) declarable()         {}
func (*ArrayType) typeName()             {}
func (*ArrayType) declarable()           {}
func (*IncompleteArrayType) typeName()   {}
func (*IncompleteArrayType) declarable() {}
func (*ConstArrayType) typeName()        {}
func (*ConstArrayType) declarable()      {}
func (*FuncType) typeName()              {}
func (*FuncType) declarable()            {}
func (*ParenType) typeName()             {}
func (*ParenType) declarable()           {}
func (*ParamVarDecl) declarable()        {}
