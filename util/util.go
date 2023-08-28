package util

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"unsafe"
)

// StructToMap 结构体转 map
func StructToMap(dest map[string]any, src any) error {
	// 使用反射获取源结构体的类型和值
	srcValue := reflect.ValueOf(src)
	if srcValue.Type().Kind() == reflect.Ptr {
		srcValue = srcValue.Elem()
	}
	srcType := srcValue.Type()

	// 确保源类型为结构体
	if srcType.Kind() != reflect.Struct {
		return fmt.Errorf("src is not a struct")
	}

	// 遍历结构体的字段
	for i := 0; i < srcType.NumField(); i++ {
		field := srcType.Field(i)
		fieldName := field.Name
		// 获取结构体字段的值
		fieldValue := srcValue.FieldByName(fieldName)

		// 解析标签
		if v, ok := field.Tag.Lookup("json"); ok {
			if befor, after, found := strings.Cut(v, ","); found {
				if befor == "-" {
					// 忽略字段
					continue
				}
				if after == "omitempty" {
					// 忽略零值
					if fieldValue.IsZero() {
						continue
					}
				}
				if befor != "" {
					fieldName = befor
				}
			} else {
				fieldName = v
			}
		}

		// 将结构体字段的值存储到目标 map 中
		if fieldValue.Kind() == reflect.Pointer {
			fieldValue = fieldValue.Elem()
		}
		dest[fieldName] = fieldValue.Interface()
	}
	return nil
}

// MapToStruct map 转结构体
func MapToStruct(dest any, src map[string]any) error {
	// 使用反射获取目标结构体的类型和值
	destValue := reflect.ValueOf(dest)
	if destValue.Type().Kind() == reflect.Pointer {
		destValue = destValue.Elem()
	}
	destType := destValue.Type()

	// 确保目标类型为结构体指针
	if destType.Kind() != reflect.Struct {
		return fmt.Errorf("dest is not a pointer to struct")
	}

	// 遍历 map 中的键值对
	for key, value := range src {
		// 获取目标结构体字段的值
		fieldValue := destValue.FieldByName(key)

		// 确保目标结构体存在对应的字段
		if fieldValue.IsValid() && fieldValue.CanSet() {
			// 将 map 中的值转换为目标结构体字段的类型，并设置给字段
			fieldValue.Set(reflect.ValueOf(value).Convert(fieldValue.Type()))
		}
	}

	return nil
}

// MapMerge 合并两个 Map
func MapMerge(dest, src map[string]any) {
	// 遍历源 map 中的键值对，将其合并到目标 map
	for key, value := range src {
		dest[key] = value
	}
}

// StructMerge 合并两个结构体
// func StructMerge(dest, src any) error {
// 	// 将源结构体转换为 map
// 	srcMap := make(map[string]any)
// 	destMap := make(map[string]any)

// 	if err := StructToMap(srcMap, src); err != nil {
// 		return err
// 	}
// 	if err := StructToMap(destMap, dest); err != nil {
// 		return err
// 	}

// 	// 合并两个 Map
// 	MapMerge(destMap, srcMap)

// 	// 将合并后的 map 设置回结构体
// 	return MapToStruct(dest, destMap)
// }

func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// Caller 获取格式化堆栈信息
func Caller(skip int) (string, bool) {
	pc, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return "", false
	}
	baseFile := filepath.Base(file)
	Func := runtime.FuncForPC(pc)
	pack, funcName, _ := strings.Cut(Func.Name(), ".")

	return fmt.Sprintf("%s/%s:%d %s()", pack, baseFile, line, funcName), true
}
