package response

// EmployeeLoginResult 登录返回结果
type EmployeeLoginResult struct {
	ID       int64  `json:"id" example:"1"`
	Name     string `json:"name" example:"管理员"`
	Token    string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	UserName string `json:"userName" example:"admin"`
}

// UserInfo 用户信息（用于登录响应，替代 gin.H 使 service 层与框架解耦）
type UserInfo struct {
	ID       int64  `json:"id" example:"101"`
	Username string `json:"username" example:"alice"`
	Nickname string `json:"nickname" example:"Alice"`
	Avatar   string `json:"avatar" example:"https://example.com/avatar/alice.png"`
}

// EmployeePageItemDTO 员工分页列表项
type EmployeePageItemDTO struct {
	ID         int64  `json:"id" example:"1"`
	Username   string `json:"username" example:"admin"`
	Name       string `json:"name" example:"管理员"`
	Phone      string `json:"phone" example:"13800138000"`
	Sex        string `json:"sex" example:"男"`
	IDNumber   string `json:"idNumber" example:"110105********8888"`
	Status     int    `json:"status" example:"1"`
	CreateTime string `json:"createTime" example:"2026-04-06 21:30:00"`
	UpdateTime string `json:"updateTime" example:"2026-04-06 21:30:00"`
	CreateUser int64  `json:"createUser" example:"1"`
	UpdateUser int64  `json:"updateUser" example:"1"`
}
