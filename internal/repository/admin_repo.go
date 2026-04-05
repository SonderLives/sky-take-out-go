package repository

import (
	"context"
	"errors"
	"strings"

	"sky-take-out-go/internal/model"
	"sky-take-out-go/internal/pkg/errcode"

	"github.com/go-sql-driver/mysql"
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

// Create 创建员工
func (r *EmployeeRepo) Create(ctx context.Context, employee *model.Employee) error {
	if err := r.db.WithContext(ctx).Create(employee).Error; err != nil {
		var me *mysql.MySQLError
		if errors.As(err, &me) {
			switch me.Number {
			case 1062: // MySQL 1062: Duplicate entry
				// 仅对用户名唯一索引返回明确业务错误，其余唯一冲突可按需扩展
				if strings.Contains(me.Message, "idx_username") {
					return errcode.ErrUserNameAlreadyExists()
				}
			}
		}
		return err
	}
	return nil
}

// GetByUsername 根据用户名查询员工
func (r *EmployeeRepo) GetByUsername(ctx context.Context, username string) (*model.Employee, error) {
	var employee model.Employee
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&employee).Error; err != nil {
		return nil, err
	}
	return &employee, nil
}
