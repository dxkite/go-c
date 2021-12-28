//go:generate stringer -type ErrCode -linecomment
package errors

type ErrCode int

const (
	ErrUnKnown  ErrCode = iota // 未知错误
	ErrReadFile                // 代码文件读取失败
	// 基础扫描错误
	scanErrBegin      ErrCode = 1000 + iota
	ErrUncloseChar            // 字符缺少关闭的 ' 符号
	ErrUncloseString          // 字符串缺少关闭的 " 符号
	ErrUncloseComment         // 多行注释缺少对应的关闭 */ 符号
	ErrHexFormat              // 符号 %c 不是一个16进制编码字符
	ErrUnicodeFormat          // 符号 %c 不是一个Unicode编码字符
	//预处理错误
	preprocessErrBegin          ErrCode = 2000 + iota
	ErrMacroHashHashPos                 // ## 不能出现在宏表达式的起始或结束位置
	ErrMacroHashHashExpr                // ## 不能用来连接 %s 和 %s
	ErrMacroHashExpr                    // # 符号后面必须跟着一个宏参数
	ErrMacroCallParamCount              // 宏调用参数数量错误，支持%d个参数，使用了%d个参数
	ErrUnexpectedElseIfMacro            // 不应该出现的 #elif 宏
	ErrUnexpectedElseMacro              // 不应该出现的 #else 宏
	ErrUnexpectedEndIfMacro             // 不应该出现的 #endif 宏
	ErrExpectedMacroIdent               // 这里应该是一个名称，不应该出现 %s 符号
	ErrExpectedMacroGot                 // 这里应该是一个 %s ，不应该出现 %s
	ErrExpectedMacroPunctuator          // 这里应该是一个 %s 符号，不应该出现 %s 符号
	ErrMacroEnd                         // 这里应该是宏结尾了，不应该出现 %s 符号
	ErrMacroExpectedTokenGotEof         // 需要符号为 %s，意外的遇到了文件尾
	ErrMacroConstExpr                   // 错误的宏常量表达式 %s
	ErrDuplicateDefine                  // 重复定义了符号 %s
	ErrInvalidIncludeString             // #include 包含错误的字符串 %s
	ErrInvalidIncludeMacro              // 错误的 #include 宏
	ErrIncludeFileRead                  // #include的文件 %s 读取错误 %s
	ErrIncludeFileNoFound               // #include的文件不存在 %s
	ErrMacroExprUnexpectedToken         // 非预期的宏表达式符号%s

	// 语法错误
	parserError                       ErrCode = 3000 + iota
	ErrSyntaxExpectedGot                      // 这里应该是一个 %s ，不应该出现 %s
	ErrSyntaxExpectedIdentGot                 // 这里应该是一个名称，不应该出现 %s 符号
	ErrSyntaxUnexpectedTypeSpecifier          // 非预期的类型定义符号 %s
	ErrSyntaxDuplicateTypeSpecifier           // 重复的类型定义符号 %s
	ErrSyntaxDuplicateTypeQualifier           // 重复的类型修饰符号 %s
	ErrSyntaxExpectedRecordMemberName         // 类型定义符号之后应该是成员变量的名称
)
