package utils

import (
	"fmt"
	"reflect"
)

func StructToMap(stc interface{}) map[string]interface{} {
	newMap := make(map[string]interface{})
	t := reflect.TypeOf(stc)
	v := reflect.ValueOf(stc)
	fields := t.NumField()
	for i := 0; i < fields; i++ {
		key := t.Field(i).Name
		// 解析注解key
		if t.Field(i).Tag.Get("json") != "" {
			key = t.Field(i).Tag.Get("json")
		}
		// 解析结构体类型
		if v.Field(i).Kind() == reflect.Struct {
			newMap[key] = StructToMap(v.Field(i).Interface())
			continue
		}
		//	解析指针类型
		if v.Field(i).Kind() == reflect.Ptr {
			newMap[key] = StructToMap(v.Field(i).Elem().Interface())
			continue
		}
		// 解析基本类型
		newMap[key] = convert(v.Field(i))

	}
	return newMap
}

func convert(field reflect.Value) interface{} {
	// todo 其它类型自行支持
	switch field.Type().Name() {
	case reflect.String.String():
		return field.String()
	case reflect.Int.String(), reflect.Int64.String():
		return field.Int()
	case reflect.Int8.String():
		return int8(field.Int())
	case reflect.Float32.String():
		return float32(field.Float())
	case reflect.Float64.String():
		return field.Float()
	case reflect.Complex64.String():
		return complex64(field.Complex())
	case reflect.Complex128.String():
		return field.Complex()
		return float32(field.Float())
	default:
		panic(fmt.Sprintf("未知的类型%s", field.Type().Kind()))
	}
	return nil
}
