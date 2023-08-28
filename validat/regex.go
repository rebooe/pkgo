package validat

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var regexMap = map[string]*regexp.Regexp{
	"mobile": regexp.MustCompile(`^(\d{11})$`), // 手机号判断
}

func regex(fl validator.FieldLevel) bool {
	val, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	re, ok := regexMap[fl.Param()]
	if !ok {
		return false
	}
	return re.MatchString(val)
}
