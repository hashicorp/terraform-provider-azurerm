package utils

import (
	"reflect"
	"runtime"
	"testing"
)

func checkTfType(o interface{}, t reflect.Kind) bool {
	out := o.([]interface{})
	for _, i := range out {
		if reflect.ValueOf(i).Kind() != t {
			return false
		}
	}
	return true
}

func compareSlices(s1, s2 interface{}) bool {
	v1 := reflect.ValueOf(s1)
	v2 := reflect.ValueOf(s2)
	if v1.Type() != v2.Type() {
		return false
	}
	if v1.Kind() == reflect.Ptr {
		v1 = v1.Elem()
	}
	if v2.Kind() == reflect.Ptr {
		v2 = v2.Elem()
	}
	if v1.Len() != v2.Len() {
		return false
	}
	for i := 0; i < v1.Len(); i++ {
		if v1.Index(i).Interface() != v2.Index(i).Interface() {
			return false
		}
	}
	return true
}

func TestExpandFlattenFunctions(t *testing.T) {
	cases := []struct {
		f      interface{}
		input  interface{}
		expect interface{}
		tfType reflect.Kind // for testing flattened to type aligned with terraform schema type
		// (not used for expand, in which case it equals to reflect.Invalid)
	}{
		{
			f:      ExpandStringSlice,
			input:  []interface{}{"a", "b"},
			expect: &[]string{"a", "b"},
			tfType: reflect.Invalid,
		},
		{
			f:      FlattenStringSlice,
			input:  &[]string{"a", "b"},
			expect: []interface{}{"a", "b"},
			tfType: reflect.String,
		},
		{
			f:      ExpandBoolSlice,
			input:  []interface{}{true, false},
			expect: &[]bool{true, false},
			tfType: reflect.Invalid,
		},
		{
			f:      FlattenBoolSlice,
			input:  &[]bool{true, false},
			expect: []interface{}{true, false},
			tfType: reflect.Bool,
		},
		{
			f:      ExpandUintSlice,
			input:  []interface{}{1, 2, 3},
			expect: &[]uint{1, 2, 3},
			tfType: reflect.Invalid,
		},
		{
			f:      FlattenUintSlice,
			input:  &[]uint{1, 2, 3},
			expect: []interface{}{1, 2, 3},
			tfType: reflect.Int,
		},
		{
			f:      ExpandUint8Slice,
			input:  []interface{}{1, 2, 3},
			expect: &[]uint8{1, 2, 3},
			tfType: reflect.Invalid,
		},
		{
			f:      FlattenUint8Slice,
			input:  &[]uint8{1, 2, 3},
			expect: []interface{}{1, 2, 3},
			tfType: reflect.Int,
		},
		{
			f:      ExpandUint16Slice,
			input:  []interface{}{1, 2, 3},
			expect: &[]uint16{1, 2, 3},
			tfType: reflect.Invalid,
		},
		{
			f:      FlattenUint16Slice,
			input:  &[]uint16{1, 2, 3},
			expect: []interface{}{1, 2, 3},
			tfType: reflect.Int,
		},
		{
			f:      ExpandUint32Slice,
			input:  []interface{}{1, 2, 3},
			expect: &[]uint32{1, 2, 3},
			tfType: reflect.Invalid,
		},
		{
			f:      FlattenUint32Slice,
			input:  &[]uint32{1, 2, 3},
			expect: []interface{}{1, 2, 3},
			tfType: reflect.Int,
		},
		{
			f:      ExpandUint64Slice,
			input:  []interface{}{1, 2, 3},
			expect: &[]uint64{1, 2, 3},
			tfType: reflect.Invalid,
		},
		{
			f:      FlattenUint64Slice,
			input:  &[]uint64{1, 2, 3},
			expect: []interface{}{1, 2, 3},
			tfType: reflect.Int,
		},
		{
			f:      ExpandIntSlice,
			input:  []interface{}{1, 2, 3},
			expect: &[]int{1, 2, 3},
			tfType: reflect.Invalid,
		},
		{
			f:      FlattenIntSlice,
			input:  &[]int{1, 2, 3},
			expect: []interface{}{1, 2, 3},
			tfType: reflect.Int,
		},
		{
			f:      ExpandInt8Slice,
			input:  []interface{}{1, 2, 3},
			expect: &[]int8{1, 2, 3},
			tfType: reflect.Invalid,
		},
		{
			f:      FlattenInt8Slice,
			input:  &[]int8{1, 2, 3},
			expect: []interface{}{1, 2, 3},
			tfType: reflect.Int,
		},
		{
			f:      ExpandInt16Slice,
			input:  []interface{}{1, 2, 3},
			expect: &[]int16{1, 2, 3},
			tfType: reflect.Invalid,
		},
		{
			f:      FlattenInt16Slice,
			input:  &[]int16{1, 2, 3},
			expect: []interface{}{1, 2, 3},
			tfType: reflect.Int,
		},
		{
			f:      ExpandInt32Slice,
			input:  []interface{}{1, 2, 3},
			expect: &[]int32{1, 2, 3},
			tfType: reflect.Invalid,
		},
		{
			f:      FlattenInt32Slice,
			input:  &[]int32{1, 2, 3},
			expect: []interface{}{1, 2, 3},
			tfType: reflect.Int,
		},
		{
			f:      ExpandInt64Slice,
			input:  []interface{}{1, 2, 3},
			expect: &[]int64{1, 2, 3},
			tfType: reflect.Invalid,
		},
		{
			f:      FlattenInt64Slice,
			input:  &[]int64{1, 2, 3},
			expect: []interface{}{1, 2, 3},
			tfType: reflect.Int,
		},
		{
			f:      ExpandFloat32Slice,
			input:  []interface{}{1.0, 2.0, 3.0},
			expect: &[]float32{1.0, 2.0, 3.0},
			tfType: reflect.Invalid,
		},
		{
			f:      FlattenFloat32Slice,
			input:  &[]float32{1.0, 2.0, 3.0},
			expect: []interface{}{1.0, 2.0, 3.0},
			tfType: reflect.Float64,
		},
		{
			f:      ExpandFloat64Slice,
			input:  []interface{}{1.0, 2.0, 3.0},
			expect: &[]float64{1.0, 2.0, 3.0},
			tfType: reflect.Invalid,
		},
		{
			f:      FlattenFloat64Slice,
			input:  &[]float64{1.0, 2.0, 3.0},
			expect: []interface{}{1.0, 2.0, 3.0},
			tfType: reflect.Float64,
		},
	}

	for _, c := range cases {
		vf := reflect.ValueOf(c.f)
		outv := vf.Call([]reflect.Value{reflect.ValueOf(c.input)})
		out := outv[0].Interface()
		if !compareSlices(out, c.expect) {
			t.Fatalf("Function %s failed:\nOutput: %v\n(Expected: %v)\n", runtime.FuncForPC(vf.Pointer()).Name(), out, c.expect)
		}
		if c.tfType != reflect.Invalid && !checkTfType(out, c.tfType) {
			t.Fatalf("Function %s failed:\nElement type in output isn't as expected(%s)\n(output: %v)", runtime.FuncForPC(vf.Pointer()).Name(), c.tfType, out)
		}
	}
}
