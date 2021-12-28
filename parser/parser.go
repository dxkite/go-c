package parser

import (
	"dxkite.cn/c/ast"
	"dxkite.cn/c/errors"
	"dxkite.cn/c/scanner"
	"dxkite.cn/c/token"
)

// 表达式解析
type Parser struct {
	// 定义的类型
	typeName map[string]*ast.TypedefDecl
	// 错误处理
	err errors.ErrorHandler
	// 当前token
	cur token.Token
	// 当前输入
	r scanner.Scanner
}

func NewParser(r scanner.Scanner, err errors.ErrorHandler) *Parser {
	p := &Parser{
		err: err,
		r:   scanner.NewTokenScan(r),
	}
	p.next()
	return p
}

const (
	LowestPrec = 0  // 最低优先级
	BinaryPrec = 2  // 二元操作符
	UnaryPrec  = 12 // 一元操作符
)

func precedence(lit string) int {
	switch lit {
	case "?":
		return 1
	case "||":
		return 2
	case "&&":
		return 3
	case "|":
		return 4
	case "^":
		return 5
	case "&":
		return 6
	case "!=", "==":
		return 7
	case ">", "<", ">=", "<=":
		return 8
	case ">>", "<<":
		return 9
	case "+", "-":
		return 10
	case "*", "/", "%":
		return 11
	}
	return LowestPrec
}

// primary-expression: identifier | constant | string-literal | ( expression )
func (p *Parser) parsePrimaryExpr() ast.Expr {
	if p.cur.Type() == token.PUNCTUATOR && p.cur.Literal() == "(" {
		p.next()
		x := p.parseExpr()
		p.exceptPunctuator(")")
		return &ast.ParenExpr{
			X: x,
		}
	}
	switch p.cur.Type() {
	case token.IDENT:
		cur := p.cur
		p.next()
		return &ast.Ident{Token: cur}
	case token.INT, token.CHAR, token.FLOAT, token.STRING:
		cur := p.cur
		p.next()
		return &ast.BasicLit{Token: cur}
	}
	exp := ast.BadExpr([]token.Token{p.cur})
	return &exp
}

func (p *Parser) parsePostfixExpr() ast.Expr {
	// ( typename ) { init-list }
	if p.cur.Type() == token.PUNCTUATOR && p.cur.Literal() == "(" {
		if t := p.peekOne(); p.isTypeNameTok(t) {
			return p.parseCompoundLitExpr()
		}
	}
	expr := p.parsePrimaryExpr()
	return p.parsePostfixExprInner(expr)
}

func (p *Parser) parsePostfixExprInner(expr ast.Expr) ast.Expr {
	switch p.cur.Literal() {
	case "[":
		p.next() // [
		idx := p.parseExpr()
		p.exceptPunctuator("]")
		expr = &ast.IndexExpr{
			Arr:   expr,
			Index: idx,
		}
	case "++", "--":
		op := p.cur
		p.next()
		expr = &ast.UnaryExpr{
			Op: op,
			X:  expr,
		}
	case "->", ".":
		op := p.cur
		p.next() // .
		name := p.expectIdent()
		expr = &ast.SelectorExpr{
			X:    expr,
			Op:   op,
			Name: &ast.Ident{Token: name},
		}
	case "(":
		arg := p.parseArgsExpr()
		expr = &ast.CallExpr{
			Fun:  expr,
			Args: arg,
		}
	default:
		return expr
	}
	return p.parsePostfixExprInner(expr)
}

func (p *Parser) parseUnaryExpr() ast.Expr {
	if p.cur.Type() == token.PUNCTUATOR || p.cur.Literal() == "sizeof" {
		switch p.cur.Literal() {
		case "++", "--", "&", "*", "+", "-", "~", "!":
			op := p.cur
			p.next() //
			return &ast.UnaryExpr{
				Op: op,
				X:  p.parseUnaryExpr(),
			}
		case "sizeof":
			op := p.cur
			p.next() // sizeof
			if t := p.peekOne(); p.cur.Literal() == "(" && p.isTypeNameTok(t) {
				p.next()                  // (
				name := p.parseTypeName() // type-name
				p.exceptPunctuator(")")   // )
				return &ast.SizeOfExpr{
					Type: name,
				}
			}
			return &ast.UnaryExpr{
				Op: op,
				X:  p.parseUnaryExpr(),
			}
		}
	}
	return p.parsePostfixExpr()
}

