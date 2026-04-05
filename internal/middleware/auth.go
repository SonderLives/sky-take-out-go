package middleware

import (
	"context"
	"fmt"
	"sky-take-out-go/internal/pkg/errcode"
	"sky-take-out-go/internal/pkg/logger"
	"sky-take-out-go/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const tokenHeader = "Token"

type AuthMiddleware struct {
	jwtSecret []byte
}

func NewAuthMiddleware(secret string) *AuthMiddleware {
	return &AuthMiddleware{jwtSecret: []byte(secret)}
}

func (m *AuthMiddleware) parseToken(c *gin.Context) (*CustomClaims, bool) {
	authHeader := c.GetHeader(tokenHeader)
	if authHeader == "" {
		logger.WithCtx(c.Request.Context()).Warnw("auth token missing",
			"header", tokenHeader,
			"path", c.Request.URL.Path,
			"method", c.Request.Method,
		)
		response.Error(c, errcode.ErrUnauthorized())
		c.Abort()
		return nil, false
	}

	//parts := strings.SplitN(authHeader, " ", 2)
	//if len(parts) != 2 || parts[0] != "Bearer" {
	//	response.Error(c, errcode.ErrUnauthorized())
	//	c.Abort()
	//	return nil, false
	//}

	if len(m.jwtSecret) == 0 {
		logger.WithCtx(c.Request.Context()).Errorw("jwt secret is empty")
		response.Error(c, errcode.ErrUnauthorized())
		c.Abort()
		return nil, false
	}

	token, err := jwt.ParseWithClaims(authHeader, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return m.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		logger.WithCtx(c.Request.Context()).Warnw("auth token invalid",
			"path", c.Request.URL.Path,
			"method", c.Request.Method,
			"error", err,
		)
		response.Error(c, errcode.ErrUnauthorized())
		c.Abort()
		return nil, false
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		logger.WithCtx(c.Request.Context()).Warnw("auth token claims type invalid")
		response.Error(c, errcode.ErrUnauthorized())
		c.Abort()
		return nil, false
	}

	return claims, true
}

func (m *AuthMiddleware) AppAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := m.parseToken(c)
		if !ok {
			return
		}

		if claims.Role != RoleUser {
			logger.WithCtx(c.Request.Context()).Warnw("app auth role forbidden",
				"role", claims.Role,
				"path", c.Request.URL.Path,
			)
			response.Error(c, errcode.ErrForbidden())
			c.Abort()
			return
		}

		ctx := context.WithValue(c.Request.Context(), string(ContextKeyClaims), claims)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func (m *AuthMiddleware) AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := m.parseToken(c)
		if !ok {
			return
		}

		if claims.Role != RoleAdmin {
			logger.WithCtx(c.Request.Context()).Warnw("admin auth role forbidden",
				"role", claims.Role,
				"path", c.Request.URL.Path,
			)
			response.Error(c, errcode.ErrForbidden())
			c.Abort()
			return
		}

		ctx := context.WithValue(c.Request.Context(), string(ContextKeyClaims), claims)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
