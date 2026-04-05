package validator

import (
	"reflect"
	"regexp"
	"strings"
	"sync"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var regexCache sync.Map // map[string]*regexp.Regexp

func registerRegexpValidation(v *validator.Validate, trans ut.Translator, message string) {
	_ = v.RegisterValidation("regexp", validateRegexp)
	_ = v.RegisterTranslation(
		"regexp",
		trans,
		func(ut ut.Translator) error {
			return ut.Add("regexp", message, true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, err := ut.T("regexp", fe.Field())
			if err != nil {
				return fe.Field() + " has an invalid format"
			}
			return t
		},
	)
}

func validateRegexp(fl validator.FieldLevel) bool {
	pattern := normalizeRegexpParam(fl.Param())
	if pattern == "" {
		return false
	}

	field := fl.Field()
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return true
		}
		field = field.Elem()
	}
	if field.Kind() != reflect.String {
		return false
	}

	val := field.String()

	var re *regexp.Regexp
	if cached, ok := regexCache.Load(pattern); ok {
		re = cached.(*regexp.Regexp)
	} else {
		compiled, err := regexp.Compile(pattern)
		if err != nil {
			return false
		}
		regexCache.Store(pattern, compiled)
		re = compiled
	}

	return re.MatchString(val)
}

func normalizeRegexpParam(param string) string {
	// validator/v10 会把 '|' 作为规则分隔符，tag 中需用 0x7C 替代；
	// 同理若正则里需要 ','，可用 0x2C 替代。
	replacer := strings.NewReplacer(
		"0x7C", "|",
		"0X7C", "|",
		"0x2C", ",",
		"0X2C", ",",
	)
	return replacer.Replace(param)
}
