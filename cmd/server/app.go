package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sky-take-out-go/internal/database"
	"sky-take-out-go/internal/middleware"
	"sky-take-out-go/internal/pkg/i18n"
	"sky-take-out-go/internal/pkg/validator"
	"sky-take-out-go/internal/router"
	"sky-take-out-go/internal/svc"
	"sky-take-out-go/migration"
	"time"

	"sky-take-out-go/internal/config"
	"sky-take-out-go/internal/mq"
	"sky-take-out-go/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type App struct {
	cfg        *config.Config
	httpServer *http.Server
	mqRouter   *mq.Router
	redis      *redis.Client
}

// InitAll 初始化所有组件
func InitAll() (*config.Config, *gorm.DB, *redis.Client, *mq.Publisher,
	*svc.ServiceContext, *gin.Engine, *mq.Router, error) {
	// 1. 加载配置
	cfg, err := config.Load()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	// 2. 初始化日志
	logger.Init(&cfg.Log)
	defer logger.Sync()

	// 3. 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 4. 初始化数据库和 Redis
	db, rdb, err := initDatabases(cfg)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	// 5. 数据库迁移
	if cfg.Server.AutoMigrate {
		if err := migration.AutoMigrate(db); err != nil {
			return nil, nil, nil, nil, nil, nil, nil, err
		}
	}

	// 6. 初始化 JWT 和认证中间件
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWT.Secret)

	// 7. 初始化验证器和翻译
	validator.Init()
	i18n.Init()

	// 8. 初始化 MQ Publisher
	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}
	mqPublisher, err := mq.NewPublisher(&cfg.MQ, rdb, sqlDB)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	// 9. 初始化服务上下文
	svcCtx := svc.NewServiceContext(cfg, db, rdb, mqPublisher)

	// 10. 初始化路由
	r := router.Setup(svcCtx, authMiddleware)

	// 11. 初始化 MQ Router
	mqRouter, err := mq.NewRouter(&cfg.MQ, rdb, sqlDB)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	// 12. 返回所有初始化的对象
	return cfg, db, rdb, mqPublisher, svcCtx, r, mqRouter, nil
}

// 初始化数据库和 Redis
func initDatabases(cfg *config.Config) (*gorm.DB, *redis.Client, error) {
	// 初始化 MySQL
	db, err := database.InitMySQL(&cfg.MySQL, &cfg.Log)
	if err != nil {
		return nil, nil, err
	}

	// 初始化 Redis
	rdb, err := database.InitRedis(&cfg.Redis, &cfg.Log)
	if err != nil {
		return nil, nil, err
	}

	return db, rdb, nil
}

func NewApp(cfg *config.Config, handler http.Handler, mqRouter *mq.Router, redisClient *redis.Client) *App {
	return &App{
		cfg: cfg,
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
			Handler:      handler,
			ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		},
		mqRouter: mqRouter,
		redis:    redisClient,
	}
}

func (a *App) Run(ctx context.Context) error {
	trimCtx, trimCancel := context.WithCancel(ctx)
	mq.StartTrimmer(trimCtx, &a.cfg.MQ, a.redis)

	errCh := make(chan error, 2)

	go func() {
		logger.Log.Infof("server starting on port %d", a.cfg.Server.Port)
		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("http server failed: %w", err)
		}
	}()

	if a.mqRouter != nil {
		go func() {
			if err := a.mqRouter.Run(ctx); err != nil && ctx.Err() == nil {
				errCh <- fmt.Errorf("mq router failed: %w", err)
			}
		}()
	}

	var runErr error
	select {
	case <-ctx.Done():
	case err := <-errCh:
		runErr = err
	}

	trimCancel()

	shutdownTimeout := 10 * time.Second
	if a.cfg.Server.ShutdownTimeout > 0 {
		shutdownTimeout = time.Duration(a.cfg.Server.ShutdownTimeout) * time.Second
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if a.mqRouter != nil {
		if err := a.mqRouter.Close(); err != nil && runErr == nil {
			runErr = fmt.Errorf("close mq router failed: %w", err)
		}
	}

	if err := a.httpServer.Shutdown(shutdownCtx); err != nil && runErr == nil {
		runErr = fmt.Errorf("shutdown http server failed: %w", err)
	}

	return runErr
}
