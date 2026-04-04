package main

import (
	"context"
	"os/signal"
	"syscall"

	"goflow/internal/database"
	"goflow/internal/pkg/logger"
)

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
