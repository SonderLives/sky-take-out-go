package service

import (
	"context"
	"errors"
	"sky-take-out-go/internal/pkg/req"
	"time"

	"sky-take-out-go/internal/config"
	"sky-take-out-go/internal/middleware"
	"sky-take-out-go/internal/model"
	"sky-take-out-go/internal/pkg/errcode"
	"sky-take-out-go/internal/pkg/response"
	"sky-take-out-go/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const defaultEmployeePassword = "123456"

type EmployeeService interface {
	Create(ctx context.Context, r *req.EmployeeCreateReq) error
	Login(ctx context.Context, username, password string) (*response.EmployeeLoginResult, error)
}

type employeeService struct {
	repo      repository.EmployeeRepository
	jwtSecret []byte
	jwtExpire time.Duration
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

// Create 创建员工
func (s *employeeService) Create(ctx context.Context, r *req.EmployeeCreateReq) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(defaultEmployeePassword), bcrypt.DefaultCost)
	if err != nil {
		return errcode.ErrInternal()
	}

	// 安全断言，避免 claims 为 nil 时 panic
	claims, ok := ctx.Value(string(middleware.ContextKeyClaims)).(*middleware.CustomClaims)
	if !ok || claims == nil {
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
			return appErr
		}
		return errcode.ErrInternal()
	}
	return nil
}

// Login 员工登录
func (s *employeeService) Login(ctx context.Context, username, password string) (*response.EmployeeLoginResult, error) {
	employee, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrUserNotFound()
		}
		return nil, errcode.ErrInternal()
	}
	if err := bcrypt.CompareHashAndPassword([]byte(employee.Password), []byte(password)); err != nil {
		return nil, errcode.ErrAdminOrPassword()
	}

	if employee.Status == model.EmployeeStatusDisabled {
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
		return nil, errcode.ErrInternal()
	}

	return &response.EmployeeLoginResult{
		Token:    tokenStr,
		ID:       employee.ID,
		Name:     employee.Name,
		UserName: employee.Username,
	}, nil
}
