//go:generate stringer -type CBuildInType -linecomment
package ast

import (
	"dxkite.cn/c/errors"
	"dxkite.cn/c/token"
)

// C语言基础类型
type CBuildInType int

const (
	CUnknownType CBuildInType = iota
	CVoid                     // void
	// bool
	CBool // _Bool
	// int8 char
	CChar // char
	// int16 short
	CShort // short
	// int32 int/long
	CInt // int
	// int64 long long
	CLongLong // long long
	// float32 float
	CFloat // float
	// float64
	CDouble // double
	// uint8
	CUnsignedChar // unsigned char
	// uint16
	CUnsignedShort // unsigned short
	// uint32
	CUnsignedInt // unsigned int
	// uint64
	CUnsignedLongLong // unsigned long long
)

func ParseBuildInType(tks []token.Token) (CBuildInType, *errors.Error) {
	var base CBuildInType
	long := 0
	signed := false
	unsigned := false

	for _, v := range tks {
		switch v.Literal() {
		case "void", "_Bool", "char", "short", "float", "double":
			if base != CUnknownType {
				return CInt, errors.New(v.Position(), errors.ErrSyntaxUnexpectedTypeSpecifier, v.Literal())
			}
		}
		switch v.Literal() {
		case "void":
			if unsigned || signed {
				return CInt, errors.New(v.Position(), errors.ErrSyntaxUnexpectedTypeSpecifier, v.Literal())
			}
			base = CVoid
		case "_Bool":
			if unsigned || signed {
				return CInt, errors.New(v.Position(), errors.ErrSyntaxUnexpectedTypeSpecifier, v.Literal())
			}
			base = CBool
		case "char":
			base = CChar
		case "int":
			base = CInt
		case "short":
			base = CShort
		case "float":
			base = CFloat
		case "double":
			base = CDouble
		case "unsigned":
			if unsigned || signed {
				return CInt, errors.New(v.Position(), errors.ErrSyntaxUnexpectedTypeSpecifier, v.Literal())
			}
			unsigned = true
		case "signed":
			if unsigned || signed {
				return CInt, errors.New(v.Position(), errors.ErrSyntaxUnexpectedTypeSpecifier, v.Literal())
			}
			signed = true
		case "long":
			long++
			if long > 2 {
				return CInt, errors.New(v.Position(), errors.ErrSyntaxUnexpectedTypeSpecifier, v.Literal())
			}
			if long == 2 {
				base = CLongLong
			}
		default:
			return CInt, errors.New(v.Position(), errors.ErrSyntaxUnexpectedTypeSpecifier, v.Literal())
		}
	}
	// 无符号数据
	if unsigned && (base >= CChar && base <= CLongLong) {
		base += CUnsignedChar - CChar
	}
	return base, nil // 默认类型int
}
