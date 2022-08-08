package response

// 响应结构体
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// 响应码
const (
	StatusOk    = 0
	StatusError = 100
)