func (p *Parser) parseCastExpr() ast.Expr {
	if t := p.peekOne(); p.cur.Literal() == "(" && p.isTypeNameTok(t) {
		p.next()                  // (
		name := p.parseTypeName() // type-name
		p.exceptPunctuator(")")   // )
		if p.cur.Literal() == "{" {
			expr := p.parseInitializerList()
			return &ast.CompoundLit{
				Type:     name,
				InitList: expr,
			}
		}
		expr := p.parseCastExpr()
		return &ast.TypeCastExpr{
			X:    expr,
			Type: name,
		}
	}
	return p.parseUnaryExpr()
}

func (p *Parser) selectExpr(prec int) (expr ast.Expr) {
	if prec >= UnaryPrec {
		return p.parseCastExpr()
	} else {
		return p.parseBinaryExpr(prec)
	}
}

func (p *Parser) parseBinaryExpr(prec int) ast.Expr {
	x := p.selectExpr(prec + 1)
	if p.cur.Type() == token.PUNCTUATOR && precedence(p.cur.Literal()) >= prec {
		op := p.cur
		p.next() // op
		y := p.selectExpr(prec + 1)
		return &ast.BinaryExpr{
			X:  x,
			Op: op,
			Y:  y,
		}
	}
	return x
}

func (p *Parser) parseCondExpr() ast.Expr {
	x := p.parseBinaryExpr(BinaryPrec)
	if p.cur.Type() == token.PUNCTUATOR && p.cur.Literal() == "?" {
		op := p.cur
		p.next()
		then := p.parseExpr()
		p.exceptPunctuator(":")
		el := p.parseCondExpr()
		return &ast.CondExpr{
			X:    x,
			Op:   op,
			Then: then,
			Else: el,
		}
	}
	return x
}

// const expr
func (p *Parser) parseConstantExpr() ast.Expr {
	return &ast.ConstantExpr{X: p.parseCondExpr()}
}

func (p *Parser) parseAssignExpr() ast.Expr {
	x := p.parseCondExpr()
	switch p.cur.Literal() {
	case "=", "*=", "/=", "%=", "+=", "-=", "<<=", "==>", "&=", "^=", "|=":
		op := p.cur
		p.next()
		return &ast.AssignExpr{
			X:  x,
			Op: op,
			Y:  p.parseAssignExpr(),
		}
	}
	return x
}

func (p *Parser) parseExpr() ast.Expr {
	x := p.parseAssignExpr()
	comma := ast.CommaExpr{x}
	if p.cur.Literal() == "," {
		for p.cur.Literal() == "," {
			p.next() // ,
			x = p.parseAssignExpr()
			comma = append(comma, x)
		}
		return &comma
	}
	return comma[0]
}

func (p *Parser) parseCompoundLitExpr() ast.Expr {
	p.next() // (
	typeName := p.parseTypeName()
	p.exceptPunctuator(")")
	expr := p.parseInitializerList()
	return &ast.CompoundLit{
		Type:     typeName,
		InitList: expr,
	}
}

func (p *Parser) parseInitializerList() *ast.InitializerExpr {
	p.exceptPunctuator("{")
	list := ast.InitializerExpr{}
	for p.cur.Literal() != "}" {
		item := p.parseInitializer()
		list = append(list, item)
		if t := p.peekOne(); p.cur.Literal() != "}" && t.Literal() != "}" {
			p.exceptPunctuator(",")
		} else {
			if p.cur.Literal() == "," {
				p.next() //,
			}
		}
	}
	p.exceptPunctuator("}")
	return &list
}

func (p *Parser) parseInitializer() ast.Expr {
	if p.cur.Literal() == "{" {
		return p.parseInitializerList()
	}
	if l := p.cur.Literal(); l == "." || l == "[" {
		return p.parseDesignationInitExpr()
	}
	return p.parseAssignExpr()
}

func (p *Parser) parseDesignationInitExpr() ast.Expr {
	designator := p.parseDesignator()
	if designator != nil {
		p.exceptPunctuator("=")
		expr := p.parseAssignExpr()
		applyDesignator(designator, expr)
		return designator
	}
	return p.parseAssignExpr()
}

