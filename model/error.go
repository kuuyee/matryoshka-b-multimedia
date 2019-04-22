package model

// Error API错误返回
//
// Error为REST-API返回的错误信息
//
// swagger:model Error
type Error struct {
	// Code 错误码
	//
	// read only: true
	// required: true
	// example: 403
	Code int `json:"code"`
	// Message 错误消息
	//
	// read only: true
	// required: true
	// example: 权限不足
	Message string `json:"message"`
}
