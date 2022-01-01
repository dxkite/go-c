package parser

import (
	"dxkite.cn/c/ast"
	"dxkite.cn/c/errors"
)

type environment struct {
	global *ast.Scope // 全局作用域
	nested *ast.Scope // 文件作用域
	label  *ast.Scope
	labels []*ast.Ident
	parser *parser
}

func newEnv(glb *ast.Scope, p *parser) *environment {
	env := &environment{}
	env.global = glb
	env.nested = ast.NewScope(ast.GlobalScope, env.global, int(ast.MaxNestedNamespace))
	env.parser = p
	return env
}

// 定义对象
func (e *environment) declare(obj *ast.Object) {
	var namespace ast.ScopeNamespace
	switch obj.Type {
	case ast.ObjectEnumName:
		namespace = ast.EnumScope
	case ast.ObjectStructName:
		namespace = ast.StructScope
	case ast.ObjectUnionName:
		namespace = ast.UnionScope
	case ast.ObjectLabelName:
		e.alterDeclare(e.label.Insert(ast.IdentScope, obj), obj, errors.ErrSyntaxRedefinedLabel)
	case ast.ObjectFunc:
		namespace = ast.IdentScope
		e.alterDeclare(e.global.Insert(ast.IdentScope, obj), obj, errors.ErrSyntaxRedefineIdent) // 函数注册到全局域
	default:
		namespace = ast.IdentScope
	}
	e.alterDeclare(e.nested.Insert(namespace, obj), obj, errors.ErrSyntaxRedefineIdent)
}

func (e *environment) declareDeclIdent(typ ast.ObjectType, decl ast.Decl) {
	if decl.Ident() == nil {
		return
	}
	e.declare(ast.NewDeclObject(typ, decl.Ident(), decl))
}

func (e *environment) declareRecord(r *ast.RecordType, completed bool) {
	if r.Name == nil {
		return
	}
	typ := ast.ObjectStructName
	if r.Type.Literal() == "union" {
		typ = ast.ObjectUnionName
	}
	obj := ast.NewObjectTypename(typ, r.Name, r)
	obj.Completed = completed
	e.declare(obj)
}

func (e *environment) declareEnum(r *ast.EnumType, completed bool) {
	if r.Name == nil {
		return
	}
	typ := ast.ObjectEnumName
	obj := ast.NewObjectTypename(typ, r.Name, r)
	obj.Completed = completed
	e.declare(obj)
}

func (e *environment) declareEnumTag(r *ast.EnumFieldDecl) {
	typ := ast.ObjectEnumTag
	obj := ast.NewDeclObject(typ, r.Name, r)
	e.declare(obj)
}

func (e *environment) alterDeclare(alt, obj *ast.Object, err errors.ErrCode) {
	if alt == nil {
		return
	}
	if alt.Type == obj.Type && obj.Type == ast.ObjectFunc {
		if !alt.Completed && obj.Completed {
			return
		}
		err = errors.ErrSyntaxRedefineFunc
	}
	if alt.Type == obj.Type && (alt.Type == ast.ObjectStructName ||
		alt.Type == ast.ObjectEnumName ||
		alt.Type == ast.ObjectUnionName) {
		if !alt.Completed && obj.Completed {
			return
		}
		switch alt.Type {
		case ast.ObjectStructName:
			err = errors.ErrSyntaxRedefinedStruct
		case ast.ObjectEnumName:
			err = errors.ErrSyntaxRedefinedEnum
		case ast.ObjectUnionName:
			err = errors.ErrSyntaxRedefinedUnion
		}
	}
	e.parser.addErr(obj.Pos, err, obj, alt)
}

// 解析变量
func (e *environment) tryResolve(space ast.ScopeNamespace, name string) *ast.Object {
	scope := e.nested
	for scope != nil {
		if obj := scope.Lookup(space, name); obj != nil {
			return obj
		}
		scope = scope.Outer
	}
	return nil
}

func (e *environment) isTypename(name string) ast.Typename {
	obj := e.tryResolve(ast.IdentScope, name)
	if obj != nil && obj.Type == ast.ObjectTypename {
		return obj.Typename
	}
	return nil
}

// 进入作用域
func (e *environment) enterScope(typ ast.ScopeType) {
	e.nested = ast.NewScope(typ, e.nested, int(ast.MaxNestedNamespace))
}

// 退出作用域
func (e *environment) leaveScope() {
	e.nested = e.nested.Outer
}

// 进入label作用域
func (e *environment) enterLabelScope() {
	e.label = ast.NewScope(ast.FuncScope, nil, 1)
}

// 解析标签
func (e *environment) tryResolveLabel(ident *ast.Ident) {
	e.labels = append(e.labels, ident)
	return
}

// 退出label作用域
func (e *environment) leaveLabelScope() (labels []*ast.Ident) {
	for _, name := range e.labels {
		if e.label.Lookup(0, name.Literal()) == nil {
			labels = append(labels, name)
		}
	}
	return
}