func applyDesignator(des ast.Expr, assign ast.Expr) {
	switch t := des.(type) {
	case *ast.ArrayDesignatorExpr:
		if t.X != nil {
			applyDesignator(t.X, assign)
			return
		}
		t.X = assign
	case *ast.RecordDesignatorExpr:
		if t.X != nil {
			applyDesignator(t.X, assign)
			return
		}
		t.X = assign
	}
}

func (p *Parser) parseDesignator() ast.Expr {
	if p.cur.Literal() == "." {
		p.next()
		ident := p.expectIdent()
		x := p.parseDesignator()
		return &ast.RecordDesignatorExpr{
			Field: &ast.Ident{ident},
			X:     x,
		}
	}
	if p.cur.Literal() == "[" {
		p.next()
		expr := p.parseConstantExpr()
		p.exceptPunctuator("]")
		x := p.parseDesignator()
		return &ast.ArrayDesignatorExpr{
			Index: expr,
			X:     x,
		}
	}
	return nil
}

func (p *Parser) parseArgsExpr() []ast.Expr {
	p.next() // (
	var list []ast.Expr
	for p.cur.Literal() != ")" {
		item := p.parseAssignExpr()
		list = append(list, item)
		if p.cur.Literal() != "," {
			break
		}
		p.exceptPunctuator(",")
	}
	p.exceptPunctuator(")")
	return list
}

var typeQualifier = []string{"const", "restrict", "volatile"}
var typeSpecifier = []string{"void", "char", "short", "int", "long", "float", "double", "signed", "unsigned", "_Bool", "_Complex"}
var functionSpecifier = []string{"inline"}
var storageClassSpecifier = []string{"typedef", "extern", "static", "auto", "register"}

// 结构化类型
var typeStructMap = map[string]bool{
	"struct": true,
	"union":  true,
	"enum":   true,
}

// 类型
// "void", "char", "short", "int", "long", "float", "double", "signed", "unsigned", "_Bool", "_Complex"
var typeSpecifierMap = map[string]bool{}

// 声明
var declarationSpecifierMap = map[string]bool{}

// "const", "restrict", "volatile"
var typeQualifierMap = map[string]bool{}
var storageClassSpecifierMap = map[string]bool{}
var typeSpecifierQualifierMap = map[string]bool{}

func init() {
	typeSpecifierQualifierMap = typeStructMap
	for _, v := range typeQualifier {
		typeQualifierMap[v] = true
		typeSpecifierQualifierMap[v] = true
	}

	typeSpecifierMap = typeStructMap
	for _, v := range typeSpecifier {
		typeSpecifierMap[v] = true
		typeSpecifierQualifierMap[v] = true
	}

	declarationSpecifierMap = typeSpecifierMap
	for _, v := range functionSpecifier {
		declarationSpecifierMap[v] = true
	}
	for _, v := range storageClassSpecifier {
		declarationSpecifierMap[v] = true
		storageClassSpecifierMap[v] = true
	}
}

// 解析类型名称
func (p *Parser) parseTypeName() ast.TypeName {
	basic := p.parseTypeQualifierSpecifierList()
	return p.parseAbstractDeclarator(basic)
}

func (p *Parser) parseAbstractDeclarator(inner ast.TypeName) ast.TypeName {
	switch p.cur.Literal() {
	case "*":
		inner = p.parsePointer(inner)
	case "(":
		if t := p.peekOne().Literal(); t == "*" || t == "[" || t == "(" {
			p.next() // (
			inner = p.parseAbstractDeclarator(&ast.ParenType{Inner: inner})
			p.exceptPunctuator(")")
		} else {
			inner = p.parseFuncType(inner)
		}
	case "[":
		inner = p.parseArrayType(inner)
	default:
		return inner
	}
	return p.parseAbstractDeclarator(inner)
}

func (p *Parser) parseFuncType(inner ast.TypeName) ast.TypeName {
	params, ellipsis := p.parseParameterList()
	return &ast.FuncType{
		Inner:    inner,
		Params:   params,
		Ellipsis: ellipsis,
	}
}

