package ast

import "dxkite.cn/c/token"

type ScopeType int

const (
	GlobalScope        ScopeType = iota // 全局件域
	FileScope                           // 文件域
	BlockScope                          // 块级别域
	FuncScope                           // 函数域
	FuncPrototypeScope                  // 参数域
)

// 定义的类型
type ObjectType int

const (
	ObjectVar ObjectType = iota
	ObjectFunc
	ObjectParamVal
	ObjectTypename
	ObjectEnumName
	ObjectStructName
	ObjectUnionName
	ObjectLabelName
	ObjectEnumTag
)

type Object struct {
	Type      ObjectType
	Name      string
	Pos       token.Position
	Completed bool // 定义完成

	Decl     Decl     // 定义语句
	Typename Typename // 类型名称
}

func NewObject(typ ObjectType, ident *Ident) *Object {
	return &Object{
		Type: typ,
		Name: ident.Literal(),
		Pos:  ident.Position(),
	}
}

func NewTypenameObject(ident *Ident, name Typename) *Object {
	return &Object{
		Type:     ObjectTypename,
		Name:     ident.Literal(),
		Pos:      ident.Position(),
		Typename: name,
	}
}

func NewObjectTypename(typ ObjectType, ident *Ident, name Typename) *Object {
	return &Object{
		Type:     typ,
		Name:     ident.Literal(),
		Pos:      ident.Position(),
		Typename: name,
	}
}

func NewDeclObject(typ ObjectType, ident *Ident, decl Decl) *Object {
	return &Object{
		Type: typ,
		Name: ident.Literal(),
		Pos:  ident.Position(),
		Decl: decl,
	}
}

type ScopeNamespace int

const (
	IdentScope         ScopeNamespace = iota // other
	StructScope                              // struct
	EnumScope                                // enum
	UnionScope                               // union
	MaxNestedNamespace                       // 可嵌套的作用域
)

type Scope struct {
	Outer   *Scope
	Type    ScopeType
	Objects []map[string]*Object
	cap     int
}

func NewScope(typ ScopeType, outer *Scope, cap int) *Scope {
	s := &Scope{
		Outer: outer,
		Type:  typ,
		cap:   cap,
	}
	s.Objects = make([]map[string]*Object, cap)
	for i := 0; i < cap; i++ {
		s.Objects[i] = map[string]*Object{}
	}
	return s
}

func (s *Scope) Lookup(name ScopeNamespace, id string) *Object {
	if int(name) >= s.cap {
		return nil
	}
	return s.Objects[name][id]
}

func (s *Scope) Insert(name ScopeNamespace, obj *Object) (alt *Object) {
	if int(name) >= s.cap {
		return nil
	}
	if alt = s.Objects[name][obj.Name]; alt == nil {
		s.Objects[name][obj.Name] = obj
	}
	return
}
