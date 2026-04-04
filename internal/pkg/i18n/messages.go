package i18n

// Init 注册所有错误码的多语言翻译
func Init() {
	// 模板消息（带动态参数，使用 fmt.Sprintf 占位符）
	RegisterTemplate("type_mismatch", map[string]string{
		LangZhCN: "字段 '%s' 类型错误：期望 %s，实际传入 %s",
		LangEnUS: "field '%s' must be %s, got %s",
	})

	// 通用错误码 10001-19999
	Register(10001, map[string]string{
		LangZhCN: "请求参数错误",
		LangEnUS: "invalid request parameters",
	})
	Register(10002, map[string]string{
		LangZhCN: "未授权",
		LangEnUS: "unauthorized",
	})
	Register(10003, map[string]string{
		LangZhCN: "禁止访问",
		LangEnUS: "forbidden",
	})
	Register(10004, map[string]string{
		LangZhCN: "资源不存在",
		LangEnUS: "resource not found",
	})
	Register(10005, map[string]string{
		LangZhCN: "服务器内部错误",
		LangEnUS: "internal server error",
	})
	Register(10006, map[string]string{
		LangZhCN: "请求过于频繁",
		LangEnUS: "too many requests",
	})

	// 用户错误码 20001-20099
	Register(20001, map[string]string{
		LangZhCN: "用户名或密码错误",
		LangEnUS: "invalid username or password",
	})
	Register(20002, map[string]string{
		LangZhCN: "用户已被禁用",
		LangEnUS: "user account is disabled",
	})
	Register(20003, map[string]string{
		LangZhCN: "管理员账号或密码错误",
		LangEnUS: "invalid admin username or password",
	})
	Register(20004, map[string]string{
		LangZhCN: "管理员账号已被禁用",
		LangEnUS: "admin account is disabled",
	})

	// 商品错误码 20101-20199
	Register(20101, map[string]string{
		LangZhCN: "商品不存在",
		LangEnUS: "product not found",
	})
	Register(20102, map[string]string{
		LangZhCN: "商品已下架",
		LangEnUS: "product is off shelf",
	})
	Register(20103, map[string]string{
		LangZhCN: "库存不足",
		LangEnUS: "insufficient stock",
	})

	// 验证器字段多语言 label
	registerFieldLabels()
}

func registerFieldLabels() {
	// LoginReq
	RegisterLabels("LoginReq", map[string]map[string]string{
		"用户名": {LangEnUS: "username"},
		"密码":  {LangEnUS: "password"},
	})

	// AdminLoginReq
	RegisterLabels("AdminLoginReq", map[string]map[string]string{
		"用户名": {LangEnUS: "username"},
		"密码":  {LangEnUS: "password"},
	})

	// CreateReq
	RegisterLabels("CreateReq", map[string]map[string]string{
		"商品名称": {LangEnUS: "product name"},
		"价格":   {LangEnUS: "price"},
		"库存":   {LangEnUS: "stock"},
	})

	// UpdateReq
	RegisterLabels("UpdateReq", map[string]map[string]string{
		"商品名称": {LangEnUS: "product name"},
		"价格":   {LangEnUS: "price"},
		"库存":   {LangEnUS: "stock"},
	})
}
