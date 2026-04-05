package handler

import (
	"context"
	"errors"
	"sky-take-out-go/internal/middleware"
	"sky-take-out-go/internal/pkg/errcode"
	"sky-take-out-go/internal/pkg/i18n"
	"sky-take-out-go/internal/pkg/req"
	"sky-take-out-go/internal/pkg/response"
	"sky-take-out-go/internal/pkg/validator"

	"github.com/gin-gonic/gin"
)

type employeeLoginService interface {
	Login(ctx context.Context, username, password string) (*response.EmployeeLoginResult, error)
	Create(ctx context.Context, r *req.EmployeeCreateReq) error
}

type EmployeeHandler struct {
	employeeSvc employeeLoginService
}

func NewEmployeeHandler(adminSvc employeeLoginService) *EmployeeHandler {
	return &EmployeeHandler{employeeSvc: adminSvc}
}

// Login godoc
// @Summary      员工登录
// @Description  通过用户名和密码登录
// @Tags         员工
// @Accept       json
// @Produce      json
// @Param        body  body  req.EmployeeLoginReq  true  "登录凭证"
// @Success      200   {object}  response.Response  "登录成功"
// @Router       /admin/employee/login [post]
func (h *EmployeeHandler) Login(c *gin.Context) {
	var r req.EmployeeLoginReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.ErrorWithMsg(c, 400, errcode.CodeBadRequest, validator.Translate(err, i18n.GetLang(c)))
		return
	}

	result, err := h.employeeSvc.Login(c.Request.Context(), r.Username, r.Password)
	if err != nil {
		var appErr *errcode.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, errcode.ErrInternal())
		return
	}

	response.Success(c, result)
}

// Create godoc
// @Summary      创建员工
// @Description  创建员工
// @Tags         员工
// @Accept       json
// @Produce      json
// @Param        body  body  req.EmployeeCreateReq  true  "员工信息"
// @Success      200   {object}  response.Response  "创建成功"
// @Router       /admin/employee [post]
func (h *EmployeeHandler) Create(c *gin.Context) {
	// 1.数据绑定
	var r req.EmployeeCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.ErrorWithMsg(c, 400, errcode.CodeBadRequest, validator.Translate(err, i18n.GetLang(c)))
		return
	}

	value, exists := c.Get(string(middleware.ContextKeyClaims))
	if !exists {
		response.Error(c, errcode.ErrUnauthorized())
		return
	}

	// value 做类型断言
	claims, ok := value.(*middleware.CustomClaims)
	if !ok || claims == nil {
		response.Error(c, errcode.ErrUnauthorized())
		return
	}

	ctx := context.WithValue(c.Request.Context(), string(middleware.ContextKeyClaims), claims)

	// 2.调用服务
	if err := h.employeeSvc.Create(ctx, &r); err != nil {
		var appErr *errcode.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, errcode.ErrInternal())
		return
	}

	// 3. 返回响应
	response.Success(c, nil)
}
