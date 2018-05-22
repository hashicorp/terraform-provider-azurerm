package azurerm

import (
	"strconv"
	"testing"
)

func TestValidateRFC3339Date(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "",
			ErrCount: 1,
		},
		{
			Value:    "Random",
			ErrCount: 1,
		},
		{
			Value:    "2017-01-01",
			ErrCount: 1,
		},
		{
			Value:    "2017-01-01T01:23:45",
			ErrCount: 1,
		},
		{
			Value:    "2017-01-01T01:23:45+00:00",
			ErrCount: 0,
		},
		{
			Value:    "2017-01-01T01:23:45Z",
			ErrCount: 0,
		},
	}

	for _, tc := range cases {
		_, errors := validateRFC3339Date(tc.Value, "example")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected validateRFC3339Date to trigger '%d' errors for '%s' - got '%d'", tc.ErrCount, tc.Value, len(errors))
		}
	}
}

func TestValidateIntInSlice(t *testing.T) {

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
		_, errors := validateIntInSlice(tc.Input)(tc.Value, "azurerm_postgresql_database")

		if len(errors) != tc.Errors {
			t.Fatalf("Expected the validateIntInSlice trigger a validation error for input: %+v looking for %+v", tc.Input, tc.Value)
		}
	}

}

func TestValidateIso8601Duration(t *testing.T) {
	cases := []struct {
		Value  string
		Errors int
	}{
		{
			// Date components only
			Value:  "P1Y2M3D",
			Errors: 0,
		},
		{
			// Time components only
			Value:  "PT7H42M3S",
			Errors: 0,
		},
		{
			// Date and time components
			Value:  "P1Y2M3DT7H42M3S",
			Errors: 0,
		},
		{
			// Invalid prefix
			Value:  "1Y2M3DT7H42M3S",
			Errors: 1,
		},
		{
			// Wrong order of components, i.e. invalid format
			Value:  "PT7H42M3S1Y2M3D",
			Errors: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateIso8601Duration()(tc.Value, "example")

		if len(errors) != tc.Errors {
			t.Fatalf("Expected validateIso8601Duration to trigger '%d' errors for '%s' - got '%d'", tc.Errors, tc.Value, len(errors))
		}
	}
}

func TestValidateIntBetweenDivisibleBy(t *testing.T) {
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
		_, errors := validateIntBetweenDivisibleBy(tc.Min, tc.Max, tc.Div)(tc.Value, strconv.Itoa(tc.Value.(int)))
		if len(errors) != tc.Errors {
			t.Fatalf("Expected intBetweenDivisibleBy to trigger '%d' errors for '%s' - got '%d'", tc.Errors, tc.Value, len(errors))
		}
	}
}

func TestValidateCollation(t *testing.T) {
	cases := []struct {
		Value  string
		Errors int
	}{
		{
			Value:  "en-US",
			Errors: 1,
		},
		{
			Value:  "en_US",
			Errors: 0,
		},
		{
			Value:  "en US",
			Errors: 0,
		},
		{
			Value:  "English_United States.1252",
			Errors: 0,
		},
	}

	for _, tc := range cases {
		_, errors := validateCollation()(tc.Value, "collation")
		if len(errors) != tc.Errors {
			t.Fatalf("Expected validateCollation to trigger '%d' errors for '%s' - got '%d'", tc.Errors, tc.Value, len(errors))
		}
	}
}
