package req

// UpdateReq 更新商品请求
type UpdateReq struct {
	Name        *string  `json:"name" binding:"omitempty,max=255" label:"商品名称" example:"宫保鸡丁"`
	Description *string  `json:"description" example:"经典川菜，微辣"`
	Price       *float64 `json:"price" binding:"omitempty,gte=0" label:"价格" example:"32.5"`
	Stock       *int     `json:"stock" binding:"omitempty,gte=0" label:"库存" example:"100"`
	Status      *int8    `json:"status" binding:"omitempty,oneof=0 1" example:"1"`
	CategoryID  *uint    `json:"category_id" example:"2"`
}

// CreateReq 创建商品请求
type CreateReq struct {
	Name        string  `json:"name" binding:"required,max=255" label:"商品名称" example:"宫保鸡丁"`
	Description string  `json:"description" example:"经典川菜，微辣"`
	Price       float64 `json:"price" binding:"required,gt=0" label:"价格" example:"32.5"`
	Stock       int     `json:"stock" binding:"gte=0" label:"库存" example:"100"`
	CategoryID  uint    `json:"category_id" example:"2"`
}
