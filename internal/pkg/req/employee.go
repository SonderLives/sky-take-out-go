package req

// EmployeeLoginReq 管理员登录请求
type EmployeeLoginReq struct {
	Username string `json:"username" binding:"required" label:"用户名" example:"admin"`
	Password string `json:"password" binding:"required" label:"密码" example:"123456"`
}
