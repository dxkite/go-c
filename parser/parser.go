package parser

import (
	"dxkite.cn/c/ast"
	"dxkite.cn/c/errors"
	"dxkite.cn/c/scanner"
	"dxkite.cn/c/token"
)

// 表达式解析
type parser struct {
	// 错误处理
	err errors.ErrorHandler
	// 当前token
	cur token.Token
	// 当前输入
	r scanner.Scanner
	// 环境
	env *environment
	// 当前文件
	file string
}

type multiparser struct {
	global *ast.Scope // 全局作用域 (extern)
	r      scanner.PeekScanner
	err    errors.ErrorHandler
}

func newMultiparser(r scanner.Scanner, err errors.ErrorHandler) *multiparser {
	return &multiparser{
		r:      scanner.NewPeekScan(r),
		err:    err,
		global: ast.NewScope(ast.GlobalScope, nil, 1),
	}
}

func (p *multiparser) parseUnit() *ast.TranslationUnit {
	unit := &ast.TranslationUnit{}
	for {
		t := p.r.PeekOne()
		if t.Type() == token.EOF {
			break
		}
		file := t.Position().Filename
		pp := newParser(file, p.r, p.global, p.err)
		ret := pp.parseFile()
		p.push(pp.cur)
		unit.Files = append(unit.Files, ret)
	}
	return unit
}

func (p *multiparser) push(tok token.Token) {
	p.r = scanner.NewPeekScan(scanner.NewMultiScan(p.r, scanner.NewArrayScan([]token.Token{tok})))
}

