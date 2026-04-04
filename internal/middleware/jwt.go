package middleware

import "github.com/golang-jwt/jwt/v5"

// Role 定义角色常量
type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// ContextKey 定义上下文中的键
type ContextKey string

const (
	ContextKeyClaims ContextKey = "claims"
)

// CustomClaims 自定义 JWT Claims，更类型安全
type CustomClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     Role   `json:"role"`
	jwt.RegisteredClaims
}
