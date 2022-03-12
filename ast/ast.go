package ast

import (
	"dxkite.cn/c/token"
	"fmt"
	"strings"
	"unicode/utf8"
)

type Node interface {
	Beg() token.Position
	End() token.Position
}

type Expr interface {
	Node
	expr()
}

type Typename interface {
	Node
	typeName()
	Qualifier() *Qualifier
	String() string
}

type Stmt interface {
	Node
	stmt()
}

type Decl interface {
	Node
	decl()
	Ident() *Ident
}

type Range struct {
	Begin token.Position
	End   token.Position
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
	// (typename) { InitializerExpr }
	CompoundLit struct {
		Lparen token.Position // (
		Type   Typename
		Rparen token.Position // )

		InitList *InitializerExpr
	}

	InitializerExpr struct {
		Lbrace token.Position // {
		List   []Expr
		Rbrace token.Position // }
	}

	// 数组编号
	ArrayDesignatorExpr struct {
		Lbrack token.Position // [
		Index  Expr
		Rbrack token.Position // ]
		X      Expr
	}

	// struct/union
	RecordDesignatorExpr struct {
		Dot   token.Position // .
		Field *Ident
		X     Expr
	}

	// 函数调用
	CallExpr struct {
		Func   Expr
		Lparen token.Position // (
		Args   []Expr
		Rparen token.Position // )
	}

	// 数组下标表达式
	IndexExpr struct {
		Arr    Expr
		Lbrack token.Position // [
		Index  Expr
		Rbrack token.Position // ]
	}

	// 选择表达式 -> .
	SelectorExpr struct {
		X    Expr
		Op   token.Token
		Name *Ident
	}

	// (typename) expr
	TypeCastExpr struct {
		Lparen token.Position // (
		Type   Typename
		Rparen token.Position // )
		X      Expr
	}

	ParenExpr struct {
		Lparen token.Position // (
		X      Expr
		Rparen token.Position // )
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
		*Range
		Type Typename
	}
)

func (*BadExpr) expr()                 {}
func (e *BadExpr) Beg() token.Position { return e.Position() }
func (e *BadExpr) End() token.Position {
	pos := e.Position()
	pos.Column += utf8.RuneCountInString(e.Literal())
	return pos
}

func (*Ident) expr()                 {}
func (e *Ident) String() string      { return e.Literal() }
func (e *Ident) Beg() token.Position { return e.Position() }
func (e *Ident) End() token.Position {
	pos := e.Position()
	pos.Column += utf8.RuneCountInString(e.Literal())
	return pos
}

func (*BasicLit) expr()                 {}
func (e *BasicLit) String() string      { return e.Literal() }
func (e *BasicLit) Beg() token.Position { return e.Position() }
func (e *BasicLit) End() token.Position {
	pos := e.Position()
	pos.Column += utf8.RuneCountInString(e.Literal())
	return pos
}

func (*CompoundLit) expr()                 {}
func (e *CompoundLit) Beg() token.Position { return e.Lparen }
func (e *CompoundLit) End() token.Position { return e.InitList.Rbrace }

func (*InitializerExpr) expr()                 {}
func (e *InitializerExpr) Beg() token.Position { return e.Lbrace }
func (e *InitializerExpr) End() token.Position { return e.Rbrace }

func (*RecordDesignatorExpr) expr()                 {}
func (e *RecordDesignatorExpr) Beg() token.Position { return e.Dot }
func (e *RecordDesignatorExpr) End() token.Position { return e.X.End() }

func (*ArrayDesignatorExpr) expr()                 {}
func (e *ArrayDesignatorExpr) Beg() token.Position { return e.Lbrack }
func (e *ArrayDesignatorExpr) End() token.Position { return e.X.End() }

func (*IndexExpr) expr()                 {}
func (e *IndexExpr) Beg() token.Position { return e.Arr.Beg() }
func (e *IndexExpr) End() token.Position { return e.Rbrack }

func (*SelectorExpr) expr()                 {}
func (e *SelectorExpr) Beg() token.Position { return e.X.Beg() }
func (e *SelectorExpr) End() token.Position { return e.Name.End() }

func (*CallExpr) expr()                 {}
func (e *CallExpr) Beg() token.Position { return e.Func.Beg() }
func (e *CallExpr) End() token.Position { return e.Rparen }

func (*ParenExpr) expr()                 {}
func (e *ParenExpr) Beg() token.Position { return e.Lparen }
func (e *ParenExpr) End() token.Position { return e.Rparen }

