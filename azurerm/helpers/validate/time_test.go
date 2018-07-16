package validate

import (
	"testing"
	"time"
)

func TestRFC3339Time(t *testing.T) {
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
		t.Run(tc.Time, func(t *testing.T) {
			_, errors := RFC3339Time(tc.Time, "test")

			if len(errors) != tc.Errors {
				t.Fatalf("Expected RFC3339Time to have %d not %d errors for %q", tc.Errors, len(errors), tc.Time)
			}
		})
	}
}

func TestRfc3339DateInFutureBy(t *testing.T) {
	cases := []struct {
		Name     string
		Time     string
		Duration time.Duration
		Errors   int
	}{
		{
			Name:     "empty",
			Time:     "",
			Duration: time.Hour,
			Errors:   1,
		},
		{
			Name:     "not a time",
			Time:     "not a time",
			Duration: time.Hour,
			Errors:   1,
		},
		{
			Name:     "now is not 1 hour ahead",
			Time:     time.Now().String(),
			Duration: time.Hour,
			Errors:   1,
		},
		{
			Name:     "now + 7 hours is not 1 hour ahead",
			Time:     time.Now().Add(time.Hour * 7).String(),
			Duration: time.Hour,
			Errors:   1,
		},
		{
			Name:     "now + 7 min is 7 min ahead",
			Time:     time.Now().Add(time.Minute).String(),
			Duration: time.Minute * 7,
			Errors:   0,
		},
		{
			Name:     "now + 8 min is at least 7 min ahead",
			Time:     time.Now().Add(time.Minute).String(),
			Duration: time.Minute * 7,
			Errors:   0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			_, errors := RFC3339DateInFutureBy(tc.Duration)(tc.Time, "test")

			if len(errors) < tc.Errors {
				t.Fatalf("Expected RFC3339DateInFutureBy to have %d not %d errors for %q in future by %q", tc.Errors, len(errors), tc.Time, tc.Duration.String())
			}
		})
	}
}
