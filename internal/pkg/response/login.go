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

// AdminInfo 管理员信息（用于登录响应）
type AdminInfo struct {
	ID       int64  `json:"id" example:"1"`
	Username string `json:"username" example:"admin"`
	Nickname string `json:"nickname" example:"系统管理员"`
	Avatar   string `json:"avatar" example:"https://example.com/avatar/admin.png"`
	Role     string `json:"role" example:"admin"`
}