func (*UnaryExpr) expr()                 {}
func (e *UnaryExpr) Beg() token.Position { return e.Op.Position() }
func (e *UnaryExpr) End() token.Position { return e.X.End() }

func (*TypeCastExpr) expr()                 {}
func (e *TypeCastExpr) Beg() token.Position { return e.Lparen }
func (e *TypeCastExpr) End() token.Position { return e.X.End() }

func (*BinaryExpr) expr()                 {}
func (e *BinaryExpr) Beg() token.Position { return e.X.Beg() }
func (e *BinaryExpr) End() token.Position { return e.Y.End() }

func (*CondExpr) expr()                 {}
func (e *CondExpr) Beg() token.Position { return e.X.Beg() }
func (e *CondExpr) End() token.Position { return e.Else.End() }

func (*ConstantExpr) expr()                 {}
func (e *ConstantExpr) Beg() token.Position { return e.X.Beg() }
func (e *ConstantExpr) End() token.Position { return e.X.End() }

func (*AssignExpr) expr()                 {}
func (e *AssignExpr) Beg() token.Position { return e.X.Beg() }
func (e *AssignExpr) End() token.Position { return e.Y.End() }

func (*CommaExpr) expr()                 {}
func (e *CommaExpr) Beg() token.Position { return (*e)[0].Beg() }
func (e *CommaExpr) End() token.Position { return (*e)[len(*e)-1].End() }

func (*SizeOfExpr) expr()                 {}
func (e *SizeOfExpr) Beg() token.Position { return e.Range.Begin }
func (e *SizeOfExpr) End() token.Position { return e.Range.End }

type (
	Qualifier map[string]token.Position

	// 结构体/联合体
	RecordType struct {
		Qua *Qualifier

		Type      token.Token
		Name      *Ident
		Completed bool
		Lbrace    token.Position // {
		Fields    []*RecordField
		Rbrace    token.Position // }
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

	EnumType struct {
		Qua *Qualifier

		Enum      token.Position // enum
		Name      *Ident
		Completed bool
		Lbrace    token.Position // {
		List      []*EnumFieldDecl
		Rbrace    token.Position // }
	}

	// 内置类型
	BuildInType struct {
		*Range
		Qua  *Qualifier
		Type BasicType
	}

	Specifier map[string]bool

	// 指针类型
	PointerType struct {
		Qua *Qualifier

		Pointer token.Position
		Type    Typename
	}

	// 数组类型
	ArrayType struct {
		Qua *Qualifier

		Type   Typename
		Static bool
		Const  bool

		// T [*]
		// T []
		Incomplete bool

		Lbrack token.Position // [
		Size   Expr
		Rbrack token.Position // ]
	}

	// 函数类型
	ParamVarDecl struct {
		Qua  *StorageSpecifier
		Name *Ident
		Type Typename
	}

	ParamList []*ParamVarDecl

	FuncType struct {
		Qua *Qualifier

		Return   Typename
		Lparen   token.Position // (
		Params   ParamList
		Ellipsis bool           // ...
		Rparen   token.Position // )
	}

	// 括号 (abstract)
	ParenType struct {
		Lparen token.Position // (
		Type   Typename
		Rparen token.Position // )
	}
)

func (*RecordType) typeName()               {}
func (t *RecordType) Qualifier() *Qualifier { return t.Qua }
func (t *RecordType) Beg() token.Position   { return t.Type.Position() }
func (t *RecordType) End() token.Position {
	if !t.Completed {
		return t.Name.End()
	}
	return t.Rbrace
}

func (*EnumType) typeName() {}
func (t *EnumType) Qualifier() *Qualifier {
	if t.Qua == nil {
		t.Qua = &Qualifier{}
	}
	return t.Qua
}
func (t *EnumType) Beg() token.Position { return t.Enum }
func (t *EnumType) End() token.Position {
	if !t.Completed {
		return t.Name.End()
	}
	return t.Rbrace
}

func (*BuildInType) typeName() {}
func (t *BuildInType) Qualifier() *Qualifier {
	if t.Qua == nil {
		t.Qua = &Qualifier{}
	}
	return t.Qua
}
func (t *BuildInType) Beg() token.Position { return t.Range.Begin }
func (t *BuildInType) End() token.Position { return t.Range.End }

func (*PointerType) typeName() {}
func (t *PointerType) Qualifier() *Qualifier {
	if t.Qua == nil {
		t.Qua = &Qualifier{}
	}
	return t.Qua
}
func (t *PointerType) Beg() token.Position { return t.Pointer }
func (t *PointerType) End() token.Position { return t.Type.End() }

func (*ArrayType) typeName() {}
func (t *ArrayType) Qualifier() *Qualifier {
	if t.Qua == nil {
		t.Qua = &Qualifier{}
	}
	return t.Qua
}
func (t *ArrayType) Beg() token.Position { return t.Type.Beg() }
func (t *ArrayType) End() token.Position { return t.Rbrack }

func (*FuncType) typeName() {}
func (t *FuncType) Qualifier() *Qualifier {
	if t.Qua == nil {
		t.Qua = &Qualifier{}
	}
	return t.Qua
}
func (t *FuncType) Beg() token.Position { return t.Return.Beg() }
func (t *FuncType) End() token.Position { return t.Rparen }

func (*ParenType) typeName()               {}
func (t *ParenType) Qualifier() *Qualifier { return t.Type.Qualifier() }
func (t *ParenType) Beg() token.Position   { return t.Lparen }
func (t *ParenType) End() token.Position   { return t.Rparen }

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
	if v, ok := t.Type.(*FuncType); ok {
		return fmt.Sprintf("%s (*%s) %s", v.Return, t.Qua.String(), funcParamString(v))
	}
	return fmt.Sprintf("%s *%s", t.Type.String(), t.Qua.String())
}