func (p *Parser) parseArrayType(inner ast.TypeName) ast.TypeName {
	if t := p.peekOne(); t.Literal() == "*" || t.Literal() == "]" {
		p.next() // [
		if p.cur.Literal() == "*" {
			p.next() // *
		}
		p.exceptPunctuator("]")
		return &ast.IncompleteArrayType{Inner: inner}
	}
	return p.parseArrayTypeExpr(inner)
}

func (p *Parser) parseArrayTypeExpr(inner ast.TypeName) ast.TypeName {
	p.next() // [
	static := false
	var qua []token.Token
	for typeQualifierMap[p.cur.Literal()] || p.cur.Literal() == "static" {
		if p.cur.Literal() == "static" {
			static = true
		} else {
			qua = append(qua, p.cur)
		}
	}
	expr := p.parseAssignExpr()
	p.exceptPunctuator("]") // ]

	if v, ok := expr.(*ast.AssignExpr); ok {
		arr := &ast.ArrayType{
			Inner:     inner,
			Qualifier: &ast.Qualifier{},
			Static:    static,
			Size:      v,
		}
		p.markQualifier(arr.Qualifier, qua)
		return arr
	}

	// 常量表达式
	return &ast.ConstArrayType{
		Inner: inner,
		Size:  expr,
	}
}

func (p *Parser) parseParameterList() (params ast.ParamList, ellipsis bool) {
	p.exceptPunctuator("(")
	for p.cur.Literal() != ")" && p.cur.Literal() != "..." {
		param := p.parseParameterDecl()
		params = append(params, param)
		if p.cur.Literal() != "," {
			break
		}
		p.exceptPunctuator(",")
	}
	if p.cur.Literal() == "..." {
		ellipsis = true
		p.next()
	}
	p.exceptPunctuator(")")
	return
}

func (p *Parser) parseParameterDecl() *ast.ParamVarDecl {
	typ := p.parseDeclarationSpecifiers()
	param := &ast.ParamVarDecl{}
	param.Type, param.Name = p.parseDeclarator(typ)
	return param
}

func (p *Parser) parseDeclarator(inner ast.TypeName) (ast.TypeName, *ast.Ident) {
	if p.cur.Literal() == "*" {
		inner = p.parsePointer(inner)
	}
	return p.parseDirectDeclarator(inner)
}

func (p *Parser) parseDirectDeclarator(inner ast.TypeName) (ast.TypeName, *ast.Ident) {
	var ident *ast.Ident
	if p.cur.Literal() == "(" {
		p.next() // (
		typ, ident := p.parseDeclarator(inner)
		p.exceptPunctuator(")")
		typ = &ast.ParenType{Inner: typ}
		typ = p.parseDirectDeclaratorInner(typ)
		return typ, ident
	}
	if p.cur.Type() == token.IDENT {
		tok := p.expectIdent()
		ident = &ast.Ident{Token: tok}
	}
	inner = p.parseDirectDeclaratorInner(inner)
	return inner, ident
}

func (p *Parser) parseDirectDeclaratorInner(typ ast.TypeName) ast.TypeName {
	switch p.cur.Literal() {
	case "(":
		typ = p.parseFuncType(typ)
	case "[":
		typ = p.parseArrayType(typ)
	default:
		return typ
	}
	return p.parseDirectDeclaratorInner(typ)
}

// ( type-specifier | type-qualifier ) +
func (p *Parser) parseTypeQualifierSpecifierList() ast.TypeName {
	var qua []token.Token
	var typ ast.TypeName
	var buildIn *ast.BuildInType

	for typeSpecifierQualifierMap[p.cur.Literal()] {
		if typeQualifierMap[p.cur.Literal()] {
			qua = append(qua, p.cur)
			p.next()
			continue
		}
		if buildIn != nil && typeStructMap[p.cur.Literal()] {
			p.addErr(p.cur, errors.ErrSyntaxUnexpectedTypeSpecifier, p.cur.Literal())
		}
		typ = p.parseTypeSpecifier()
		if v, ok := typ.(*ast.BuildInType); ok {
			if buildIn == nil {
				buildIn = v
			} else {
				buildIn.Type = append(buildIn.Type, v.Type...)
			}
		}
	}

	if len(qua) > 0 {
		t := &ast.TypeQualifier{
			Qualifier: &ast.Qualifier{},
			Inner:     typ,
		}
		t.Qualifier = &ast.Qualifier{}
		p.markQualifier(t.Qualifier, qua)
		typ = t
	}
	return typ
}

