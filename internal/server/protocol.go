package server

import (
	"fmt"
	"reflect"
	"time"
)

func DataConversion(v interface{}) interface{} {
	val := reflect.ValueOf(v)

	switch val.Kind() {
	case reflect.Ptr:
		return DataConversion(val.Elem().Interface())
	case reflect.Struct:
		fmt.Printf("val: %v\n", val)
		t := reflect.TypeOf(val)
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)

			if field.PkgPath == "" {
				fmt.Printf("field: %v\n", field.Name)
				fmt.Printf("field: %v\n", reflect.ValueOf(val).Field(i).Interface())
				fmt.Printf("field.PkgPath: %v\n", field.PkgPath)
			}
		}
	case reflect.Map:
		for _, key := range val.MapKeys() {
			mapVal := val.MapIndex(key)
			val.SetMapIndex(key, reflect.ValueOf(DataConversion(mapVal.Interface())))
		}
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			sliceVal := val.Index(i)
			sliceVal.Set(reflect.ValueOf(DataConversion(sliceVal.Interface())))
		}
	case reflect.Interface:
		return DataConversion(val.Elem().Interface())
	}
	if val.Type().String() == "time.Time" {
		return val.Interface().(time.Time).Unix()
	}
	return v
}
