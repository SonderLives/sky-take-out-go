package service

import (
	"context"
	"errors"
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

type EmployeeService interface {
	CreateAdmin(ctx context.Context, username, password, nickname, role string) error
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

// CreateAdmin 创建管理员
func (s *employeeService) CreateAdmin(ctx context.Context, username, password, nickname, role string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errcode.ErrInternal()
	}
	return s.repo.Create(ctx, &model.Employee{
		Username: username,
		Password: string(hashed),
		Status:   1,
	})
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

	if employee.Status != 1 {
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