// (('*') typeQualifierList?)+
func (p *Parser) parsePointer(inner ast.TypeName) (t ast.TypeName) {
	p.next() // *
	tt := &ast.PointerType{Inner: inner}
	tks := p.scanTypeQualifierTok()
	tt.Qualifier = &ast.Qualifier{}
	p.markQualifier(tt.Qualifier, tks)
	t = tt
	for p.cur.Literal() == "*" {
		t = p.parsePointer(t)
	}
	return t
}

// 扫描类型
func (p *Parser) parseDeclarationSpecifiers() ast.TypeName {
	var qua []token.Token
	var typ ast.TypeName
	var buildIn *ast.BuildInType
	var spec []token.Token

	for declarationSpecifierMap[p.cur.Literal()] {
		if typeQualifierMap[p.cur.Literal()] {
			qua = append(qua, p.cur)
			p.next()
			continue
		}
		if storageClassSpecifierMap[p.cur.Literal()] {
			spec = append(spec, p.cur)
			p.next()
			continue
		}
		if buildIn != nil && typeStructMap[p.cur.Literal()] {
			p.addErr(p.cur, errors.ErrSyntaxUnexpectedTypeSpecifier, p.cur.Literal())
		}
		typ = p.parseTypeSpecifier()
		if v, ok := typ.(*ast.BuildInType); ok {
			if buildIn == nil {
				buildIn = v
			} else {
				buildIn.Type = append(buildIn.Type, v.Type...)
			}
		}
	}

	if len(qua) > 0 {
		t := &ast.TypeQualifier{
			Qualifier: &ast.Qualifier{},
			Inner:     typ,
		}
		p.markQualifier(t.Qualifier, qua)
		typ = t
	}

	if len(spec) > 0 {
		st := &ast.TypeStorageSpecifier{
			Specifier: &ast.Specifier{},
			Inner:     typ,
		}
		p.markSpecifier(st.Specifier, spec)
		typ = st
	}
	return typ
}

func (p *Parser) parseTypeSpecifier() ast.TypeName {
	switch p.cur.Literal() {
	case "struct", "union":
		return p.parseRecordType()
	case "enum":
		return p.parseEnumType()
	default:
		// 用户定义的类型
		if p.cur.Type() == token.IDENT && p.isTypedefName(p.cur) {
			return p.typeName[p.cur.Literal()].Type
		}
	}
	return p.parseBuildInType()
}

// 扫描内置类型
func (p *Parser) parseBuildInType() ast.TypeName {
	var spec []token.Token
	for p.cur.Type() == token.KEYWORD && typeSpecifierMap[p.cur.Literal()] {
		if typeSpecifierMap[p.cur.Literal()] {
			spec = append(spec, p.cur)
		}
		p.next()
	}
	return &ast.BuildInType{
		Type: spec,
	}
}

func (p *Parser) markQualifier(q *ast.Qualifier, qua []token.Token) {
	for _, t := range qua {
		if (*q)[t.Literal()] {
			p.addWarn(t, errors.ErrSyntaxDuplicateTypeQualifier, t.Literal())
		}
		(*q)[t.Literal()] = true
	}
}

func (p *Parser) markSpecifier(q *ast.Specifier, qua []token.Token) {
	for _, t := range qua {
		if (*q)[t.Literal()] {
			p.addWarn(t, errors.ErrSyntaxDuplicateTypeSpecifier, t.Literal())
		}
		(*q)[t.Literal()] = true
	}
}

func (p *Parser) parseRecordType() *ast.RecordType {
	t := p.cur
	p.next() // struct union
	r := &ast.RecordType{Type: t}

	if p.cur.Literal() != "{" {
		tok := p.expectIdent()
		r.Name = &ast.Ident{Token: tok}
	}

	if p.cur.Literal() != "{" {
		return r
	}

	p.next() // {
	for p.cur.Literal() != "}" {
		typ := p.parseTypeQualifierSpecifierList()
		for {
			f := &ast.RecordField{}
			typ, ident := p.parseDeclarator(typ)
			f.Type = typ
			f.Name = ident
			// bit-field
			if p.cur.Literal() == ":" {
				p.next() // :
				expr := p.parseConstantExpr()
				f.Bit = expr
			}
			// bit field
			// TODO 递归类型引用检查
			if f.Bit == nil && f.Name == nil && !isRecordType(typ) {
				p.addErr(p.cur, errors.ErrSyntaxExpectedRecordMemberName)
				break
			}
			r.Fields = append(r.Fields, f)
			if p.cur.Literal() != "," {
				break
			}
			p.exceptPunctuator(",")
		}
		p.exceptPunctuator(";")
	}
	p.exceptPunctuator("}") // }
	return r
}

