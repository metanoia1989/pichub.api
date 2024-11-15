package utils

//********************************************************
// 语言语法糖 syntactic sugar
//********************************************************

// 三元运算符 ternary operator
func If[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}
