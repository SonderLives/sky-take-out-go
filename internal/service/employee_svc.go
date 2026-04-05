package service

import (
	"context"
	"errors"
	"sky-take-out-go/internal/pkg/req"
	"strings"
	"time"

	"sky-take-out-go/internal/config"
	"sky-take-out-go/internal/middleware"
	"sky-take-out-go/internal/model"
	"sky-take-out-go/internal/pkg/errcode"
	"sky-take-out-go/internal/pkg/logger"
	"sky-take-out-go/internal/pkg/response"
	"sky-take-out-go/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const defaultEmployeePassword = "123456"
const employeeTimeLayout = "2006-01-02 15:04:05"

type EmployeeService interface {
	Create(ctx context.Context, r *req.EmployeeCreateReq) error
	Login(ctx context.Context, r *req.EmployeeLoginReq) (*response.EmployeeLoginResult, error)
	PageQuery(ctx context.Context, r *req.EmployeePageReq) (*response.PageResult, error)
}

type employeeService struct {
	repo      repository.EmployeeRepository
	jwtSecret []byte
	jwtExpire time.Duration
}

func (s *employeeService) SetStatus(ctx context.Context, status string, id string) error {
	//TODO implement me
	panic("implement me")
}

func NewEmployeeService(repo repository.EmployeeRepository, cfg *config.Config) EmployeeService {
	jwtExpire := cfg.JWT.Expire
	if jwtExpire <= 0 {
		jwtExpire = 7200
	}
	return &employeeService{
		repo:      repo,
		jwtSecret: []byte(cfg.JWT.Secret),
		jwtExpire: time.Duration(jwtExpire) * time.Second,
	}
}

// Login 员工登录
func (s *employeeService) Login(ctx context.Context, r *req.EmployeeLoginReq) (*response.EmployeeLoginResult, error) {
	employee, err := s.repo.GetByUsername(ctx, r.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.WithCtx(ctx).Warnw("employee login user not found", "username", r.Username)
			return nil, errcode.ErrUserNotFound()
		}
		return nil, errcode.ErrInternal()
	}
	if err := bcrypt.CompareHashAndPassword([]byte(employee.Password), []byte(r.Password)); err != nil {
		logger.WithCtx(ctx).Warnw("employee login password mismatch", "username", r.Username)
		return nil, errcode.ErrAdminOrPassword()
	}

	if employee.Status == model.EmployeeStatusDisabled {
		logger.WithCtx(ctx).Warnw("employee login account disabled",
			"username", r.Username,
			"employee_id", employee.ID,
		)
		return nil, errcode.ErrAdminDisabled()
	}

	expireAt := time.Now().Add(s.jwtExpire)
	claims := middleware.CustomClaims{
		UserID:   employee.ID,
		Username: employee.Username,
		Role:     middleware.RoleAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(s.jwtSecret)
	if err != nil {
		logger.WithCtx(ctx).Errorw("employee login sign token failed", "employee_id", employee.ID, "error", err)
		return nil, errcode.ErrInternal()
	}

	logger.WithCtx(ctx).Infow("employee login success",
		"employee_id", employee.ID,
		"username", employee.Username,
	)

	return &response.EmployeeLoginResult{
		Token:    tokenStr,
		ID:       employee.ID,
		Name:     employee.Name,
		UserName: employee.Username,
	}, nil
}

// Create 创建员工
func (s *employeeService) Create(ctx context.Context, r *req.EmployeeCreateReq) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(defaultEmployeePassword), bcrypt.DefaultCost)
	if err != nil {
		logger.WithCtx(ctx).Errorw("employee create hash password failed", "error", err)
		return errcode.ErrInternal()
	}

	// 安全断言，避免 claims 为 nil 时 panic
	claims, ok := ctx.Value(string(middleware.ContextKeyClaims)).(*middleware.CustomClaims)
	if !ok || claims == nil {
		logger.WithCtx(ctx).Warnw("employee create missing claims")
		return errcode.ErrInternal()
	}

	now := time.Now()
	employee := &model.Employee{
		Username:   r.Username,
		Name:       r.Name,
		Password:   string(hashed),
		Phone:      r.Phone,
		Sex:        r.Sex,
		IDNumber:   r.IDNumber,
		Status:     model.EmployeeStatusEnabled,
		CreateTime: now,
		UpdateTime: now,
		CreateUser: claims.UserID,
		UpdateUser: claims.UserID,
	}

	if err := s.repo.Create(ctx, employee); err != nil {
		var appErr *errcode.AppError
		if errors.As(err, &appErr) {
			logger.WithCtx(ctx).Warnw("employee create failed",
				"username", r.Username,
				"operator_id", claims.UserID,
				"code", appErr.Code,
				"error", appErr.Message,
			)
			return appErr
		}
		logger.WithCtx(ctx).Errorw("employee create failed with unknown error",
			"username", r.Username,
			"operator_id", claims.UserID,
			"error", err,
		)
		return errcode.ErrInternal()
	}
	return nil
}

// PageQuery 分页查询员工
func (s *employeeService) PageQuery(ctx context.Context, r *req.EmployeePageReq) (*response.PageResult, error) {
	records, total, err := s.repo.PageQuery(ctx, r)
	if err != nil {
		var appErr *errcode.AppError
		if errors.As(err, &appErr) {
			return nil, appErr
		}
		logger.WithCtx(ctx).Errorw("employee page query failed", "error", err)
		return nil, errcode.ErrInternal()
	}

	items := make([]response.EmployeePageItemDTO, 0, len(records))
	for _, e := range records {
		items = append(items, response.EmployeePageItemDTO{
			ID:         e.ID,
			Username:   e.Username,
			Name:       e.Name,
			Phone:      e.Phone,
			Sex:        e.Sex,
			IDNumber:   maskIDNumber(e.IDNumber),
			Status:     int(e.Status),
			CreateTime: formatTime(e.CreateTime),
			UpdateTime: formatTime(e.UpdateTime),
			CreateUser: e.CreateUser,
			UpdateUser: e.UpdateUser,
		})
	}

	return &response.PageResult{
		Records: items,
		Total:   total,
	}, nil
}

func formatTime(t time.Time) string {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return t.Format(employeeTimeLayout)
	}
	return t.In(loc).Format(employeeTimeLayout)
}

func maskIDNumber(idNumber string) string {
	if len(idNumber) < 8 {
		return idNumber
	}
	head := idNumber[:4]
	tail := idNumber[len(idNumber)-4:]
	return head + strings.Repeat("*", len(idNumber)-8) + tail
}
