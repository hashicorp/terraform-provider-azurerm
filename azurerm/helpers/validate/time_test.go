package validate

import (
	"testing"
)

func TestISO8601DateTime(t *testing.T) {
	cases := []struct {
		Time   string
		Errors int
	}{
		{
			Time:   "",
			Errors: 1,
		},
		{
			Time:   "this is not a date",
			Errors: 1,
		},
		{
			Time:   "2000-06-31", // No 31st of 6th
			Errors: 1,
		},
		{
			Time:   "01/21/2015", // not valid US date with slashes
			Errors: 1,
		},
		{
			Time:   "01-21-2015", // not valid US date with dashes
			Errors: 1,
		},
		{
			Time:   "2000-01-01",
			Errors: 0,
		},
		{
			Time:   "2000-01-01T01:23:45",
			Errors: 0,
		},
		{
			Time:   "2000-01-01T01:23:45Z",
			Errors: 0,
		},
		{
			Time:   "2000-01-01T01:23:45+00:00",
			Errors: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Time, func(t *testing.T) {
			_, errors := ISO8601DateTime(tc.Time, "test")

			if len(errors) != tc.Errors {
				t.Fatalf("Expected ISO8601DateTime to have %d but got %d errors for %q", tc.Errors, len(errors), tc.Time)
			}
		})
	}
}

func TestISO8601Duration(t *testing.T) {
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
		_, errors := ISO8601Duration(tc.Value, "example")

		if len(errors) != tc.Errors {
			t.Fatalf("Expected ISO8601Duration to trigger '%d' errors for '%s' - got '%d'", tc.Errors, tc.Value, len(errors))
		}
	}
}
