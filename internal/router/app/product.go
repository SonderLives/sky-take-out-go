package app

import (
	"sky-take-out-go/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(rg *gin.RouterGroup, h *handler.ProductHandler) {
	rg.GET("/products", h.List)
	rg.GET("/products/:id", h.Get)
}