func (t *ArrayType) String() string {
	return fmt.Sprintf("%s %s[]", t.Qualifier().String(), t.Type.String())
}

func funcParamString(t *FuncType) string {
	var params []string
	for _, v := range t.Params {
		params = append(params, v.Type.String())
	}
	if t.Ellipsis {
		params = append(params, "...")
	}
	return "(" + strings.Join(params, ",") + ")"
}

func (t *FuncType) String() string {
	var params []string
	for _, v := range t.Params {
		params = append(params, v.Type.String())
	}
	if t.Ellipsis {
		params = append(params, "...")
	}
	return fmt.Sprintf("%s %s", t.Return.String(), funcParamString(t))
}

func (t *ParenType) String() string {
	return fmt.Sprintf("(%s)", t.Type.String())
}

type (
	// 定义语句
	DeclStmt []Decl

	// typedef /extern /static /auto /register
	StorageSpecifier map[string]token.Position

	LabelStmt struct {
		Id   *Ident
		Stmt Stmt
	}

	CaseStmt struct {
		Case token.Position // case
		Expr Expr
		Stmt Stmt
	}

	DefaultStmt struct {
		Default token.Position // case
		Stmt    Stmt
	}

	CompoundStmt struct {
		Lbrace token.Position // {
		Stmts  []Stmt
		Rbrace token.Position // }
	}

	ExprStmt struct {
		Expr      Expr
		Semicolon token.Position // ;
	}

	IfStmt struct {
		If   token.Position // if
		X    Expr
		Then Stmt
		Else Stmt
	}

	SwitchStmt struct {
		Switch token.Position // switch
		X      Expr
		Stmt   Stmt
	}

	WhileStmt struct {
		While token.Position // while
		X     Expr
		Stmt  Stmt
	}

	DoWhileStmt struct {
		Do        token.Position // do
		Stmt      Stmt
		X         Expr
		Semicolon token.Position // ;
	}

	ForStmt struct {
		For  token.Position
		Init Expr
		Decl Stmt
		Cond Expr
		Post Expr
		Stmt Stmt
	}

	GotoStmt struct {
		Goto      token.Position
		Id        *Ident
		Semicolon token.Position // ;
	}

	ContinueStmt struct {
		Continue  token.Position
		Semicolon token.Position // ;
	}

	BreakStmt struct {
		Break     token.Position
		Semicolon token.Position // ;
	}

	ReturnStmt struct {
		Return    token.Position
		X         Expr
		Semicolon token.Position // ;
	}
)

func (*DeclStmt) stmt()                 {}
func (s *DeclStmt) Beg() token.Position { return (*s)[0].Beg() }
func (s *DeclStmt) End() token.Position { return (*s)[len(*s)-1].Beg() }

func (*LabelStmt) stmt()                 {}
func (s *LabelStmt) Beg() token.Position { return s.Id.Beg() }
func (s *LabelStmt) End() token.Position { return s.Stmt.End() }

