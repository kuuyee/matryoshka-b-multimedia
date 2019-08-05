package model

// Meta 文件meta信息
//
// Meta为文件的基本信息（identifier，大小，hash等）
//
// swagger:model Meta
type Meta struct {
	// Ident 文件identifier
	//
	// read only: true
	// required: true
	// example: abcde.jpg
	Ident string `json:"ident"`
	// Length 文件字节数
	//
	// read only: true
	// required: true
	// example: 1024
	Length int64 `json:"length"`
	// Type 文件类型
	//
	// image: 图片
	//
	// read only: true
	// required: true
	// example: image
	Type string `json:"type"`
}
