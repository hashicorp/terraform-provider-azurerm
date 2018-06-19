package validate

import (
	"testing"
	"time"
)

func TestHelper_Validate_RFC3339Time(t *testing.T) {
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
			t.Fatalf("Expected Rfc3339Time to have an error for %q", tc.Time)
		}
	}
}

func TestHelper_Validate_Rfc3339DateInFutureBy(t *testing.T) {
	cases := []struct {
		Time     string
		Duration time.Duration
		Errors   int
	}{
		{
			Time:     "",
			Duration: time.Hour,
			Errors:   1,
		},
		{
			Time:     "not a time",
			Duration: time.Hour,
			Errors:   1,
		},
		{
			Time:     time.Now().String(),
			Duration: time.Hour,
			Errors:   1,
		},
		{
			Time:     time.Now().Add(time.Hour * 7).String(),
			Duration: time.Hour,
			Errors:   1,
		},
		{
			Time:     time.Now().Add(time.Minute).String(),
			Duration: time.Minute * 7,
			Errors:   0,
		},
	}

	for _, tc := range cases {
		_, errors := Rfc3339DateInFutureBy(tc.Duration)(tc.Time, "test")

		if len(errors) < tc.Errors {
			t.Fatalf("Expected Rfc3339DateInFutureBy to have an error for %q in future by %q", tc.Time, tc.Duration.String())
		}
	}
}
