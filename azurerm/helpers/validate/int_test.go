package validate

import (
	"strconv"
	"testing"
)

func TestIntBetweenAndDivisibleBy(t *testing.T) {
	cases := []struct {
		Min    int
		Max    int
		Div    int
		Value  interface{}
		Errors int
	}{
		{
			Min:    1025,
			Max:    2048,
			Div:    1024,
			Value:  1024,
			Errors: 1,
		},
		{
			Min:    1025,
			Max:    2048,
			Div:    3,
			Value:  1024,
			Errors: 1,
		},
		{
			Min:    1024,
			Max:    2048,
			Div:    1024,
			Value:  3072,
			Errors: 1,
		},
		{
			Min:    1024,
			Max:    2048,
			Div:    1024,
			Value:  2049,
			Errors: 1,
		},
		{
			Min:    1024,
			Max:    2048,
			Div:    1024,
			Value:  1024,
			Errors: 0,
		},
	}

	for _, tc := range cases {
		_, errors := IntBetweenAndDivisibleBy(tc.Min, tc.Max, tc.Div)(tc.Value, strconv.Itoa(tc.Value.(int)))
		if len(errors) != tc.Errors {
			t.Fatalf("Expected ValidateIntBetweenAndDivisibleBy to trigger '%d' errors for '%s' - got '%d'", tc.Errors, tc.Value, len(errors))
		}
	}
}

func TestIntDivisibleBy(t *testing.T) {
	cases := []struct {
		Div    int
		Value  interface{}
		Errors int
	}{
		{
			Div:    1024,
			Value:  2048,
			Errors: 0,
		},
		{
			Div:    3,
			Value:  1024,
			Errors: 1,
		},
		{
			Div:    1024,
			Value:  3333,
			Errors: 1,
		},
		{
			Div:    1024,
			Value:  2049,
			Errors: 1,
		},
		{
			Div:    1024,
			Value:  1024,
			Errors: 0,
		},
	}

	for _, tc := range cases {
		_, errors := IntDivisibleBy(tc.Div)(tc.Value, strconv.Itoa(tc.Value.(int)))
		if len(errors) != tc.Errors {
			t.Fatalf("Expected ValidateIntDivisibleBy to trigger '%d' errors for '%s' - got '%d'", tc.Errors, tc.Value, len(errors))
		}
	}
}

func TestIntInSlice(t *testing.T) {

	cases := []struct {
		Input  []int
		Value  int
		Errors int
	}{
		{
			Input:  []int{},
			Value:  0,
			Errors: 1,
		},
		{
			Input:  []int{1},
			Value:  1,
			Errors: 0,
		},
		{
			Input:  []int{1, 2, 3, 4, 5},
			Value:  3,
			Errors: 0,
		},
		{
			Input:  []int{1, 3, 5},
			Value:  3,
			Errors: 0,
		},
		{
			Input:  []int{1, 3, 5},
			Value:  4,
			Errors: 1,
		},
	}

	for _, tc := range cases {
		_, errors := IntInSlice(tc.Input)(tc.Value, "azurerm_postgresql_database")

		if len(errors) != tc.Errors {
			t.Fatalf("Expected the validateIntInSlice trigger a validation error for input: %+v looking for %+v", tc.Input, tc.Value)
		}
	}

}
