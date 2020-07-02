package utils

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestExpandSlice(t *testing.T) {
	type T1 string
	type T2 struct {
		S string
	}

	cases := []struct {
		input   []interface{}
		t       interface{}
		convert func(interface{}) interface{}
		output  interface{}
	}{
		// slice of string -> slice of string
		{
			input:  []interface{}{"a", "b"},
			t:      "",
			output: ToPtr([]string{"a", "b"}),
		},
		// slice of string -> slice of customized string type
		{
			input: []interface{}{"a", "b"},
			t:     T1(""),
			convert: func(x interface{}) interface{} {
				return T1(x.(string))
			},
			output: ToPtr([]T1{"a", "b"}),
		},
		// slice of string -> slice of customized structure containing string member
		{
			input: []interface{}{"a", "b"},
			t:     T2{},
			convert: func(x interface{}) interface{} {
				return T2{S: x.(string)}
			},
			output: ToPtr([]T2{{"a"}, {"b"}}),
		},
		// slice of int -> slice of int32
		{
			input: []interface{}{1, 2},
			t:     int32(0),
			convert: func(x interface{}) interface{} {
				return int32(x.(int))
			},
			output: ToPtr([]int32{1, 2}),
		},
		// slice of int -> slice of int64
		{
			input: []interface{}{1, 2},
			t:     int64(0),
			convert: func(x interface{}) interface{} {
				return int64(x.(int))
			},
			output: ToPtr([]int64{1, 2}),
		},
		// slice of int -> slice of int
		{
			input:  []interface{}{1, 2},
			t:      0,
			output: ToPtr([]int{1, 2}),
		},
		// slice of float64 -> slice of float32
		{
			input: []interface{}{float64(1), float64(2)},
			t:     float32(0),
			convert: func(x interface{}) interface{} {
				return float32(x.(float64))
			},
			output: ToPtr([]float32{1, 2}),
		},
		// slice of float64 -> slice of float64
		{
			input:  []interface{}{float64(1), float64(2)},
			t:      float64(0),
			output: ToPtr([]float64{1, 2}),
		},
		// slice of string contains nil -> slice of string
		{
			input:  []interface{}{"a", "b", nil},
			t:      "",
			output: ToPtr([]string{"a", "b", ""}),
		},
		// empty slice of string -> slice of string
		{
			input:  []interface{}{},
			t:      "",
			output: ToPtr([]string{}),
		},
	}

	for idx, c := range cases {
		out := ExpandSlice(c.input, c.t, c.convert)
		if !reflect.DeepEqual(out, c.output) {
			t.Fatalf("%d failed\nexpected:\n%s\nactual:\n%s\n", idx, spew.Sdump(c.output), spew.Sdump(out))
		}
	}
}