func newParser(file string, r scanner.Scanner, glb *ast.Scope, err errors.ErrorHandler) *parser {
	p := &parser{
		err:  err,
		file: file,
		r:    scanner.NewTokenScan(r),
	}
	p.next()
	p.env = newEnv(glb, p)
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
func (p *parser) parsePrimaryExpr() ast.Expr {
	if p.cur.Type() == token.PUNCTUATOR && p.cur.Literal() == "(" {
		l := p.cur
		p.next()
		x := p.parseExpr()
		r := p.exceptPunctuator(")")
		return &ast.ParenExpr{
			Lparen: l.Position(),
			X:      x,
			Rparen: r.Position(),
		}
	}
	switch p.cur.Type() {
	case token.IDENT:
		cur := p.cur
		p.next()
		ident := &ast.Ident{Token: cur}
		if obj := p.env.resolveIdent(ident); obj != nil {
			ident.Type = obj.Typename
		}
		return ident
	case token.INT, token.CHAR, token.FLOAT, token.STRING:
		cur := p.cur
		p.next()
		return &ast.BasicLit{Token: cur}
	}
	exp := ast.BadExpr{Token: p.cur}
	p.next()
	return &exp
}

func (p *parser) parsePostfixExpr() ast.Expr {
	// ( typename ) { init-list }
	if p.cur.Type() == token.PUNCTUATOR && p.cur.Literal() == "(" {
		if t := p.peekOne(); p.isTypeNameTok(t) {
			return p.parseCompoundLitExpr()
		}
	}
	expr := p.parsePrimaryExpr()
	return p.parsePostfixExprInner(expr)
}

func (p *parser) parsePostfixExprInner(expr ast.Expr) ast.Expr {
	switch p.cur.Literal() {
	case "[":
		lb := p.cur.Position()
		p.next() // [
		idx := p.parseExpr()
		t := p.exceptPunctuator("]")
		expr = &ast.IndexExpr{
			Arr:    expr,
			Lbrack: lb,
			Index:  idx,
			Rbrack: t.Position(),
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
		l := p.cur
		p.next() // (
		arg := p.parseArgsExpr()
		r := p.exceptPunctuator(")")
		expr = &ast.CallExpr{
			Func:   expr,
			Lparen: l.Position(),
			Args:   arg,
			Rparen: r.Position(),
		}
	default:
		return expr
	}
	return p.parsePostfixExprInner(expr)
}

func (p *parser) parseUnaryExpr() ast.Expr {
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

func (p *parser) parseCastExpr() ast.Expr {
	if t := p.peekOne(); p.cur.Literal() == "(" && p.isTypeNameTok(t) {
		lp := p.cur
		p.next()                      // (
		name := p.parseTypeName()     // type-name
		rp := p.exceptPunctuator(")") // )
		if p.cur.Literal() == "{" {
			expr := p.parseInitializerList()
			return &ast.CompoundLit{
				Lparen:   lp.Position(),
				Type:     name,
				Rparen:   rp.Position(),
				InitList: expr,
			}
		}
		expr := p.parseCastExpr()
		return &ast.TypeCastExpr{
			Lparen: lp.Position(),
			Type:   name,
			Rparen: rp.Position(),
			X:      expr,
		}
	}
	return p.parseUnaryExpr()
}

func (p *parser) selectExpr(prec int) (expr ast.Expr) {
	if prec >= UnaryPrec {
		return p.parseCastExpr()
	} else {
		return p.parseBinaryExpr(prec)
	}
}

func (p *parser) parseBinaryExpr(prec int) ast.Expr {
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

func (p *parser) parseCondExpr() ast.Expr {
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
func (p *parser) parseConstantExpr() ast.Expr {
	return &ast.ConstantExpr{X: p.parseCondExpr()}
}

func (p *parser) parseAssignExpr() ast.Expr {
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

func (p *parser) parseExpr() ast.Expr {
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

func (p *parser) parseCompoundLitExpr() ast.Expr {
	lp := p.cur
	p.next() // (
	typeName := p.parseTypeName()
	rp := p.exceptPunctuator(")")
	expr := p.parseInitializerList()
	return &ast.CompoundLit{
		Lparen:   lp.Position(),
		Type:     typeName,
		Rparen:   rp.Position(),
		InitList: expr,
	}
}

func (p *parser) parseInitializerList() *ast.InitializerExpr {
	l := p.exceptPunctuator("{")
	expr := ast.InitializerExpr{}
	expr.Lbrace = l.Position()
	for p.until("}") {
		item := p.parseInitializer()
		expr.List = append(expr.List, item)
		if t := p.peekOne(); p.cur.Literal() != "}" && t.Literal() != "}" {
			p.exceptPunctuator(",")
		} else {
			if p.cur.Literal() == "," {
				p.next() //,
			}
		}
	}
	r := p.exceptPunctuator("}")
	expr.Rbrace = r.Position()
	return &expr
}

func (p *parser) parseInitializer() ast.Expr {
	if p.cur.Literal() == "{" {
		return p.parseInitializerList()
	}
	if l := p.cur.Literal(); l == "." || l == "[" {
		return p.parseDesignationInitExpr()
	}
	return p.parseAssignExpr()
}

func (p *parser) parseDesignationInitExpr() ast.Expr {
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

func (p *parser) parseDesignator() ast.Expr {
	bg := p.cur.Position()
	// .xxx
	if p.cur.Literal() == "." {
		p.next()
		ident := p.expectIdent()
		x := p.parseDesignator()
		return &ast.RecordDesignatorExpr{
			Dot:   bg,
			Field: &ast.Ident{Token: ident},
			X:     x,
		}
	}
	if p.cur.Literal() == "[" {
		p.next()
		expr := p.parseConstantExpr()
		rb := p.exceptPunctuator("]")
		x := p.parseDesignator()
		return &ast.ArrayDesignatorExpr{
			Lbrack: bg,
			Index:  expr,
			Rbrack: rb.Position(),
			X:      x,
		}
	}
	return nil
}

func (p *parser) parseArgsExpr() []ast.Expr {
	var list []ast.Expr
	for p.until(")") {
		item := p.parseAssignExpr()
		list = append(list, item)
		if p.cur.Literal() != "," {
			break
		}
		p.exceptPunctuator(",")
	}
	return list
}

var typeQualifier = []string{"const", "restrict", "volatile"}
var typeSpecifier = []string{"void", "char", "short", "int", "long", "float", "double", "signed", "unsigned", "_Bool" /*"_Complex"*/}
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
func (p *parser) parseTypeName() ast.Typename {
	basic := p.parseTypeQualifierSpecifierList()
	return p.parseAbstractDeclarator(basic)
}

func (p *parser) parseAbstractDeclarator(inner ast.Typename) ast.Typename {
	switch p.cur.Literal() {
	case "*":
		inner = p.parsePointer(inner)
	case "(":
		if t := p.peekOne().Literal(); t == "*" || t == "[" || t == "(" {
			p.next() // (
			inner = p.parseAbstractDeclarator(&ast.ParenType{Type: inner})
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

func (p *parser) parseFuncType(inner ast.Typename) ast.Typename {
	p.env.enterScope(ast.FuncPrototypeScope)
	lp := p.exceptPunctuator("(")
	params, ellipsis := p.parseParameterList()
	rp := p.exceptPunctuator(")")
	p.env.leaveScope()
	return &ast.FuncType{
		Return:   inner,
		Lparen:   lp.Position(),
		Params:   params,
		Ellipsis: ellipsis,
		Rparen:   rp.Position(),
	}
}

func (p *parser) parseArrayType(inner ast.Typename) ast.Typename {
	if t := p.peekOne(); t.Literal() == "*" || t.Literal() == "]" {
		lb := p.exceptPunctuator("[")
		if p.cur.Literal() == "*" {
			p.next() // *
		}
		rb := p.exceptPunctuator("]")
		return &ast.ArrayType{
			Type:       inner,
			Incomplete: true,
			Lbrack:     lb.Position(),
			Rbrack:     rb.Position(),
		}
	}
	return p.parseArrayTypeSize(inner)
}

func (p *parser) parseArrayTypeSize(inner ast.Typename) ast.Typename {
	lb := p.exceptPunctuator("[") // [
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
	rb := p.exceptPunctuator("]") // ]
	arr := &ast.ArrayType{
		Type:   inner,
		Static: static,
		Lbrack: lb.Position(),
		Rbrack: rb.Position(),
	}

	if v, ok := expr.(*ast.AssignExpr); ok {
		arr.Size = v
		return p.makeTypeQualifier(arr, qua)
	}

	// 常量表达式
	arr.Const = true
	return arr
}

func (p *parser) parseParameterList() (params ast.ParamList, ellipsis bool) {
	for p.cur.Literal() != ")" && p.cur.Literal() != "..." && p.cur.Type() != token.EOF {
		param := p.parseParameterDecl()
		p.env.declareIdent(ast.ObjectParamVal, param)
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
	return
}

func (p *parser) parseParameterDecl() *ast.ParamVarDecl {
	typ, spec := p.parseDeclarationSpecifiers()
	param := &ast.ParamVarDecl{Qua: spec}
	param.Type, param.Name = p.parseDeclarator(typ)
	return param
}

func (p *parser) parseDeclarator(inner ast.Typename) (ast.Typename, *ast.Ident) {
	if p.cur.Literal() == "*" {
		inner = p.parsePointer(inner)
	}
	return p.parseDirectDeclarator(inner)
}

func (p *parser) parseDirectDeclarator(inner ast.Typename) (ast.Typename, *ast.Ident) {
	var ident *ast.Ident
	if p.cur.Literal() == "(" {
		lp := p.exceptPunctuator("(")
		typ, ident := p.parseDeclarator(inner)
		rp := p.exceptPunctuator(")")
		typ = &ast.ParenType{Lparen: lp.Position(), Type: typ, Rparen: rp.Position()}
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

func (p *parser) parseDirectDeclaratorInner(typ ast.Typename) ast.Typename {
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
func (p *parser) parseTypeQualifierSpecifierList() ast.Typename {
	var qua []token.Token
	var typ ast.Typename
	var buildIn []token.Token

	for typeSpecifierQualifierMap[p.cur.Literal()] {
		if typeQualifierMap[p.cur.Literal()] {
			qua = append(qua, p.cur)
			p.next()
			continue
		}
		if len(buildIn) > 0 && typeStructMap[p.cur.Literal()] {
			p.addErr(p.cur.Position(), errors.ErrSyntaxUnexpectedTypeSpecifier, p.cur.Literal())
		}
		if t, s := p.parseTypeSpecifier(); t != nil {
			typ = t
		} else {
			buildIn = append(buildIn, s...)
		}
	}

	if len(buildIn) != 0 {
		tp := &ast.BuildInType{
			Type: ast.Int,
		}
		if t, err := ast.ParseBuildInType(buildIn); err != nil {
			p.addErr(err.Pos, err.Code, err.Params...)
		} else {
			tp.Type = t
		}
		rg := &ast.Range{}
		rg.Begin = buildIn[0].Position()
		rg.End = buildIn[len(buildIn)-1].Position()
		tp.Range = rg
		typ = tp
	}

	if len(qua) > 0 && typ != nil {
		p.markQualifier(typ.Qualifier(), qua)
	}
	return typ
}

// (('*') typeQualifierList?)+
func (p *parser) parsePointer(inner ast.Typename) (t ast.Typename) {
	pk := p.exceptPunctuator("*")
	tt := &ast.PointerType{Pointer: pk.Position(), Type: inner}
	tks := p.scanTypeQualifierTok()
	t = p.makeTypeQualifier(tt, tks)
	for p.cur.Literal() == "*" {
		t = p.parsePointer(t)
	}
	return t
}

func (p *parser) makeTypeQualifier(typ ast.Typename, qua []token.Token) ast.Typename {
	p.markQualifier(typ.Qualifier(), qua)
	return typ
}

func (p *parser) isDeclarationSpecifier(tok token.Token) bool {
	return declarationSpecifierMap[tok.Literal()] || p.isTypeNameTok(tok)
}

// 扫描类型
func (p *parser) parseDeclarationSpecifiers() (ast.Typename, *ast.StorageSpecifier) {
	var qua []token.Token
	var typ ast.Typename
	var buildIn []token.Token
	var spec []token.Token

	for p.isDeclarationSpecifier(p.cur) {
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
		if len(buildIn) > 0 && typeStructMap[p.cur.Literal()] {
			p.addErr(p.cur.Position(), errors.ErrSyntaxUnexpectedTypeSpecifier, p.cur.Literal())
		}
		if t, s := p.parseTypeSpecifier(); t != nil {
			typ = t
		} else {
			buildIn = append(buildIn, s...)
		}
	}

	if len(buildIn) != 0 {
		tp := &ast.BuildInType{
			Type: ast.Int,
		}
		if t, err := ast.ParseBuildInType(buildIn); err != nil {
			p.addErr(err.Pos, err.Code, err.Params...)
		} else {
			tp.Type = t
		}
		typ = tp
	}

	if len(qua) > 0 && typ != nil {
		p.markQualifier(typ.Qualifier(), qua)
	}

	storage := &ast.StorageSpecifier{}
	if len(spec) > 0 {
		p.markSpecifier(storage, spec)
	}
	return typ, storage
}

func (p *parser) markSpecifier(q *ast.StorageSpecifier, qua []token.Token) {
	for _, t := range qua {
		if _, ok := (*q)[t.Literal()]; ok {
			p.addWarn(t.Position(), errors.ErrSyntaxDuplicateTypeSpecifier, t.Literal())
		}
		(*q)[t.Literal()] = t.Position()
	}
}

func (p *parser) parseTypeSpecifier() (ast.Typename, []token.Token) {
	switch p.cur.Literal() {
	case "struct", "union":
		return p.parseRecordType(), nil
	case "enum":
		return p.parseEnumType(), nil
	default:
		// 用户定义的类型
		if p.cur.Type() == token.IDENT {
			if t := p.isTypedefName(p.cur); t != nil {
				p.next()
				return t, nil
			}
		}
	}
	return nil, p.parseBuildInSpec()
}

// 扫描内置类型
func (p *parser) parseBuildInSpec() []token.Token {
	var spec []token.Token
	for p.cur.Type() == token.KEYWORD && typeSpecifierMap[p.cur.Literal()] {
		if typeSpecifierMap[p.cur.Literal()] {
			spec = append(spec, p.cur)
		}
		p.next()
	}
	return spec
}

func (p *parser) markQualifier(q *ast.Qualifier, qua []token.Token) {
	if q == nil {
		return
	}
	for _, t := range qua {
		if _, ok := (*q)[t.Literal()]; ok {
			p.addWarn(t.Position(), errors.ErrSyntaxDuplicateTypeQualifier, t.Literal())
		}
		(*q)[t.Literal()] = t.Position()
	}
}

func (p *parser) parseRecordType() *ast.RecordType {
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

	p.env.declareRecord(r, false)

	r.Completed = true
	r.Lbrace = p.exceptPunctuator("{").Position()
	for p.until("}") {
		typ := p.parseTypeQualifierSpecifierList()
		for p.cur.Type() != token.EOF {
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
			if f.Bit == nil && f.Name == nil && !isRecordType(typ) {
				p.addErr(p.cur.Position(), errors.ErrSyntaxExpectedRecordMemberName)
				break
			}
			p.env.resolveType(f.Type, true)
			r.Fields = append(r.Fields, f)
			if p.cur.Literal() != "," {
				break
			}
			p.exceptPunctuator(",")
		}
		p.exceptPunctuator(";")
	}

	r.Rbrace = p.exceptPunctuator("}").Position() // }
	p.env.declareRecord(r, true)
	return r
}

func isRecordType(typ ast.Typename) bool {
	switch v := typ.(type) {
	case *ast.RecordType:
		return true
	case *ast.PointerType:
		return isRecordType(v.Type)
	}
	return false
}

func (p *parser) parseEnumType() *ast.EnumType {
	pk := p.exceptKeyword("enum")
	t := &ast.EnumType{Enum: pk.Position()}
	if p.cur.Type() == token.IDENT {
		t.Name = &ast.Ident{Token: p.cur}
		p.next()
	}
	if p.cur.Literal() != "{" {
		return t
	}
	p.env.declareEnum(t, false)

	t.Lbrace = p.exceptPunctuator("{").Position()
	t.Completed = true
	for p.until("}") {
		ident := p.expectIdent()
		var expr ast.Expr
		if p.cur.Literal() == "=" {
			p.next()
			expr = p.parseConstantExpr()
		}
		tag := &ast.EnumFieldDecl{
			Name: &ast.Ident{Token: ident},
			Val:  expr,
		}
		p.env.declareEnumTag(tag)
		t.List = append(t.List, tag)
		if p.cur.Literal() != "," {
			break
		}
		p.exceptPunctuator(",")
	}

	t.Rbrace = p.exceptPunctuator("}").Position()
	p.env.declareEnum(t, true)
	return t
}

// 特殊限定符
func (p *parser) scanTypeQualifierTok() (qua []token.Token) {
	for p.cur.Type() == token.KEYWORD && typeQualifierMap[p.cur.Literal()] {
		qua = append(qua, p.cur)
		p.next()
	}
	return
}

func (p *parser) parseDeclStmt() *ast.DeclStmt {
	stmt := ast.DeclStmt(p.parseDeclaration())
	return &stmt
}

func (p *parser) parseDeclaration() []ast.Decl {
	typ, spec := p.parseDeclarationSpecifiers()
	var decls []ast.Decl
	for p.until(";") {
		decl := p.parserInitDeclarator(typ, spec)
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

func (p *parser) until(lit string) bool {
	return p.cur.Literal() != lit && p.cur.Type() != token.EOF
}

func (p *parser) parserInitDeclarator(inner ast.Typename, spec *ast.StorageSpecifier) ast.Decl {
	decl, _ := p.parseExternalDeclOrInit(inner, spec, false)
	return decl
}

func (p *parser) parseStmt() ast.Stmt {
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
		pk := p.exceptKeyword("goto")
		id := p.expectIdent()
		se := p.exceptPunctuator(";")
		stmt := &ast.GotoStmt{
			Goto:      pk.Position(),
			Id:        &ast.Ident{Token: id},
			Semicolon: se.Position(),
		}
		p.env.tryResolveLabel(stmt.Id)
		return stmt
	case "break":
		pk := p.exceptKeyword("break")
		se := p.exceptPunctuator(";")
		return &ast.BreakStmt{Break: pk.Position(), Semicolon: se.Position()}
	case "continue":
		pk := p.exceptKeyword("continue")
		se := p.exceptPunctuator(";")
		return &ast.ContinueStmt{Continue: pk.Position(), Semicolon: se.Position()}
	case "return":
		pk := p.exceptKeyword("return")
		stmt := &ast.ReturnStmt{Return: pk.Position()}
		if p.cur.Literal() != ";" {
			stmt.X = p.parseExpr()
		}
		se := p.exceptPunctuator(";")
		stmt.Semicolon = se.Position()
		return stmt
	case "{":
		stmt := p.parseCompoundStmt()
		return stmt
	}
	if p.cur.Type() == token.IDENT && p.peekOne().Literal() == ":" {
		return p.parseLabeledStmt()
	}
	return p.parseExprStmt()
}

func (p *parser) parseLabeledStmt() ast.Stmt {
	// case const-expr:
	if p.cur.Literal() == "case" {
		pk := p.exceptKeyword("case")
		expr := p.parseConstantExpr()
		p.exceptPunctuator(":")
		stmt := p.parseStmt()
		return &ast.CaseStmt{
			Case: pk.Position(),
			Expr: expr,
			Stmt: stmt,
		}
	}

	// default:
	if p.cur.Literal() == "default" {
		pk := p.exceptKeyword("default")
		p.exceptPunctuator(":")
		stmt := p.parseStmt()
		return &ast.DefaultStmt{
			Default: pk.Position(),
			Stmt:    stmt,
		}
	}

	// id:
	if p.cur.Type() == token.IDENT && p.peekOne().Literal() == ":" {
		ident := p.expectIdent()
		p.exceptPunctuator(":")
		stmt := p.parseStmt()
		st := &ast.LabelStmt{
			Id:   &ast.Ident{Token: ident},
			Stmt: stmt,
		}
		p.env.declare(ast.NewObject(ast.ObjectLabelName, st.Id))
		return st
	}
	return p.parseStmt()
}

func (p *parser) parseSwitchStmt() ast.Stmt {
	pk := p.exceptKeyword("switch")
	p.exceptPunctuator("(")
	expr := p.parseExpr()
	p.exceptPunctuator(")")
	stmt := p.parseStmt()
	return &ast.SwitchStmt{
		Switch: pk.Position(),
		X:      expr,
		Stmt:   stmt,
	}
}

func (p *parser) parseIfStmt() ast.Stmt {
	pk := p.exceptKeyword("if")
	p.exceptPunctuator("(")
	expr := p.parseExpr()
	p.exceptPunctuator(")")
	then := p.parseStmt()
	var elseStmt ast.Stmt
	if p.cur.Literal() == "else" {
		elseStmt = p.parseStmt()
	}
	return &ast.IfStmt{
		If:   pk.Position(),
		X:    expr,
		Then: then,
		Else: elseStmt,
	}
}

func (p *parser) parseWhileStmt() ast.Stmt {
	pk := p.exceptKeyword("while")
	p.exceptPunctuator("(")
	expr := p.parseExpr()
	p.exceptPunctuator(")")
	stmt := p.parseStmt()
	return &ast.WhileStmt{
		While: pk.Position(),
		X:     expr,
		Stmt:  stmt,
	}
}

func (p *parser) parseDoWhileStmt() ast.Stmt {
	pk := p.exceptKeyword("do")
	stmt := p.parseStmt()
	p.exceptKeyword("while")
	p.exceptPunctuator("(")
	expr := p.parseExpr()
	p.exceptPunctuator(")")
	se := p.exceptPunctuator(";")
	return &ast.DoWhileStmt{
		Do:        pk.Position(),
		Stmt:      stmt,
		X:         expr,
		Semicolon: se.Position(),
	}
}

func (p *parser) parseForStmt() ast.Stmt {
	pk := p.exceptKeyword("for")
	forStmt := &ast.ForStmt{For: pk.Position()}
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

func (p *parser) parseCompoundStmt() *ast.CompoundStmt {
	comp := ast.CompoundStmt{}
	lb := p.exceptPunctuator("{")
	for p.until("}") {
		stmt := p.parseBlockItem()
		comp.Stmts = append(comp.Stmts, stmt)
	}
	rb := p.exceptPunctuator("}")
	comp.Lbrace = lb.Position()
	comp.Rbrace = rb.Position()
	return &comp
}

func (p *parser) parseBlockItem() ast.Stmt {
	if declarationSpecifierMap[p.cur.Literal()] || p.isTypeNameTok(p.cur) {
		return p.parseDeclStmt()
	}
	return p.parseStmt()
}

func (p *parser) parseExprStmt() ast.Stmt {
	expr := p.parseExpr()
	sm := p.exceptPunctuator(";")
	return &ast.ExprStmt{Expr: expr, Semicolon: sm.Position()}
}

func (p *parser) parseExternalDecl() ast.Decl {
	typ, spec := p.parseDeclarationSpecifiers()
	decl, comma := p.parseExternalDeclOrInit(typ, spec, true)
	if comma {
		p.exceptPunctuator(";")
	}
	return decl
}

func (p *parser) parseExternalDeclOrInit(inner ast.Typename, specifier *ast.StorageSpecifier, external bool) (ast.Decl, bool) {
	isTypedef := false
	if _, ok := (*specifier)["typedef"]; ok {
		isTypedef = true
	}

	typ, ident := p.parseDeclarator(inner)
	if isTypedef {
		decl := &ast.TypedefDecl{
			Typedef: (*specifier)["typedef"],
			Type:    typ,
			Name:    ident,
		}
		p.env.declareType(decl)
		return decl, external
	}

	if v, ok := typ.(*ast.FuncType); ok && external {
		fn := &ast.FuncDecl{
			Type: v,
			Name: ident,
		}

		obj := ast.NewDeclObject(ast.ObjectFunc, ident, fn)
		if p.cur.Literal() == ";" {
			p.env.declare(obj)
			return fn, external
		}

		// 如果函数中未定义参数类型 则在后续尝试解析语句定义
		if len(v.Params) > 0 && v.Params[0].Type == nil {
			for declarationSpecifierMap[p.cur.Literal()] {
				fn.Decl = append(fn.Decl, p.parseDeclaration()...)
			}
		}

		p.env.enterLabelScope()
		p.env.enterScope(ast.FuncScope)
		p.env.copyParameterScope(fn)
		fn.Body = p.parseCompoundStmt()
		p.env.leaveScope()
		p.reportUnResolveLabel(p.env.leaveLabelScope())

		obj.Completed = true
		p.env.declare(obj)
		return fn, false
	}

	decl := &ast.VarDecl{Type: typ, Name: ident}

	if p.cur.Literal() == "=" {
		p.exceptPunctuator("=")
		decl.Init = p.parseInitializer()
	}

	p.env.declareIdent(ast.ObjectVar, decl)
	return decl, external
}

func (p *parser) parseFile() *ast.File {
	unit := &ast.File{}
	unit.Name = p.file
	var decls []ast.Decl
	for p.cur.Type() != token.EOF && p.cur.Position().Filename == p.file {
		if p.isDeclarationSpecifier(p.cur) {
			decl := p.parseDecl()
			decls = append(decls, decl)
		} else {
			p.addErr(p.cur.Position(), errors.ErrSyntaxExpectedGot, "定义语句", p.cur.Literal())
			p.next()
		}
	}
	unit.Decl = decls
	return unit
}

func (p *parser) parseDecl() ast.Decl {
	return p.parseExternalDecl()
}

func (p *parser) isTypeNameTok(tok token.Token) bool {
	name := tok.Literal()
	if typeSpecifierQualifierMap[name] {
		return true
	}
	return p.isTypedefName(tok) != nil
}

func (p *parser) isTypedefName(tok token.Token) ast.Typename {
	return p.env.isTypename(tok.Literal())
}

func (p *parser) expectIdent() token.Token {
	if p.cur.Type() == token.IDENT {
		tok := p.cur
		p.next()
		return tok
	}
	p.addErr(p.cur.Position(), errors.ErrSyntaxExpectedIdentGot, p.cur.Literal())
	return p.cur
}

// 获取下一个Token
func (p *parser) next() token.Token {
	p.cur = p.r.Scan()
	return p.cur
}

func (p *parser) exceptKeyword(lit string) (t token.Token) {
	t = p.cur
	if p.cur.Type() == token.KEYWORD && p.cur.Literal() == lit {
		p.next()
		return
	}
	p.addErr(p.cur.Position(), errors.ErrSyntaxExpectedGot, lit, p.cur.Literal())
	return
}

func (p *parser) exceptPunctuator(lit string) (t token.Token) {
	t = p.cur
	if p.cur.Type() == token.PUNCTUATOR && p.cur.Literal() == lit {
		p.next()
		return
	}
	p.addErr(p.cur.Position(), errors.ErrSyntaxExpectedGot, lit, p.cur.Literal())
	return
}

func (p *parser) addErr(pos token.Position, code errors.ErrCode, args ...interface{}) {
	p.err(pos, errors.ErrTypeError, code, args...)
}

func (p *parser) addWarn(pos token.Position, code errors.ErrCode, args ...interface{}) {
	p.err(pos, errors.ErrTypeError, code, args...)
}

func (p *parser) reportUnResolveLabel(labels []*ast.Ident) {
	for _, v := range labels {
		p.addErr(v.Position(), errors.ErrSyntaxUndefinedLabel, v.Literal())
	}
}

func (p *parser) peekOne() token.Token {
	if ps, ok := p.r.(scanner.PeekScanner); ok {
		return ps.PeekOne()
	}
	pp := scanner.NewPeekScan(p.r)
	tok := pp.PeekOne()
	p.r = pp
	return tok
}
