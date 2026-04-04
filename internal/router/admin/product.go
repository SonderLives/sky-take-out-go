package admin

import (
	"goflow/internal/handler"
	"goflow/internal/middleware"
	"goflow/internal/pkg/ratelimit"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(rg *gin.RouterGroup, r *ratelimit.RateLimiter, h *handler.ProductHandler) {
	rg.GET("/products/:id", h.Get)
	rg.POST("/products", middleware.RateLimit(r, 1, 3*time.Second), h.Create)
	rg.GET("/products", middleware.RateLimit(r, 1, 1*time.Second), h.List)
	rg.PUT("/products/:id", h.Update)
	rg.DELETE("/products/:id", h.Delete)
}