func TestExpandMap(t *testing.T) {
	type T1 string
	type T2 struct {
		S string
	}

	cases := []struct {
		input   map[string]interface{}
		t       interface{}
		convert func(interface{}) interface{}
		output  interface{}
	}{
		// map[string]string -> map[string]string
		{
			input: map[string]interface{}{
				"a": "b",
			},
			t: "",
			output: map[string]string{
				"a": "b",
			},
		},
		// map[string]string -> map[string](customized string type)
		{
			input: map[string]interface{}{
				"a": "b",
			},
			t: T1(""),
			convert: func(x interface{}) interface{} {
				return T1(x.(string))
			},
			output: map[string]T1{
				"a": "b",
			},
		},
		// map[string]string -> map[string](customized type structure containing string member)
		{
			input: map[string]interface{}{
				"a": "b",
			},
			t: T2{},
			convert: func(x interface{}) interface{} {
				return T2{S: x.(string)}
			},
			output: map[string]T2{
				"a": {"b"},
			},
		},
		// map[string]int -> map[string]int32
		{
			input: map[string]interface{}{
				"a": 1,
			},
			t: int32(0),
			convert: func(x interface{}) interface{} {
				return int32(x.(int))
			},
			output: map[string]int32{
				"a": 1,
			},
		},
		// map[string]int -> map[string]int64
		{
			input: map[string]interface{}{
				"a": 1,
			},
			t: int64(0),
			convert: func(x interface{}) interface{} {
				return int64(x.(int))
			},
			output: map[string]int64{
				"a": 1,
			},
		},
		// map[string]int -> map[string]int
		{
			input: map[string]interface{}{
				"a": 1,
			},
			t: 0,
			output: map[string]int{
				"a": 1,
			},
		},
		// map[string]float64 -> map[string]float32
		{
			input: map[string]interface{}{
				"a": float64(0),
			},
			t: float32(0),
			convert: func(x interface{}) interface{} {
				return float32(x.(float64))
			},
			output: map[string]float32{
				"a": 0,
			},
		},
		// map[string]float64 -> map[string]float64
		{
			input: map[string]interface{}{
				"a": float64(0),
			},
			t: float64(0),
			output: map[string]float64{
				"a": 0,
			},
		},
		// map[string]string contains nil -> map[string]string
		{
			input: map[string]interface{}{
				"a": "b",
				"c": nil,
			},
			t: "",
			output: map[string]string{
				"a": "b",
				"c": "",
			},
		},
		// empty map[string]string -> map[string]string
		{
			input:  map[string]interface{}{},
			t:      "",
			output: map[string]string{},
		},
	}

	for idx, c := range cases {
		out := ExpandMap(c.input, c.t, c.convert)
		if !reflect.DeepEqual(out, c.output) {
			t.Fatalf("%d failed\nexpected:\n%s\nactual:\n%s\n", idx, spew.Sdump(c.output), spew.Sdump(out))
		}
	}
}

func TestFlattenSlice(t *testing.T) {
	type T1 string
	type T2 struct {
		S string
	}

	cases := []struct {
		input   interface{}
		convert func(interface{}) interface{}
		output  []interface{}
	}{
		// slice of string -> slice of string
		{
			input:  ToPtr([]string{"a", "b"}),
			output: []interface{}{"a", "b"},
		},
		// slice of customized string type -> slice of string
		{
			input: ToPtr([]T1{"a", "b"}),
			convert: func(x interface{}) interface{} {
				return string(x.(T1))
			},
			output: []interface{}{"a", "b"},
		},
		// slice of customized structure containing string member -> slice of string
		{
			input: ToPtr([]T2{{"a"}, {"b"}}),
			convert: func(x interface{}) interface{} {
				return x.(T2).S
			},
			output: []interface{}{"a", "b"},
		},
		// slice of int32 -> slice of int
		{
			input: ToPtr([]int32{1, 2}),
			convert: func(x interface{}) interface{} {
				return int(x.(int32))
			},
			output: []interface{}{1, 2},
		},
		// slice of int64 -> slice of int
		{
			input: ToPtr([]int64{1, 2}),
			convert: func(x interface{}) interface{} {
				return int(x.(int64))
			},
			output: []interface{}{1, 2},
		},
		// slice of int -> slice of int
		{
			input:  ToPtr([]int{1, 2}),
			output: []interface{}{1, 2},
		},
		//  slice of float32 -> slice of float64
		{
			input: ToPtr([]float32{1, 2}),
			convert: func(x interface{}) interface{} {
				return float64(x.(float32))
			},
			output: []interface{}{1.0, 2.0},
		},
		// slice of float64 -> slice of float64
		{
			input:  ToPtr([]float64{1, 2}),
			output: []interface{}{1.0, 2.0},
		},
		// slice of string pointer -> slice of string
		{
			input: ToPtr([]*string{String("a"), String("b")}),
			convert: func(x interface{}) interface{} {
				return *(x.(*string))
			},
			output: []interface{}{"a", "b"},
		},
		// slice of string pointer contains nil -> slice of string
		{
			input: ToPtr([]*string{String("a"), String("b"), nil}),
			convert: func(x interface{}) interface{} {
				return *(x.(*string))
			},
			output: []interface{}{"a", "b", ""},
		},
		// empty slice of string pointer -> slice of string
		{
			input:  ToPtr([]*string{}),
			output: []interface{}{},
		},
	}

	for idx, c := range cases {
		out := FlattenSlicePtr(c.input, c.convert)
		if !reflect.DeepEqual(out, c.output) {
			t.Fatalf("%d failed\nexpected:\n%s\nactual:\n%s\n", idx, spew.Sdump(c.output), spew.Sdump(out))
		}
	}
}