func (*CaseStmt) stmt()                 {}
func (s *CaseStmt) Beg() token.Position { return s.Case }
func (s *CaseStmt) End() token.Position { return s.Stmt.End() }

func (*DefaultStmt) stmt()                 {}
func (s *DefaultStmt) Beg() token.Position { return s.Default }
func (s *DefaultStmt) End() token.Position { return s.Stmt.End() }

func (*CompoundStmt) stmt()                 {}
func (s *CompoundStmt) Beg() token.Position { return s.Lbrace }
func (s *CompoundStmt) End() token.Position { return s.Rbrace }

func (*ExprStmt) stmt()                 {}
func (s *ExprStmt) Beg() token.Position { return s.Expr.Beg() }
func (s *ExprStmt) End() token.Position { return s.Semicolon }

func (*IfStmt) stmt()                 {}
func (s *IfStmt) Beg() token.Position { return s.If }
func (s *IfStmt) End() token.Position {
	if s.Else == nil {
		return s.Then.End()
	}
	return s.Else.End()
}

func (*SwitchStmt) stmt()                 {}
func (s *SwitchStmt) Beg() token.Position { return s.Switch }
func (s *SwitchStmt) End() token.Position { return s.Stmt.End() }

func (*WhileStmt) stmt()                 {}
func (s *WhileStmt) Beg() token.Position { return s.While }
func (s *WhileStmt) End() token.Position { return s.Stmt.End() }

func (*DoWhileStmt) stmt()                 {}
func (s *DoWhileStmt) Beg() token.Position { return s.Do }
func (s *DoWhileStmt) End() token.Position { return s.Semicolon }

func (*ForStmt) stmt()                 {}
func (s *ForStmt) Beg() token.Position { return s.For }
func (s *ForStmt) End() token.Position { return s.Stmt.End() }

func (*GotoStmt) stmt()                 {}
func (s *GotoStmt) Beg() token.Position { return s.Goto }
func (s *GotoStmt) End() token.Position { return s.Semicolon }

func (*ContinueStmt) stmt()                 {}
func (s *ContinueStmt) Beg() token.Position { return s.Continue }
func (s *ContinueStmt) End() token.Position { return s.Semicolon }

func (*BreakStmt) stmt()                 {}
func (s *BreakStmt) Beg() token.Position { return s.Break }
func (s *BreakStmt) End() token.Position { return s.Semicolon }

func (*ReturnStmt) stmt()                 {}
func (s *ReturnStmt) Beg() token.Position { return s.Return }
func (s *ReturnStmt) End() token.Position { return s.Semicolon }

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
		Name *Ident
		Type *FuncType
		Decl []Decl
		Body *CompoundStmt
	}

	// 变量定义
	VarDecl struct {
		Type Typename
		Name *Ident
		Init Expr
	}

	// 类型定义
	// typedef type-name ident;
	TypedefDecl struct {
		Typedef token.Position
		Type    Typename
		Name    *Ident
	}
)

func (*FuncDecl) decl() {}
func (t *FuncDecl) Ident() *Ident {
	return t.Name
}
func (t *FuncDecl) Beg() token.Position { return t.Type.Beg() }
func (t *FuncDecl) End() token.Position {
	return t.Body.End()
}

func (*VarDecl) decl() {}
func (t *VarDecl) Ident() *Ident {
	return t.Name
}
func (t *VarDecl) Beg() token.Position { return t.Type.Beg() }
func (t *VarDecl) End() token.Position {
	if t.Init != nil {
		return t.Init.End()
	}
	return t.Name.End()
}

func (*TypedefDecl) decl() {}
func (t *TypedefDecl) Ident() *Ident {
	return t.Name
}
func (t *TypedefDecl) Beg() token.Position { return t.Typedef }
func (t *TypedefDecl) End() token.Position {
	return t.Name.End()
}

func (*ParamVarDecl) decl() {}
func (t *ParamVarDecl) Ident() *Ident {
	return t.Name
}
func (t *ParamVarDecl) Beg() token.Position { return t.Type.Beg() }
func (t *ParamVarDecl) End() token.Position {
	return t.Name.End()
}

func (*EnumFieldDecl) decl() {}
func (t *EnumFieldDecl) Ident() *Ident {
	return t.Name
}
func (t *EnumFieldDecl) Beg() token.Position { return t.Name.Beg() }
func (t *EnumFieldDecl) End() token.Position {
	if t.Val != nil {
		return t.Val.End()
	}
	return t.Name.End()
}
