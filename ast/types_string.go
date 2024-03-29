// Code generated by "stringer -type BasicType -linecomment"; DO NOT EDIT.

package ast

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[UnknownType-0]
	_ = x[Void-1]
	_ = x[Bool-2]
	_ = x[Char-3]
	_ = x[Short-4]
	_ = x[Int-5]
	_ = x[LongLong-6]
	_ = x[Float-7]
	_ = x[Double-8]
	_ = x[UnsignedChar-9]
	_ = x[UnsignedShort-10]
	_ = x[UnsignedInt-11]
	_ = x[UnsignedLongLong-12]
}

const _BuildInType_name = "UnknownTypevoid_Boolcharshortintlong longfloatdoubleunsigned charunsigned shortunsigned intunsigned long long"

var _BuildInType_index = [...]uint8{0, 11, 15, 20, 24, 29, 32, 41, 46, 52, 65, 79, 91, 109}

func (i BasicType) String() string {
	if i < 0 || i >= BasicType(len(_BuildInType_index)-1) {
		return "BasicType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _BuildInType_name[_BuildInType_index[i]:_BuildInType_index[i+1]]
}
