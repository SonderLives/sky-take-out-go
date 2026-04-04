package handler

import (
	"context"
	"errors"
	"sky-take-out-go/internal/pkg/errcode"
	"sky-take-out-go/internal/pkg/i18n"
	"sky-take-out-go/internal/pkg/req"
	"sky-take-out-go/internal/pkg/response"
	"sky-take-out-go/internal/pkg/validator"

	"github.com/gin-gonic/gin"
)

type employeeLoginService interface {
	Login(ctx context.Context, username, password string) (*response.EmployeeLoginResult, error)
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
