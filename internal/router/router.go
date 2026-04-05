package router

import (
	"context"
	"net/http"
	"sky-take-out-go/internal/handler"
	"sky-take-out-go/internal/middleware"
	"sky-take-out-go/internal/pkg/response"
	"sky-take-out-go/internal/router/admin"
	"sky-take-out-go/internal/router/app"
	"sky-take-out-go/internal/svc"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Setup 接收 ServiceContext，内部创建所有 handler 并注册路由。
// 新增模块只需在此处添加 handler 和路由，无需修改函数签名。
func Setup(svcCtx *svc.ServiceContext, auth *middleware.AuthMiddleware) *gin.Engine {
	r := gin.New()

	// 全局中间件
	r.Use(middleware.RequestID())
	r.Use(middleware.I18n())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS(svcCtx.Config.Server.CORSOrigins))

	if svcCtx.Config.Server.Mode == "debug" {
		// 1.swagger 文档
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		// 2.健康检查：检测 MySQL 和 Redis 的真实连通性
		r.GET("/health", func(c *gin.Context) {
			ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
			defer cancel()

			if err := svcCtx.HealthCheck(ctx); err != nil {
				c.JSON(http.StatusServiceUnavailable, response.Response{
					Code:    -1,
					Message: "unhealthy: " + err.Error(),
				})
				return
			}
			response.Success(c, gin.H{"status": "ok"})
		})
	}

	// 创建 handler（从 ServiceContext 获取依赖）
	productHandler := handler.NewProductHandler(svcCtx.ProductSvc, svcCtx.MQPublisher)
	userHandler := handler.NewUserHandler(svcCtx.UserSvc)
	employeeHandler := handler.NewEmployeeHandler(svcCtx.EmployeeSvc)

	// App 客户端接口
	appGroup := r.Group("/app/v1")
	app.RegisterUserRoutes(appGroup, userHandler)

	// App 商品接口
	appProductGroup := appGroup.Group("")
	app.RegisterProductRoutes(appProductGroup, productHandler)

	// App 需要认证的接口
	appAuth := appGroup.Group("")
	appAuth.Use(auth.AppAuth())
	// 后续扩展 App 需认证的接口

	// 员工管理后台接口
	employeeGroup := r.Group("/admin/employee")
	admin.RegisterEmployeeRoutes(employeeGroup, auth, employeeHandler)

	//admin.RegisterProductRoutes(adminAuth, svcCtx.RateLimiter, productHandler)

	return r
}
