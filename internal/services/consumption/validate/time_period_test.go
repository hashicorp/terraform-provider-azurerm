package validate

import (
	"testing"
	"time"
)

func TestConsumptionBudgetTimePeriodStartDate(t *testing.T) {
	// Set up time for testing
	now := time.Now()
	validTime := time.Date(
		now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		Input         string
		ExpectError   bool
		ExpectWarning bool
	}{
		{
			Input:         "",
			ExpectError:   true,
			ExpectWarning: false,
		},
		{
			Input:         "2006-01-02",
			ExpectError:   true,
			ExpectWarning: false,
		},
		{
			// Not on the first of a month
			Input:         "2020-11-02T00:00:00Z",
			ExpectError:   true,
			ExpectWarning: false,
		},
		{
			// Before June 1, 2017
			Input:         "2000-01-01T00:00:00Z",
			ExpectError:   true,
			ExpectWarning: false,
		},
		{
			// Valid date and time
			Input:         validTime.Format(time.RFC3339),
			ExpectError:   false,
			ExpectWarning: false,
		},
		{
			// More than 12 months in the future
			Input:         validTime.AddDate(2, 0, 0).Format(time.RFC3339),
			ExpectError:   false,
			ExpectWarning: true,
		},
	}

	for _, tc := range cases {
		warnings, errors := ConsumptionBudgetTimePeriodStartDate(tc.Input, "start_date")
		if errors != nil {
			if !tc.ExpectError {
				t.Fatalf("Got error for input %q: %+v", tc.Input, errors)
			}

			return
		}

		if warnings != nil {
			if !tc.ExpectWarning {
				t.Fatalf("Got warnings for input %q: %+v", tc.Input, warnings)
			}

			return
		}

		if tc.ExpectError && len(errors) == 0 {
			t.Fatalf("Got no errors for input %q but expected some", tc.Input)
		} else if !tc.ExpectError && len(errors) > 0 {
			t.Fatalf("Got %d errors for input %q when didn't expect any", len(errors), tc.Input)
		}

		if tc.ExpectWarning && len(warnings) == 0 {
			t.Fatalf("Got no warnings for input %q but expected some", tc.Input)
		} else if !tc.ExpectWarning && len(warnings) > 0 {
			t.Fatalf("Got %d warnings for input %q when didn't expect any", len(warnings), tc.Input)
		}
	}
}