func isRecordType(typ ast.TypeName) bool {
	switch v := typ.(type) {
	case *ast.RecordType:
		return true
	case *ast.PointerType:
		return isRecordType(v.Inner)
	case *ast.TypeQualifier:
		return isRecordType(v.Inner)
	case *ast.TypeStorageSpecifier:
		return isRecordType(v.Inner)
	}
	return false
}

func (p *Parser) parseEnumType() *ast.EnumType {
	p.next() // enum
	t := &ast.EnumType{}
	if p.cur.Type() == token.IDENT {
		t.Name = &ast.Ident{Token: p.cur}
		p.next()
	}
	if p.cur.Type() == token.PUNCTUATOR && p.cur.Literal() == "{" {
		p.next() // {
		for p.cur.Literal() != "}" {
			ident := p.expectIdent()
			var expr ast.Expr
			if p.cur.Literal() == "=" {
				p.next()
				expr = p.parseConstantExpr()
			}
			t.List = append(t.List, &ast.EnumField{
				Name: &ast.Ident{Token: ident},
				Val:  expr,
			})
			if p.cur.Literal() == "," {
				p.next() // ,
			}
		}
		p.exceptPunctuator("}") // }
	}
	return t
}

// 特殊限定符
func (p *Parser) scanTypeQualifierTok() (qua []token.Token) {
	for p.cur.Type() == token.KEYWORD && typeQualifierMap[p.cur.Literal()] {
		qua = append(qua, p.cur)
		p.next()
	}
	return
}

func (p *Parser) parseDeclStmt() *ast.DeclStmt {
	stmt := ast.DeclStmt(p.parseDeclaration())
	return &stmt
}

func (p *Parser) parseDeclaration() []ast.Decl {
	typ := p.parseDeclarationSpecifiers()
	var decls []ast.Decl
	for p.cur.Literal() != ";" {
		decl := p.parserInitDeclarator(typ)
		decls = append(decls, decl)
		if p.cur.Literal() == "," {
			p.next() //,
		} else {
			break
		}
	}
	p.exceptPunctuator(";")
	return decls
}

func (p *Parser) parserInitDeclarator(inner ast.TypeName) ast.Decl {
	isTypedef := false
	if v, ok := inner.(*ast.TypeStorageSpecifier); ok && (*v.Specifier)["typedef"] {
		isTypedef = true
		if len(*v.Specifier) == 1 {
			inner = v.Inner
		}
	}
	typ, ident := p.parseDeclarator(inner)
	if isTypedef {
		decl := &ast.TypedefDecl{
			Type: typ,
			Name: ident,
		}
		p.defineType(decl)
		return decl
	}
	decl := &ast.VarDecl{Type: typ, Name: ident}
	if p.cur.Literal() == "=" {
		p.next()
		decl.Init = p.parseInitializer()
	}
	return decl
}

func (p *Parser) defineType(decl *ast.TypedefDecl) {
	p.typeName[decl.Name.Literal()] = decl
}

func (p *Parser) parseStmt() ast.Stmt {
	switch p.cur.Literal() {
	case "case", "default":
		return p.parseLabeledStmt()
	case "if":
		return p.parseIfStmt()
	case "switch":
		return p.parseSwitchStmt()
	case "while":
		return p.parseWhileStmt()
	case "do":
		return p.parseDoWhileStmt()
	case "for":
		return p.parseForStmt()
	case "goto":
		p.next() // goto
		id := p.expectIdent()
		p.exceptPunctuator(";")
		return &ast.GotoStmt{Id: &ast.Ident{Token: id}}
	case "break":
		p.next() // break
		p.exceptPunctuator(";")
		return &ast.BreakStmt{}
	case "continue":
		p.next() // continue
		p.exceptPunctuator(";")
		return &ast.ContinueStmt{}
	case "return":
		p.next() // return
		stmt := &ast.ReturnStmt{}
		if p.cur.Literal() != ";" {
			stmt.X = p.parseExpr()
		}
		p.exceptPunctuator(";")
		return stmt
	case "{":
		return p.parseCompoundStmt()
	}
	if p.cur.Type() == token.IDENT && p.peekOne().Literal() == ":" {
		return p.parseLabeledStmt()
	}
	return p.parseExprStmt()
}

