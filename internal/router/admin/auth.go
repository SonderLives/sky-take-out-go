package admin

import (
	"goflow/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, h *handler.AdminHandler) {
	rg.POST("/login", h.Login)
}
