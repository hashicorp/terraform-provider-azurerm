// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package terraform

import (
	"reflect"

	"github.com/hashicorp/terraform-plugin-testing/internal/configs/hcl2shim"
)

// unknownValueWalk is a reimplementation of the prior walk() logic from
// github.com/mitchellh/reflectwalk and the only walker implemented in this
// module that checked values for hcl2shim.UnknownVariableValue.
//
// Using reflection instead of known logic here is a Go anti-pattern, however
// this logic will be removed in the next major version, so the reflection
// approach is preserved to minimize reimplementation effort.
func unknownValueWalk(v reflect.Value) bool {
	for {
		switch v.Kind() {
		case reflect.Interface:
			v = v.Elem()

			continue
		case reflect.Pointer:
			v = reflect.Indirect(v)

			continue
		}

		break
	}

	switch v.Kind() {
	case reflect.Bool,
		reflect.Complex128,
		reflect.Complex64,
		reflect.Float32,
		reflect.Float64,
		reflect.Int,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Int8,
		reflect.Uint,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uint8,
		reflect.Uintptr,
		reflect.String:
		value := v.Interface()

		return value == hcl2shim.UnknownVariableValue
	case reflect.Map:
		for _, k := range v.MapKeys() {
			value := v.MapIndex(k)

			if foundUnknown := unknownValueWalk(value); foundUnknown {
				return true
			}
		}
	case reflect.Array, reflect.Slice:
		for index := 0; index < v.Len(); index++ {
			value := v.Index(index)

			if foundUnknown := unknownValueWalk(value); foundUnknown {
				return true
			}
		}
	case reflect.Struct:
		for index := 0; index < v.Type().NumField(); index++ {
			value := v.FieldByIndex([]int{index})

			if foundUnknown := unknownValueWalk(value); foundUnknown {
				return true
			}
		}
	default:
		panic("unsupported reflect type: " + v.Kind().String())
	}

	return false
}
