package main

import (
	"context"
	"os/signal"
	"syscall"

	_ "sky-take-out-go/cmd/server/docs"
	"sky-take-out-go/internal/database"
	"sky-take-out-go/internal/pkg/logger"
)

// @title           苍穹外卖
// @version         1.0
// @description     苍穹外卖后端 API 文档
// @BasePath        /
// @schemes         http
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Token
func main() {
	// 1. InitAll 函数初始化所有资源
	cfg, db, rdb, mqPublisher, _, r, mqRouter, err := InitAll()
	if err != nil {
		logger.Log.Fatalf("failed to initialize all components: %v", err)
	}

	defer mqPublisher.Close()
	defer database.CloseMySQL(db)
	defer database.CloseRedis(rdb)

	app := NewApp(cfg, r, mqRouter, rdb)

	runCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.Run(runCtx); err != nil {
		logger.Log.Fatalf("app run failed: %v", err)
	}

	logger.Log.Info("server exited")
}
