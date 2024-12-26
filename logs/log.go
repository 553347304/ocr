package logs

import (
	"fmt"
	"log"
	"reflect"
)

func Structs(s interface{}) {
	val := reflect.ValueOf(s)
	subStructSlice(val)
	subStruct(val)
}

func subStructSlice(val reflect.Value) bool {
	if val.Kind() == reflect.Slice {
		for i := 0; i < val.Len(); i++ {
			is1 := subStruct(val.Index(i))
			is2 := subStructSlice(val.Index(i))
			if is1 || is2 {
				continue
			}
			return true
		}
		return false
	}
	return false
}
func subStruct(val reflect.Value) bool {
	typ := val.Type()
	if val.Kind() == reflect.Struct {
		// 遍历结构体字段
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			fieldName := typ.Field(i).Name
			log.Print(fmt.Sprintf("%s: %v", fieldName, field.Interface()))
			is1 := subStruct(field)
			is2 := subStructSlice(field)
			if is1 || is2 {
				continue
			}
		}
		return true
	}
	return false
}
