package constants

// 用户状态码
const (
	StatusActive   uint8 = 1
	StatusInactive uint8 = 0
)

// github 相关默认值
const (
	DefaultRepoBranch = "master"
	DefaultPageSize   = 10
	MaxPageSize       = 100
	RepositoryHost    = "https://github.com"
)
