package i18n

import (
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	LangZhCN = "zh-CN"
	LangEnUS = "en-US"

	langKey     = "lang"
	defaultLang = LangZhCN
)

// messages: errorCode -> lang -> message
var messages = make(map[int]map[string]string)

// templates: templateKey -> lang -> template
var templates = make(map[string]map[string]string)

// RegisterTemplate 注册模板消息（用于需要动态参数的提示）
func RegisterTemplate(key string, langs map[string]string) {
	templates[key] = langs
}

// GetTemplate 获取指定语言的模板，不存在时 fallback 到默认语言
func GetTemplate(key, lang string) string {
	if t, ok := templates[key]; ok {
		if msg, ok := t[lang]; ok {
			return msg
		}
		if msg, ok := t[defaultLang]; ok {
			return msg
		}
	}
	return ""
}

// fieldLabels: structName.fieldName -> lang -> label
var fieldLabels = make(map[string]map[string]string)

// Register 注册错误码的多语言翻译
func Register(code int, langs map[string]string) {
	messages[code] = langs
}

// RegisterLabels 注册结构体字段的多语言 label
func RegisterLabels(structName string, fields map[string]map[string]string) {
	for field, langs := range fields {
		key := structName + "." + field
		fieldLabels[key] = langs
	}
}

// GetFieldLabel 获取字段的多语言 label
// namespace 格式如 "LoginReq.用户名"，zhLabel 为中文字段名
func GetFieldLabel(namespace string, zhLabel string, lang string) (string, bool) {
	// 先尝试完整 namespace 匹配
	if labels, ok := fieldLabels[namespace]; ok {
		if label, ok := labels[lang]; ok {
			return label, true
		}
	}
	// 从 namespace 中提取 StructName.Field 部分（去掉包名前缀）
	// e.StructNamespace() 可能返回 "LoginReq.用户名" 或更长的路径
	if idx := strings.LastIndex(namespace, "."); idx > 0 {
		// 取 StructName 部分：找倒数第二个 "." 之前的内容
		prefix := namespace[:idx]
		if dotIdx := strings.LastIndex(prefix, "."); dotIdx >= 0 {
			prefix = prefix[dotIdx+1:]
		}
		key := prefix + "." + zhLabel
		if labels, ok := fieldLabels[key]; ok {
			if label, ok := labels[lang]; ok {
				return label, true
			}
		}
	}
	return "", false
}

// GetMessage 根据错误码和语言获取翻译后的消息
func GetMessage(code int, lang string) string {
	if langs, ok := messages[code]; ok {
		if msg, ok := langs[lang]; ok {
			return msg
		}
		// fallback 到默认语言
		if msg, ok := langs[defaultLang]; ok {
			return msg
		}
	}
	return ""
}

// DefaultLang 返回系统默认语言
func DefaultLang() string {
	return defaultLang
}

// SetLang 设置语言到 context
func SetLang(c *gin.Context, lang string) {
	c.Set(langKey, lang)
}

// GetLang 从 context 获取当前语言
func GetLang(c *gin.Context) string {
	if lang, exists := c.Get(langKey); exists {
		if l, ok := lang.(string); ok {
			return l
		}
	}
	return defaultLang
}

// ParseAcceptLanguage 解析 Accept-Language header，返回最匹配的语言
func ParseAcceptLanguage(header string) string {
	if header == "" {
		return defaultLang
	}

	supported := map[string]string{
		"zh":    LangZhCN,
		"zh-cn": LangZhCN,
		"en":    LangEnUS,
		"en-us": LangEnUS,
	}

	// 按逗号分割，取第一个匹配的
	for _, part := range strings.Split(header, ",") {
		tag := strings.TrimSpace(strings.SplitN(part, ";", 2)[0])
		tag = strings.ToLower(tag)
		if lang, ok := supported[tag]; ok {
			return lang
		}
		if idx := strings.Index(tag, "-"); idx > 0 {
			if lang, ok := supported[tag[:idx]]; ok {
				return lang
			}
		}
	}

	return defaultLang
}
