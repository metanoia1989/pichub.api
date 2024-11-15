package constants

// 错误码定义
const (
	ErrCodeSuccess      = 0
	ErrCodeParamInvalid = 4001
	ErrCodeUnauthorized = 4002
	ErrCodeForbidden    = 4003
	ErrCodeNotFound     = 4004
	ErrCodeServerError  = 5000
)

// 错误信息映射
var ErrorMessages = map[int]string{
	ErrCodeSuccess:      "success",
	ErrCodeParamInvalid: "invalid parameters",
	ErrCodeUnauthorized: "unauthorized",
	ErrCodeForbidden:    "forbidden",
	ErrCodeNotFound:     "resource not found",
	ErrCodeServerError:  "internal server error",
}

// golang 真讨厌啊，抽象程度太低了，连个枚举都没有
// 想当C，又不如C的底层和灵活，不过是因为现代语言，所以自带包机制，占了优势
