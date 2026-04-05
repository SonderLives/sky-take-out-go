package admin

import (
	"sky-take-out-go/internal/handler"
	"sky-take-out-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterEmployeeRoutes(rg *gin.RouterGroup, auth *middleware.AuthMiddleware, h *handler.EmployeeHandler) {
	rg.POST("/login", h.Login)
	// 管理后台全部需认证
	rga := rg.Use(auth.AdminAuth())
	{
		rga.POST("", h.Create)
		rga.GET("/page", h.PageQuery)
		rga.GET("/status/:status", h.SetStatus)
	}
}
