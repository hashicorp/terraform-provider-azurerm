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
			output: []string{"a", "b"},
		},
		// slice of string -> slice of customized string type
		{
			input: []interface{}{"a", "b"},
			t:     T1(""),
			convert: func(x interface{}) interface{} {
				return T1(x.(string))
			},
			output: []T1{"a", "b"},
		},
		// slice of string -> slice of customized structure containing string member
		{
			input: []interface{}{"a", "b"},
			t:     T2{},
			convert: func(x interface{}) interface{} {
				return T2{S: x.(string)}
			},
			output: []T2{{"a"}, {"b"}},
		},
		// slice of int -> slice of int32
		{
			input: []interface{}{1, 2},
			t:     int32(0),
			convert: func(x interface{}) interface{} {
				return int32(x.(int))
			},
			output: []int32{1, 2},
		},
		// slice of int -> slice of int64
		{
			input: []interface{}{1, 2},
			t:     int64(0),
			convert: func(x interface{}) interface{} {
				return int64(x.(int))
			},
			output: []int64{1, 2},
		},
		// slice of int -> slice of int
		{
			input:  []interface{}{1, 2},
			t:      0,
			output: []int{1, 2},
		},
		// slice of float64 -> slice of float32
		{
			input: []interface{}{float64(1), float64(2)},
			t:     float32(0),
			convert: func(x interface{}) interface{} {
				return float32(x.(float64))
			},
			output: []float32{1, 2},
		},
		// slice of float64 -> slice of float64
		{
			input:  []interface{}{float64(1), float64(2)},
			t:      float64(0),
			output: []float64{1, 2},
		},
	}

	for idx, c := range cases {
		out := ExpandSlice(c.input, c.t, c.convert)
		if !reflect.DeepEqual(out, c.output) {
			t.Fatalf("%d failed\nexpected:\n%s\nactual:\n%s\n", idx, spew.Sdump(out), spew.Sdump(c.output))
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
		// map[string]int  -> map[string]int32
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
		// map[string]int  -> map[string]int64
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
	}

	for idx, c := range cases {
		out := ExpandMap(c.input, c.t, c.convert)
		if !reflect.DeepEqual(out, c.output) {
			t.Fatalf("%d failed\nexpected:\n%s\nactual:\n%s\n", idx, spew.Sdump(out), spew.Sdump(c.output))
		}
	}
}
