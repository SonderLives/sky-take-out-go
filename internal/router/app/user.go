package app

import (
	"goflow/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup, h *handler.UserHandler) {
	rg.POST("/login", h.Login)
}
