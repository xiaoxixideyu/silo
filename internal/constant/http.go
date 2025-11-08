package constant

var (
	HttpMethodGet    = "GET"
	HttpMethodPost   = "POST"
	HttpMethodPut    = "PUT"
	HttpMethodDelete = "DELETE"
)

const (
	CtxKeyBoSession = "bo_session" // 存储session信息
	CtxKeyUser      = "user"       // 存储用户信息
)

// const .
const (
	BoPrefixV1 = "/api/v1"
	BoPrefixV2 = "/api/v2"
)
