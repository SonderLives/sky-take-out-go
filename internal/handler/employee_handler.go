package handler

import (
	"context"
	"errors"
	"sky-take-out-go/internal/middleware"
	"sky-take-out-go/internal/pkg/errcode"
	"sky-take-out-go/internal/pkg/i18n"
	"sky-take-out-go/internal/pkg/logger"
	"sky-take-out-go/internal/pkg/req"
	"sky-take-out-go/internal/pkg/response"
	"sky-take-out-go/internal/pkg/validator"

	"github.com/gin-gonic/gin"
)

type employeeLoginService interface {
	Login(ctx context.Context, r *req.EmployeeLoginReq) (*response.EmployeeLoginResult, error)
	Create(ctx context.Context, r *req.EmployeeCreateReq) error
	PageQuery(ctx context.Context, r *req.EmployeePageReq) (*response.PageResult, error)
	SetStatus(ctx context.Context, status string, id string) error
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
// @Success      200   {object}  response.Response{data=response.EmployeeLoginResult}  "登录成功"
// @Failure      400   {object}  response.Response  "请求参数错误"
// @Failure      500   {object}  response.Response  "服务器内部错误"
// @Router       /admin/employee/login [post]
func (h *EmployeeHandler) Login(c *gin.Context) {
	var r req.EmployeeLoginReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.ErrorWithMsg(c, 400, errcode.CodeBadRequest, validator.Translate(err, i18n.GetLang(c)))
		return
	}

	result, err := h.employeeSvc.Login(c.Request.Context(), &r)
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
// @Description  创建员工（需管理员权限）
// @Tags         员工
// @Accept       json
// @Produce      json
// @Param        body  body  req.EmployeeCreateReq  true  "员工信息"
// @Success      200   {object}  response.Response  "创建成功"
// @Failure      400   {object}  response.Response  "请求参数错误"
// @Failure      401   {object}  response.Response  "未授权"
// @Failure      500   {object}  response.Response  "服务器内部错误"
// @Security     BearerAuth
// @Router       /admin/employee [post]
func (h *EmployeeHandler) Create(c *gin.Context) {
	// 1.数据绑定
	var r req.EmployeeCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.ErrorWithMsg(c, 400, errcode.CodeBadRequest, validator.Translate(err, i18n.GetLang(c)))
		return
	}

	claims, ok := c.Request.Context().Value(string(middleware.ContextKeyClaims)).(*middleware.CustomClaims)
	if !ok || claims == nil {
		response.Error(c, errcode.ErrUnauthorized())
		return
	}

	// 2.调用服务
	if err := h.employeeSvc.Create(c.Request.Context(), &r); err != nil {
		var appErr *errcode.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, errcode.ErrInternal())
		return
	}

	// 3. 返回响应
	logger.WithCtx(c.Request.Context()).Infow("employee created",
		"username", r.Username,
		"operator_id", claims.UserID,
	)
	response.Success(c, nil)
}

// PageQuery godoc
// @Summary      分页查询员工
// @Description  分页查询员工（需管理员权限）
// @Tags         员工
// @Accept       json
// @Produce      json
// @Param        name  query  string  false  "员工姓名"  default("张三")
// @Param        page  query  int  true  "页码"  default(1)
// @Param        pageSize  query  int  true  "每页数量"  default(10)
// @Success      200   {object}  response.Response{data=response.PageResult{records=[]response.EmployeePageItemDTO}}  "查询成功"
// @Failure      400   {object}  response.Response  "请求参数错误"
// @Failure      401   {object}  response.Response  "未授权"
// @Failure      500   {object}  response.Response  "服务器内部错误"
// @Security     BearerAuth
// @Router       /admin/employee/page [get]
func (h *EmployeeHandler) PageQuery(c *gin.Context) {
	// 1.数据绑定
	var r req.EmployeePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.ErrorWithMsg(c, 400, errcode.CodeBadRequest, validator.Translate(err, i18n.GetLang(c)))
		return
	}

	claims, ok := c.Request.Context().Value(string(middleware.ContextKeyClaims)).(*middleware.CustomClaims)
	if !ok || claims == nil {
		response.Error(c, errcode.ErrUnauthorized())
		return
	}

	// 2.调用服务
	result, err := h.employeeSvc.PageQuery(c.Request.Context(), &r)
	if err != nil {
		var appErr *errcode.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, errcode.ErrInternal())
		return
	}

	logger.WithCtx(c.Request.Context()).Infow("employee page query",
		"name", r.Name,
		"page", r.Page,
		"page_size", r.PageSize,
		"operator_id", claims.UserID,
	)
	response.Success(c, result)
}

func (h *EmployeeHandler) SetStatus(c *gin.Context) {
	// 1.数据绑定
	status := c.Param("status")
	id := c.Query("id")

	if err := h.employeeSvc.SetStatus(c.Request.Context(), status, id); err != nil {
		var appErr *errcode.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, errcode.ErrInternal())
		return
	}
	response.Success(c, nil)
}
