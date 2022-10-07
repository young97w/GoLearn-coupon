package custom_error

const (
	TokenExpired        = "TokenExpired"
	TokenNotValidYet    = "TokenNotValidYet"
	TokenMalformed      = "TokenMalformed"
	TokenInvalid        = "TokenInvalid"
	TokenGenerateFailed = "TokenGenerateFailed"

	AccountNotFind      = "账户不存在"
	AccountExists       = "账户已存在"
	AccountCreateFailed = "账户创建失败"
	InternalError       = "服务器内部错误"
	SaltEmpty           = "盐值为空"
)
