package admin

import (
	"sky-take-out-go/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, h *handler.EmployeeHandler) {
	rg.POST("/employee/login", h.Login)
}
