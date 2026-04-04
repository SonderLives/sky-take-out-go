package migration

import (
	"fmt"

	"sky-take-out-go/internal/model"
	"sky-take-out-go/internal/pkg/logger"

	"gorm.io/gorm"
)

// AutoMigrate 自动迁移数据库表结构，返回 error 让调用方统一处理
func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&model.User{},
		&model.Employee{},
		&model.Product{},
	)
	if err != nil {
		return fmt.Errorf("auto migrate failed: %w", err)
	}
	logger.Log.Info("database migration completed")
	return nil
}
