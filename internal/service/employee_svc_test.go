package service

import (
	"context"
	"errors"
	"testing"

	"sky-take-out-go/internal/config"
	"sky-take-out-go/internal/middleware"
	"sky-take-out-go/internal/model"
	"sky-take-out-go/internal/pkg/errcode"
	"sky-take-out-go/internal/pkg/req"
)

type mockEmployeeRepo struct {
	createErr error
}

func (m *mockEmployeeRepo) Create(_ context.Context, _ *model.Employee) error {
	return m.createErr
}

func (m *mockEmployeeRepo) GetByUsername(_ context.Context, _ string) (*model.Employee, error) {
	return nil, errors.New("not implemented")
}

func buildEmployeeCreateCtx() context.Context {
	claims := &middleware.CustomClaims{
		UserID:   1,
		Username: "admin",
		Role:     middleware.RoleAdmin,
	}
	return context.WithValue(context.Background(), string(middleware.ContextKeyClaims), claims)
}

func TestEmployeeServiceCreate_PropagateAppError(t *testing.T) {
	repo := &mockEmployeeRepo{createErr: errcode.ErrUserNameAlreadyExists()}
	svc := NewEmployeeService(repo, &config.Config{JWT: config.JWTConfig{Secret: "test", Expire: 3600}})

	err := svc.Create(buildEmployeeCreateCtx(), &req.EmployeeCreateReq{
		Username: "admin",
		Name:     "管理员",
		Phone:    "13800000000",
		Sex:      "男",
		IDNumber: "110105199003078888",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	var appErr *errcode.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T", err)
	}
	if appErr.Code != errcode.CodeUserNameAlreadyExists {
		t.Fatalf("expected code %d, got %d", errcode.CodeUserNameAlreadyExists, appErr.Code)
	}
}

func TestEmployeeServiceCreate_MapUnknownErrorToInternal(t *testing.T) {
	repo := &mockEmployeeRepo{createErr: errors.New("db down")}
	svc := NewEmployeeService(repo, &config.Config{JWT: config.JWTConfig{Secret: "test", Expire: 3600}})

	err := svc.Create(buildEmployeeCreateCtx(), &req.EmployeeCreateReq{
		Username: "admin",
		Name:     "管理员",
		Phone:    "13800000000",
		Sex:      "男",
		IDNumber: "110105199003078888",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	var appErr *errcode.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T", err)
	}
	if appErr.Code != errcode.CodeInternal {
		t.Fatalf("expected code %d, got %d", errcode.CodeInternal, appErr.Code)
	}
}
