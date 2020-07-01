package utils

import "reflect"

func Bool(input bool) *bool {
	return &input
}

func Int(input int) *int {
	return &input
}

func Int32(input int32) *int32 {
	return &input
}

func Int64(input int64) *int64 {
	return &input
}

func Float(input float64) *float64 {
	return &input
}

func String(input string) *string {
	return &input
}

// ToPtr create a new object from the passed in "obj" and return its address back.
func ToPtr(obj interface{}) interface{} {
	v := reflect.ValueOf(obj)
	vp := reflect.New(v.Type())
	vp.Elem().Set(v)
	return vp.Interface()
}
