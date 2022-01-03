//go:generate stringer -type ErrCode -linecomment
package errors

type ErrCode int

const (
	// 其他错误
	ErrUnKnown  ErrCode = iota // 未知错误
	ErrReadFile                // 代码文件读取失败
	// 基础扫描错误
	scanErr               ErrCode = 1000 + iota
	ErrScanUncloseChar            // 字符缺少关闭的 ' 符号
	ErrScanUncloseString          // 字符串缺少关闭的 " 符号
	ErrScanUncloseComment         // 多行注释缺少对应的关闭 */ 符号
	ErrScanHexFormat              // 符号 %c 不是一个16进制编码字符
	ErrScanUnicodeFormat          // 符号 %c 不是一个Unicode编码字符
	//预处理错误
	macroErr                     ErrCode = 2000 + iota
	ErrMacroHashHashPos                  // ## 不能出现在宏表达式的起始或结束位置
	ErrMacroHashHashExpr                 // ## 不能用来连接 %s 和 %s
	ErrMacroHashExpr                     // # 符号后面必须跟着一个宏参数
	ErrMacroCallParamCount               // 宏调用参数数量错误，支持%d个参数，使用了%d个参数
	ErrMacroUnexpectedElseIf             // 不应该出现的 #elif 宏
	ErrMacroUnexpectedElse               // 不应该出现的 #else 宏
	ErrMacroUnexpectedEndIf              // 不应该出现的 #endif 宏
	ErrMacroExpectedIdent                // 这里应该是一个名称，不应该出现 %s 符号
	ErrMacroExpectedGot                  // 这里应该是一个 %s ，不应该出现 %s
	ErrMacroExpectedPunctuator           // 这里应该是一个 %s 符号，不应该出现 %s 符号
	ErrMacroEnd                          // 这里应该是宏结尾了，不应该出现 %s 符号
	ErrMacroExpectedTokenGotEof          // 需要符号为 %s，意外的遇到了文件尾
	ErrMacroConstExpr                    // 错误的宏常量表达式 %s
	ErrMacroDuplicateIdent               // 重复定义了符号 %s
	ErrMacroInvalidIncludeString         // #include 包含错误的字符串 %s
	ErrMacroInvalidIncludeMacro          // 错误的 #include 宏
	ErrMacroIncludeFileRead              // #include的文件 %s 读取错误 %s
	ErrMacroIncludeFileNoFound           // #include的文件不存在 %s
	ErrMacroExprUnexpectedToken          // 非预期的宏表达式符号%s
	// 语法错误
	syntaxError                       ErrCode = 3000 + iota
	ErrSyntaxExpectedGot                      // 这里应该是一个 %s ，不应该出现 %s
	ErrSyntaxExpectedIdentGot                 // 这里应该是一个名称，不应该出现 %s 符号
	ErrSyntaxUnexpectedTypeSpecifier          // 非预期的类型定义符号 %s
	ErrSyntaxDuplicateTypeSpecifier           // 重复的类型定义符号 %s
	ErrSyntaxDuplicateTypeQualifier           // 重复的类型修饰符号 %s
	ErrSyntaxExpectedRecordMemberName         // 类型定义符号之后应该是成员变量的名称
	ErrSyntaxRedefineFunc                     // 重复声明函数 %s，上次声明的位置 %s
	ErrSyntaxRedefineVar                      // 重复声明的变量名 %s，上次声明的位置 %s
	ErrSyntaxRedefineIdent                    // 重复的标识符 %s，上次声明的位置 %s
	ErrSyntaxRedefinedType                    // 重复定义的类型 %s，上次定义的位置 %s
	ErrSyntaxRedefinedStruct                  // 重复定义的结构体 %s，上次定义的位置 %s
	ErrSyntaxRedefinedUnion                   // 重复定义的联合体 %s，上次定义的位置 %s
	ErrSyntaxRedefinedEnum                    // 重复定义的枚举 %s，上次定义的位置 %s
	ErrSyntaxRedefinedLabel                   // 重复定义的标签 %s，上次定义的位置 %s
	ErrSyntaxUndefinedIdent                   // 未定义的标识符 %s
	ErrSyntaxUndefinedLabel                   // 未定义的标签 %s
	ErrSyntaxIncompleteStruct                 // 不完全的结构体类型 %s
	ErrSyntaxIncompleteUnion                  // 不完全的联合体类型 %s
	typeError                         ErrCode = 4000 + iota
	ErrTypeImmediateMakeAddress               // 无法对临时变量进行取地址操作
)
