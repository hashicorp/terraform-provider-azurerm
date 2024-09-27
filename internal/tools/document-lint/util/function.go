// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"reflect"
	"runtime"
)

func FuncFileLine(f interface{}) (file string, line int) {
	vf, ok := f.(reflect.Value)
	if !ok {
		vf = reflect.ValueOf(f)
	}
	if vf.IsNil() {
		return
	}
	return PointerFileLine(vf)
}

func PointerFileLine(v reflect.Value) (file string, line int) {
	pc := v.Pointer()
	file, line = runtime.FuncForPC(pc).FileLine(pc)
	return
}