func (p *Parser) parseLabeledStmt() ast.Stmt {
	// case const-expr:
	if p.cur.Literal() == "case" {
		p.next() // case
		expr := p.parseConstantExpr()
		p.exceptPunctuator(":")
		stmt := p.parseStmt()
		return &ast.CaseStmt{
			Expr: expr,
			Stmt: stmt,
		}
	}

	// default:
	if p.cur.Literal() == "default" {
		p.next() // default
		p.exceptPunctuator(":")
		stmt := p.parseStmt()
		return &ast.DefaultStmt{
			Stmt: stmt,
		}
	}

	// id:
	if p.cur.Type() == token.IDENT && p.peekOne().Literal() == ":" {
		ident := p.expectIdent()
		p.exceptPunctuator(":")
		stmt := p.parseStmt()
		return &ast.LabelStmt{
			Id:   &ast.Ident{Token: ident},
			Stmt: stmt,
		}
	}
	return p.parseStmt()
}

func (p *Parser) parseSwitchStmt() ast.Stmt {
	p.next() // switch
	p.exceptPunctuator("(")
	expr := p.parseExpr()
	p.exceptPunctuator(")")
	stmt := p.parseStmt()
	return &ast.SwitchStmt{
		X:    expr,
		Stmt: stmt,
	}
}

func (p *Parser) parseIfStmt() ast.Stmt {
	p.next() // if
	p.exceptPunctuator("(")
	expr := p.parseExpr()
	p.exceptPunctuator(")")
	then := p.parseStmt()
	var elseStmt ast.Stmt
	if p.cur.Literal() == "else" {
		elseStmt = p.parseStmt()
	}
	return &ast.IfStmt{
		X:    expr,
		Then: then,
		Else: elseStmt,
	}
}

func (p *Parser) parseWhileStmt() ast.Stmt {
	p.next() // switch
	p.exceptPunctuator("(")
	expr := p.parseExpr()
	p.exceptPunctuator(")")
	stmt := p.parseStmt()
	return &ast.WhileStmt{
		X:    expr,
		Stmt: stmt,
	}
}

func (p *Parser) parseDoWhileStmt() ast.Stmt {
	p.next() // do
	stmt := p.parseStmt()
	p.exceptKeyword("while")
	p.exceptPunctuator("(")
	expr := p.parseExpr()
	p.exceptPunctuator(")")
	p.exceptPunctuator(";")
	return &ast.DoWhileStmt{
		Stmt: stmt,
		X:    expr,
	}
}

func (p *Parser) parseForStmt() ast.Stmt {
	p.next() // for
	forStmt := &ast.ForStmt{}
	if p.isTypeNameTok(p.cur) {
		forStmt.Decl = p.parseDeclStmt()
	} else {
		forStmt.Init = p.parseExpr()
		p.exceptPunctuator(";")
	}
	if p.cur.Literal() != ";" {
		expr := p.parseExpr()
		forStmt.Cond = expr
	}
	p.exceptPunctuator(";")
	if p.cur.Literal() != ";" {
		expr := p.parseExpr()
		forStmt.Post = expr
	}
	p.exceptPunctuator(";")
	stmt := p.parseStmt()
	forStmt.Stmt = stmt
	return forStmt
}

func (p *Parser) parseCompoundStmt() *ast.CompoundStmt {
	stmts := ast.CompoundStmt{}
	p.exceptPunctuator("{")
	for p.cur.Literal() != "}" {
		stmt := p.parseBlockItem()
		stmts = append(stmts, stmt)
	}
	p.exceptPunctuator("}")
	return &stmts
}

func (p *Parser) parseBlockItem() ast.Stmt {
	if declarationSpecifierMap[p.cur.Literal()] {
		return p.parseDeclStmt()
	}
	return p.parseStmt()
}

