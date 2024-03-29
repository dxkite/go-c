// Code generated by "stringer -type ErrCode -linecomment"; DO NOT EDIT.

package errors

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ErrUnKnown-0]
	_ = x[ErrReadFile-1]
	_ = x[scanErr-1002]
	_ = x[ErrScanUncloseChar-1003]
	_ = x[ErrScanUncloseString-1004]
	_ = x[ErrScanUncloseComment-1005]
	_ = x[ErrScanHexFormat-1006]
	_ = x[ErrScanUnicodeFormat-1007]
	_ = x[macroErr-2008]
	_ = x[ErrMacroHashHashPos-2009]
	_ = x[ErrMacroHashHashExpr-2010]
	_ = x[ErrMacroHashExpr-2011]
	_ = x[ErrMacroCallParamCount-2012]
	_ = x[ErrMacroUnexpectedElseIf-2013]
	_ = x[ErrMacroUnexpectedElse-2014]
	_ = x[ErrMacroUnexpectedEndIf-2015]
	_ = x[ErrMacroExpectedIdent-2016]
	_ = x[ErrMacroExpectedGot-2017]
	_ = x[ErrMacroExpectedPunctuator-2018]
	_ = x[ErrMacroEnd-2019]
	_ = x[ErrMacroExpectedTokenGotEof-2020]
	_ = x[ErrMacroConstExpr-2021]
	_ = x[ErrMacroDuplicateIdent-2022]
	_ = x[ErrMacroInvalidIncludeString-2023]
	_ = x[ErrMacroInvalidIncludeMacro-2024]
	_ = x[ErrMacroIncludeFileRead-2025]
	_ = x[ErrMacroIncludeFileNoFound-2026]
	_ = x[ErrMacroExprUnexpectedToken-2027]
	_ = x[syntaxError-3028]
	_ = x[ErrSyntaxExpectedGot-3029]
	_ = x[ErrSyntaxExpectedIdentGot-3030]
	_ = x[ErrSyntaxUnexpectedTypeSpecifier-3031]
	_ = x[ErrSyntaxDuplicateTypeSpecifier-3032]
	_ = x[ErrSyntaxDuplicateTypeQualifier-3033]
	_ = x[ErrSyntaxExpectedRecordMemberName-3034]
	_ = x[ErrSyntaxRedefineFunc-3035]
	_ = x[ErrSyntaxRedefineVar-3036]
	_ = x[ErrSyntaxRedefineIdent-3037]
	_ = x[ErrSyntaxRedefinedType-3038]
	_ = x[ErrSyntaxRedefinedStruct-3039]
	_ = x[ErrSyntaxRedefinedUnion-3040]
	_ = x[ErrSyntaxRedefinedEnum-3041]
	_ = x[ErrSyntaxRedefinedLabel-3042]
	_ = x[ErrSyntaxUndefinedIdent-3043]
	_ = x[ErrSyntaxUndefinedLabel-3044]
	_ = x[ErrSyntaxIncompleteStruct-3045]
	_ = x[ErrSyntaxIncompleteUnion-3046]
}

const (
	_ErrCode_name_0 = "未知错误代码文件读取失败"
	_ErrCode_name_1 = "scanErr字符缺少关闭的 ' 符号字符串缺少关闭的 \" 符号多行注释缺少对应的关闭 */ 符号符号 %c 不是一个16进制编码字符符号 %c 不是一个Unicode编码字符"
	_ErrCode_name_2 = "macroErr## 不能出现在宏表达式的起始或结束位置## 不能用来连接 %s 和 %s# 符号后面必须跟着一个宏参数宏调用参数数量错误，支持%d个参数，使用了%d个参数不应该出现的 #elif 宏不应该出现的 #else 宏不应该出现的 #endif 宏这里应该是一个名称，不应该出现 %s 符号这里应该是一个 %s ，不应该出现 %s这里应该是一个 %s 符号，不应该出现 %s 符号这里应该是宏结尾了，不应该出现 %s 符号需要符号为 %s，意外的遇到了文件尾错误的宏常量表达式 %s重复定义了符号 %s#include 包含错误的字符串 %s错误的 #include 宏#include的文件 %s 读取错误 %s#include的文件不存在 %s非预期的宏表达式符号%s"
	_ErrCode_name_3 = "syntaxError这里应该是一个 %s ，不应该出现 %s这里应该是一个名称，不应该出现 %s 符号非预期的类型定义符号 %s重复的类型定义符号 %s重复的类型修饰符号 %s类型定义符号之后应该是成员变量的名称重复声明函数 %s，上次声明的位置 %s重复声明的变量名 %s，上次声明的位置 %s重复的标识符 %s，上次声明的位置 %s重复定义的类型 %s，上次定义的位置 %s重复定义的结构体 %s，上次定义的位置 %s重复定义的联合体 %s，上次定义的位置 %s重复定义的枚举 %s，上次定义的位置 %s重复定义的标签 %s，上次定义的位置 %s未定义的标识符 %s未定义的标签 %s不完全的结构体类型 %s不完全的联合体类型 %s"
)

var (
	_ErrCode_index_0 = [...]uint8{0, 12, 36}
	_ErrCode_index_1 = [...]uint8{0, 7, 37, 70, 113, 155, 196}
	_ErrCode_index_2 = [...]uint16{0, 8, 62, 93, 134, 204, 232, 260, 289, 344, 390, 449, 504, 552, 582, 606, 642, 664, 700, 729, 761}
	_ErrCode_index_3 = [...]uint16{0, 11, 57, 112, 145, 175, 205, 259, 307, 361, 409, 460, 514, 568, 619, 670, 694, 715, 745, 775}
)

func (i ErrCode) String() string {
	switch {
	case 0 <= i && i <= 1:
		return _ErrCode_name_0[_ErrCode_index_0[i]:_ErrCode_index_0[i+1]]
	case 1002 <= i && i <= 1007:
		i -= 1002
		return _ErrCode_name_1[_ErrCode_index_1[i]:_ErrCode_index_1[i+1]]
	case 2008 <= i && i <= 2027:
		i -= 2008
		return _ErrCode_name_2[_ErrCode_index_2[i]:_ErrCode_index_2[i+1]]
	case 3028 <= i && i <= 3046:
		i -= 3028
		return _ErrCode_name_3[_ErrCode_index_3[i]:_ErrCode_index_3[i+1]]
	default:
		return "ErrCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
