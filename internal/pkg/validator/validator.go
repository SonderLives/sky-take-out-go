package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"goflow/internal/pkg/i18n"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var translators = make(map[string]ut.Translator)

// Init 初始化多语言验证器
func Init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册 label tag 作为默认字段名（中文）
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get("label")
			if name == "" {
				name = fld.Tag.Get("json")
				if idx := strings.Index(name, ","); idx != -1 {
					name = name[:idx]
				}
			}
			return name
		})

		zhLocale := zh.New()
		enLocale := en.New()
		uni := ut.New(zhLocale, zhLocale, enLocale)

		zhTrans, _ := uni.GetTranslator("zh")
		_ = zh_translations.RegisterDefaultTranslations(v, zhTrans)
		translators[i18n.LangZhCN] = zhTrans

		enTrans, _ := uni.GetTranslator("en")
		_ = en_translations.RegisterDefaultTranslations(v, enTrans)
		translators[i18n.LangEnUS] = enTrans
	}
}

// Translate 将 validator 错误翻译为指定语言
func Translate(err error, lang string) string {
	trans, ok := translators[lang]
	if !ok {
		trans = translators[i18n.DefaultLang()]
	}

	// JSON 类型不匹配错误，如前端传字符串但期望 number
	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		tmpl := i18n.GetTemplate("type_mismatch", lang)
		if tmpl == "" {
			tmpl = "field '%s' must be %s, got %s"
		}
		return fmt.Sprintf(tmpl, typeErr.Field, typeErr.Type, typeErr.Value)
	}

	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		msgs := make([]string, 0, len(errs))
		for _, e := range errs {
			msg := e.Translate(trans)
			// 非默认语言时，替换字段名为对应语言的 label
			if lang != i18n.DefaultLang() {
				fieldName := e.Field()
				if label, ok := i18n.GetFieldLabel(e.StructNamespace(), fieldName, lang); ok {
					msg = strings.Replace(msg, fieldName, label, 1)
				}
			}
			msgs = append(msgs, msg)
		}
		return strings.Join(msgs, "; ")
	}
	return err.Error()
}
