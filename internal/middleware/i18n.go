package middleware

import (
	"goflow/internal/pkg/i18n"

	"github.com/gin-gonic/gin"
)

// I18n 国际化中间件，解析 Accept-Language 并写入 context
func I18n() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := i18n.ParseAcceptLanguage(c.GetHeader("Accept-Language"))
		i18n.SetLang(c, lang)
		c.Next()
	}
}
