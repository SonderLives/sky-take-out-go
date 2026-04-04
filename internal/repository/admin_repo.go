package repository

import (
	"context"

	"sky-take-out-go/internal/model"

	"gorm.io/gorm"
)

type EmployeeRepository interface {
	Create(ctx context.Context, admin *model.Employee) error
	GetByUsername(ctx context.Context, username string) (*model.Employee, error)
}

type EmployeeRepo struct {
	db *gorm.DB
}

func NewEmployeeRepo(db *gorm.DB) EmployeeRepository {
	return &EmployeeRepo{db: db}
}

// Create 创建管理员
func (r *EmployeeRepo) Create(ctx context.Context, employee *model.Employee) error {
	return r.db.WithContext(ctx).Create(employee).Error
}

// GetByUsername 根据用户名查询管理员
func (r *EmployeeRepo) GetByUsername(ctx context.Context, username string) (*model.Employee, error) {
	var employee model.Employee
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&employee).Error; err != nil {
		return nil, err
	}
	return &employee, nil
}
