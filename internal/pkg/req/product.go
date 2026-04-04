package req

// UpdateReq 更新商品请求
type UpdateReq struct {
	Name        *string  `json:"name" binding:"omitempty,max=255" label:"商品名称"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price" binding:"omitempty,gte=0" label:"价格"`
	Stock       *int     `json:"stock" binding:"omitempty,gte=0" label:"库存"`
	Status      *int8    `json:"status" binding:"omitempty,oneof=0 1"`
	CategoryID  *uint    `json:"category_id"`
}

// CreateReq 创建商品请求
type CreateReq struct {
	Name        string  `json:"name" binding:"required,max=255" label:"商品名称"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0" label:"价格"`
	Stock       int     `json:"stock" binding:"gte=0" label:"库存"`
	CategoryID  uint    `json:"category_id"`
}
