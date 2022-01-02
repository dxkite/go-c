//go:generate stringer -type BuildInType -linecomment
package ast

import (
	"dxkite.cn/c/errors"
	"dxkite.cn/c/token"
)

// C语言基础类型
type BasicType int

const (
	UnknownType BasicType = iota

	Void // void
	// bool
	Bool // _Bool
	// int8 char
	Char // char
	// int16 short
	Short // short
	// int32 int/long
	Int // int
	// int64 long long
	LongLong // long long
	// float32 float
	Float // float
	// float64
	Double // double
	// uint8
	UnsignedChar // unsigned char
	// uint16
	UnsignedShort // unsigned short
	// uint32
	UnsignedInt // unsigned int
	// uint64
	UnsignedLongLong // unsigned long long
	// uintptr
	UnsignedPointer // 无符号指针
)

var sizeof = map[BasicType]int{
	Void:             1,
	Bool:             1,
	Char:             1,
	Short:            2,
	Int:              4,
	LongLong:         8,
	Float:            4,
	Double:           8,
	UnsignedChar:     1,
	UnsignedShort:    2,
	UnsignedInt:      4,
	UnsignedLongLong: 8,
	UnsignedPointer:  4,
}

func (t BasicType) Size() int {
	if v, ok := sizeof[t]; ok {
		return v
	}
	return 4
}

// 解析内置类型
func ParseBuildInType(lit []token.Token) (BasicType, *errors.Error) {
	var base BasicType
	long := 0
	signed := false
	unsigned := false

	for _, v := range lit {
		switch v.Literal() {
		case "void", "_Bool", "char", "short", "float", "double":
			if base != UnknownType {
				return Int, errors.New(v.Position(), errors.ErrSyntaxUnexpectedTypeSpecifier, v.Literal())
			}
		}
		switch v.Literal() {
		case "void":
			if unsigned || signed {
				return Int, errors.New(v.Position(), errors.ErrSyntaxUnexpectedTypeSpecifier, v.Literal())
			}
			base = Void
		case "_Bool":
			if unsigned || signed {
				return Int, errors.New(v.Position(), errors.ErrSyntaxUnexpectedTypeSpecifier, v.Literal())
			}
			base = Bool
		case "char":
			base = Char
		case "int":
			base = Int
		case "short":
			base = Short
		case "float":
			base = Float
		case "double":
			base = Double
		case "unsigned":
			if unsigned || signed {
				return Int, errors.New(v.Position(), errors.ErrSyntaxUnexpectedTypeSpecifier, v.Literal())
			}
			unsigned = true
		case "signed":
			if unsigned || signed {
				return Int, errors.New(v.Position(), errors.ErrSyntaxUnexpectedTypeSpecifier, v.Literal())
			}
			signed = true
		case "long":
			long++
			if long > 2 {
				return Int, errors.New(v.Position(), errors.ErrSyntaxUnexpectedTypeSpecifier, v.Literal())
			}
			if long == 2 {
				base = LongLong
			}
		default:
			return Int, errors.New(v.Position(), errors.ErrSyntaxUnexpectedTypeSpecifier, v.Literal())
		}
	}
	// 无符号数据
	if unsigned && (base >= Char && base <= LongLong) {
		base += UnsignedChar - Char
	}
	return base, nil // 默认类型int
}
