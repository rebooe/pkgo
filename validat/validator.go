// 基于 gin 使用的过滤器 validator 的封装
package validat

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = binding.Validator.Engine().(*validator.Validate)

	Validate.RegisterValidation("regex", regex)
}

func Struct(obj any) error {
	if err := Validate.Struct(obj); err != nil {
		return ParseError(err, obj)
	}
	return nil
}

func ParseError(err error, obj any) error {
	// 如果是输入参数无效，则直接返回输入参数错误
	if invalid, ok := err.(*validator.InvalidValidationError); ok {
		return fmt.Errorf("输入参数错误: %s" + invalid.Error())
	}

	// 判断是否验证器错误
	validationErrs, ok := err.(validator.ValidationErrors)
	if ok {
		// 解析默认验证器错误
		if err := defaultErr(validationErrs[0]); err != nil {
			return err
		}

		for _, validationErr := range validationErrs {
			fieldName := validationErr.Field()               // 获取字段名称
			value := validationErr.Value()                   // 获取字段值
			fieldName, index, isArr := fieldArray(fieldName) // 判断是否是数组(返回字段名称,报错的数据索引号)

			if field, ok := reflect.TypeOf(obj).Elem().FieldByName(fieldName); ok { // 通过反射获取filed
				// 获取标签的字段名称设置
				if val, ok := field.Tag.Lookup("form"); ok {
					fieldName = val
				}
				if val, ok := field.Tag.Lookup("json"); ok {
					fieldName = val
				}

				if errorInfo, ok := field.Tag.Lookup("errmsg"); ok { // 获取field对应的 errmsg tag 值
					if isArr {
						return fmt.Errorf("参数 %s[%s]=%v %s", fieldName, index, value, errorInfo)
					}
					return fmt.Errorf("参数 %s=%v %s", fieldName, value, errorInfo)
				}
			}
		}
	}

	// 时间格式解析错误
	parseError, ok := err.(*time.ParseError)
	if ok {
		return fmt.Errorf("时间 %s 格式错误, 格式应为: %s", parseError.Value, parseError.Layout)
	}
	return err
}

// 分割数组字段名和索引
func fieldArray(fileName string) (string, string, bool) {
	if match, _ := regexp.MatchString(`.*\[.*\]`, fileName); !match {
		return fileName, "", false
	}

	name, index, _ := strings.Cut(fileName, "[")
	index = strings.TrimSuffix(index, "]")
	return name, index, true
}

var defaultErrorText = map[string]string{
	"required": "参数 %s 必填%s",
	"datetime": "参数 %s 日期格式错误(示例:%s)",
	"len":      "参数 %s 长度须为%s",
	"max":      "参数 %s 大小或长度须大于%s",
	"min":      "参数 %s 大小或长度须小于%s",
	"oneof":    "参数 %s 只允许输入 %s",
}

// 解析默认错误
func defaultErr(fielderr validator.FieldError) error {
	tag := fielderr.Tag()     // 过滤器名称
	param := fielderr.Param() // 过滤器参数
	field := fielderr.Field() // 参数名称

	if errText, ok := defaultErrorText[tag]; ok {
		return fmt.Errorf(errText, field, param)
	}
	return nil
}