func (p *Parser) parseExprStmt() ast.Stmt {
	expr := p.parseExpr()
	p.exceptPunctuator(";")
	return &ast.ExprStmt{Expr: expr}
}

func (p *Parser) parseExternalDecl() ast.Decl {
	typ := p.parseDeclarationSpecifiers()
	isTypedef := false
	if v, ok := typ.(*ast.TypeStorageSpecifier); ok && (*v.Specifier)["typedef"] {
		isTypedef = true
		if len(*v.Specifier) == 1 {
			typ = v.Inner
		}
	}
	typ, ident := p.parseDeclarator(typ)
	if isTypedef {
		decl := &ast.TypedefDecl{
			Type: typ,
			Name: ident,
		}
		p.defineType(decl)
		return decl
	}

	if v, ok := typ.(*ast.FuncType); ok {
		fn := &ast.FuncDecl{
			Name:     ident,
			Return:   v.Inner,
			Params:   v.Params,
			Ellipsis: v.Ellipsis,
		}
		if p.cur.Literal() == ";" {
			p.exceptPunctuator(";")
			return fn
		}

		// 如果函数中未定义参数类型 则在后续定义语句定义
		if len(fn.Params) > 0 && fn.Params[0].Type == nil {
			for declarationSpecifierMap[p.cur.Literal()] {
				fn.Decl = append(fn.Decl, p.parseDeclaration()...)
			}
		}

		if p.cur.Literal() == "{" {
			fn.Body = p.parseCompoundStmt()
		}
		return fn
	}

	decl := &ast.VarDecl{Type: typ, Name: ident}
	if p.cur.Literal() == "=" {
		p.next()
		decl.Init = p.parseInitializer()
	}
	return decl
}

func (p *Parser) ParseUnit() *ast.TranslationUnitDecl {
	decls := &ast.TranslationUnitDecl{}
	for p.cur.Type() != token.EOF {
		decls.Decl = append(decls.Decl, p.ParseDecl())
	}
	return decls
}

func (p *Parser) ParseDecl() ast.Decl {
	return p.parseExternalDecl()
}

func (p *Parser) isTypeNameTok(tok token.Token) bool {
	name := tok.Literal()
	if typeSpecifierQualifierMap[name] {
		return true
	}
	return p.isTypedefName(tok)
}

func (p *Parser) isTypedefName(tok token.Token) bool {
	if _, ok := p.typeName[tok.Literal()]; ok {
		return true
	}
	return false
}

func (p *Parser) expectIdent() token.Token {
	if p.cur.Type() == token.IDENT {
		tok := p.cur
		p.next()
		return tok
	}
	p.addErr(p.cur, errors.ErrSyntaxExpectedIdentGot, p.cur.Literal())
	return p.cur
}

// 获取下一个Token
func (p *Parser) next() token.Token {
	p.cur = p.r.Scan()
	return p.cur
}

func (p *Parser) exceptKeyword(lit string) (t token.Token) {
	t = p.cur
	if p.cur.Type() == token.KEYWORD && p.cur.Literal() == lit {
		p.next()
		return
	}
	p.addErr(p.cur, errors.ErrSyntaxExpectedGot, lit, p.cur.Literal())
	return
}

func (p *Parser) exceptPunctuator(lit string) (t token.Token) {
	t = p.cur
	if p.cur.Type() == token.PUNCTUATOR && p.cur.Literal() == lit {
		p.next()
		return
	}
	p.addErr(p.cur, errors.ErrSyntaxExpectedGot, lit, p.cur.Literal())
	return
}

func (p *Parser) addErr(pos token.Token, code errors.ErrCode, args ...interface{}) {
	p.err(pos.Position(), errors.ErrTypeError, code, args...)
}

func (p *Parser) addWarn(pos token.Token, code errors.ErrCode, args ...interface{}) {
	p.err(pos.Position(), errors.ErrTypeError, code, args...)
}

func (p *Parser) peekOne() token.Token {
	if ps, ok := p.r.(scanner.PeekScanner); ok {
		return ps.PeekOne()
	}
	pp := scanner.NewPeekScan(p.r)
	tok := pp.PeekOne()
	p.r = pp
	return tok
}