func TestFlattenSliceGuard(t *testing.T) {
	inputs := []interface{}{
		[]int{},
		1,
		[]string{},
		[]*string{},
		String("a"),
	}
	for idx, input := range inputs {
		shouldPanic(t, func() { FlattenSlicePtr(input, nil) }, idx)
	}
}

func TestFlattenStringMap(t *testing.T) {
	type T1 string
	type T2 struct {
		S string
	}

	cases := []struct {
		input   interface{}
		convert func(interface{}) interface{}
		output  map[string]interface{}
	}{
		// map[string]string -> map[string]string
		{
			input: map[string]string{
				"a": "b",
			},
			output: map[string]interface{}{
				"a": "b",
			},
		},
		// map[string](customized string type) -> map[string]string
		{
			input: map[string]T1{
				"a": "b",
			},
			convert: func(x interface{}) interface{} {
				return string(x.(T1))
			},
			output: map[string]interface{}{
				"a": "b",
			},
		},
		// map[string](customized type structure containing string member) -> map[string]string
		{
			input: map[string]T2{
				"a": {"b"},
			},
			convert: func(x interface{}) interface{} {
				return x.(T2).S
			},
			output: map[string]interface{}{
				"a": "b",
			},
		},
		// map[string]int32 -> map[string]int
		{
			input: map[string]int32{
				"a": 1,
			},
			convert: func(x interface{}) interface{} {
				return int(x.(int32))
			},
			output: map[string]interface{}{
				"a": 1,
			},
		},
		// map[string]int64 -> map[string]int
		{
			input: map[string]int64{
				"a": 1,
			},
			convert: func(x interface{}) interface{} {
				return int(x.(int64))
			},
			output: map[string]interface{}{
				"a": 1,
			},
		},
		// map[string]int -> map[string]int
		{
			input: map[string]int{
				"a": 1,
			},
			output: map[string]interface{}{
				"a": 1,
			},
		},
		// map[string]float32 -> map[string]float64
		{
			input: map[string]float32{
				"a": 0,
			},
			convert: func(x interface{}) interface{} {
				return float64(x.(float32))
			},
			output: map[string]interface{}{
				"a": 0.0,
			},
		},
		// map[string]float64 -> map[string]float64
		{
			input: map[string]float64{
				"a": 0,
			},
			output: map[string]interface{}{
				"a": 0.0,
			},
		},
		// map[string]*string -> map[string]string
		{
			input: map[string]*string{
				"a": String("b"),
			},
			convert: func(x interface{}) interface{} {
				return *(x.(*string))
			},
			output: map[string]interface{}{
				"a": "b",
			},
		},
		// map[string]*string contains nil -> map[string]string
		{
			input: map[string]*string{
				"a": String("b"),
				"c": nil,
			},
			convert: func(x interface{}) interface{} {
				return *(x.(*string))
			},
			output: map[string]interface{}{
				"a": "b",
				"c": "",
			},
		},
		// empty map[string]*string -> map[string]string
		{
			input:  map[string]*string{},
			output: map[string]interface{}{},
		},
	}

	for idx, c := range cases {
		out := FlattenStringMap(c.input, c.convert)
		if !reflect.DeepEqual(out, c.output) {
			t.Fatalf("%d failed\nexpected:\n%s\nactual:\n%s\n", idx, spew.Sdump(c.output), spew.Sdump(out))
		}
	}
}

func TestFlattenStringMapGuard(t *testing.T) {
	inputs := []interface{}{
		[]int{},
		1,
		String("a"),
		map[int]interface{}{},
	}
	for idx, input := range inputs {
		shouldPanic(t, func() { FlattenStringMap(input, nil) }, idx)
	}
}

func shouldPanic(t *testing.T, f func(), id int) {
	defer func() { recover() }()
	f()
	t.Errorf("%d should have panicked", id)
}
