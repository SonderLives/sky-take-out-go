package req

// EmployeeLoginReq 管理员登录请求
type EmployeeLoginReq struct {
	Username string `json:"username" binding:"required" label:"用户名" example:"admin"`
	Password string `json:"password" binding:"required" label:"密码" example:"123456"`
}

// EmployeeCreateReq 管理员创建请求
type EmployeeCreateReq struct {
	// ID 字段为自增且后端生成
	ID int `json:"id" label:"ID" example:"1"`

	// 身份证号：必填 + 18位 + 基础正则（支持末尾X/x）
	// 正则支持 0x7C -> | 正则支持 0x2C -> ,
	IDNumber string `json:"idNumber" binding:"required,len=18,regexp=^[1-9]\\d{5}(180x7C190x7C20)\\d{2}(0[1-9]0x7C1[0-2])(0[1-9]0x7C[12]\\d0x7C3[01])\\d{3}[\\dXx]$" label:"身份证号" example:"110105199003078888"`

	Name string `json:"name" binding:"required" label:"姓名" example:"张三"`

	// 手机号：必填 + 11位 + 大陆手机号正则
	Phone string `json:"phone" binding:"required,len=11,regexp=^1[3-9]\\d{9}$" label:"手机号" example:"13800000000"`

	// 性别：限制枚举值
	Sex string `json:"sex" binding:"required" label:"性别" example:"男"`
	// 用户名：必填，唯一
	Username string `json:"username" binding:"required" label:"用户名" example:"admin"`
}
