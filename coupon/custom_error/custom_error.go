package custom_error

const (
	TokenExpired        = "TokenExpired"
	TokenNotValidYet    = "TokenNotValidYet"
	TokenMalformed      = "TokenMalformed"
	TokenInvalid        = "TokenInvalid"
	TokenGenerateFailed = "TokenGenerateFailed"

	AddCoffeeFailed    = "添加咖啡失败"
	DeleteCoffeeFailed = "删除咖啡失败"
	UpdateCoffeeFailed = "更新咖啡失败"
	CannotFindCoffee   = "获取咖啡失败"

	AddCouponFailed    = "添加优惠券失败"
	GetCouponFailed    = "获取优惠券失败"
	UpdateCouponFailed = "更新优惠券失败"

	ParameterIncorrect = "参数错误"
)
