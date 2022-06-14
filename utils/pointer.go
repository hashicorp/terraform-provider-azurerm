package utils

import (
	"reflect"
)

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

// Ptr input MUST NOT be nil object
func Ptr[T any](input T) *T {
	//v := reflect.ValueOf(input)
	return &input
}

// TryPtr input can be nil of any type
func TryPtr[T any](input T) *T {
	v := reflect.ValueOf(input)
	switch v.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Interface:
		if v.IsNil() {
			return nil
		}
	}
	return &input
}

func Value[T any](ptr *T) T {
	if ptr == nil {
		t := new(T)
		return *t
	}
	return *ptr
}
