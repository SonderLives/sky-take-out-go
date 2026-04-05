package repository

import (
	"context"
	"errors"
	"sky-take-out-go/internal/pkg/req"
	"strings"

	"sky-take-out-go/internal/model"
	"sky-take-out-go/internal/pkg/errcode"
	"sky-take-out-go/internal/pkg/logger"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

const (
	defaultPage     = 1
	defaultPageSize = 10
	maxPageSize     = 100
)

type EmployeeRepository interface {
	Create(ctx context.Context, admin *model.Employee) error
	GetByUsername(ctx context.Context, username string) (*model.Employee, error)
	PageQuery(ctx context.Context, r *req.EmployeePageReq) ([]model.Employee, int64, error)
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
					logger.WithCtx(ctx).Warnw("employee create duplicate username",
						"username", employee.Username,
						"mysql_code", me.Number,
					)
					return errcode.ErrUserNameAlreadyExists()
				}
			}
		}
		logger.WithCtx(ctx).Errorw("employee create db error",
			"username", employee.Username,
			"error", err,
		)
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

// PageQuery 分页查询员工
func (r *EmployeeRepo) PageQuery(ctx context.Context, pageReq *req.EmployeePageReq) ([]model.Employee, int64, error) {
	var list []model.Employee
	var total int64

	page, pageSize := sanitizePagination(pageReq.Page, pageReq.PageSize)

	// 1. 构建基础查询（从 db 开始，避免污染全局）
	query := r.db.WithContext(ctx).Model(&model.Employee{})

	// 2. 动态添加条件：姓名模糊查询
	if pageReq.Name != "" {
		query = query.Where("name LIKE ?", "%"+pageReq.Name+"%")
	}

	// 3. 查询总数（用于分页）
	// 注意：Count 会复用前面的 Where 条件
	if err := query.Count(&total).Error; err != nil {
		var me *mysql.MySQLError
		if errors.As(err, &me) {
			logger.WithCtx(ctx).Warnw("employee page count mysql error",
				"name", pageReq.Name,
				"page", page,
				"page_size", pageSize,
				"mysql_code", me.Number,
				"mysql_message", me.Message,
			)
		}
		return nil, 0, err
	}

	// total 为 0 时提前返回，避免一次无意义查询
	if total == 0 {
		return nil, 0, nil
	}

	// 4. 添加分页 + 排序 + 执行查询
	if err := query.
		Offset((page - 1) * pageSize). // 页码从 1 开始
		Limit(pageSize).
		Order("create_time DESC").
		Find(&list).
		Error; err != nil {
		var me *mysql.MySQLError
		if errors.As(err, &me) {
			logger.WithCtx(ctx).Warnw("employee page query mysql error",
				"name", pageReq.Name,
				"page", page,
				"page_size", pageSize,
				"mysql_code", me.Number,
				"mysql_message", me.Message,
			)
		}
		return nil, 0, err
	}

	return list, total, nil
}

func sanitizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = defaultPage
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	return page, pageSize
}
