package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"sky-take-out-go/internal/config"
	"sky-take-out-go/internal/middleware"
	"sky-take-out-go/internal/model"
	"sky-take-out-go/internal/pkg/errcode"
	"sky-take-out-go/internal/pkg/req"
	"sky-take-out-go/internal/pkg/response"
)

type mockEmployeeRepo struct {
	createErr error
	pageResp  *response.PageResult
	pageErr   error
}

func (m *mockEmployeeRepo) Create(_ context.Context, _ *model.Employee) error {
	return m.createErr
}

func (m *mockEmployeeRepo) GetByUsername(_ context.Context, _ string) (*model.Employee, error) {
	return nil, errors.New("not implemented")
}

func (m *mockEmployeeRepo) PageQuery(_ context.Context, _ *req.EmployeePageReq) (*response.PageResult, error) {
	return m.pageResp, m.pageErr
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

func TestEmployeeServicePageQuery_FormatTimeAndMaskIDNumber(t *testing.T) {
	repo := &mockEmployeeRepo{
		pageResp: &response.PageResult{
			Records: []model.Employee{
				{
					ID:         7,
					Username:   "u1",
					Name:       "张三",
					Phone:      "13800000000",
					Sex:        "男",
					IDNumber:   "110105199003078888",
					Status:     model.EmployeeStatusEnabled,
					CreateTime: time.Date(2026, 4, 6, 10, 30, 0, 0, time.UTC),
					UpdateTime: time.Date(2026, 4, 6, 11, 0, 0, 0, time.UTC),
					CreateUser: 1,
					UpdateUser: 1,
				},
			},
			Total: 1,
		},
	}

	svc := NewEmployeeService(repo, &config.Config{JWT: config.JWTConfig{Secret: "test", Expire: 3600}})
	page, err := svc.PageQuery(context.Background(), &req.EmployeePageReq{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	items, ok := page.Records.([]response.EmployeePageItemDTO)
	if !ok {
		t.Fatalf("expected []response.EmployeePageItemDTO, got %T", page.Records)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].IDNumber != "1101**********8888" {
		t.Fatalf("unexpected masked id number: %s", items[0].IDNumber)
	}
	if items[0].CreateTime != "2026-04-06 18:30:00" {
		t.Fatalf("unexpected create time: %s", items[0].CreateTime)
	}
}
