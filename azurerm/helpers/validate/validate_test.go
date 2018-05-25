package validate

import "testing"

func TestInternal_ValidateRFC3339Time(t *testing.T) {
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
			Time:   "2000-01-01",
			Errors: 1,
		},
		{
			Time:   "2000-01-01T01:23:45",
			Errors: 1,
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
		_, errors := Rfc3339Time(tc.Time, "test")

		if len(errors) != tc.Errors {
			t.Fatalf("Expected Rfc3339Time to have an error for '%q'", tc.Time)
		}
	}
}
