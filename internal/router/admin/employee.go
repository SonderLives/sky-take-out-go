package admin

import (
	"sky-take-out-go/internal/handler"
	"sky-take-out-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterEmployeeRoutes(rg *gin.RouterGroup, auth *middleware.AuthMiddleware, h *handler.EmployeeHandler) {
	rg.POST("/login", h.Login)
	// 管理后台全部需认证
	rg.Use(auth.AdminAuth()).POST("", h.Create)
	//rg.POST("", h.Create)
}
