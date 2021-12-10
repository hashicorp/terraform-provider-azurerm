package validate

import "testing"

func TestNetworkConnectionMonitorValidStatusCodeRanges(t *testing.T) {
	cases := []struct {
		Value  string
		Errors int
	}{
		{
			Value:  "",
			Errors: 1,
		},
		{
			Value:  "100",
			Errors: 0,
		},
		{
			Value:  "599",
			Errors: 0,
		},
		{
			Value:  "600",
			Errors: 0,
		},
		{
			Value:  "1xx",
			Errors: 0,
		},
		{
			Value:  "10x",
			Errors: 0,
		},
		{
			Value:  "1x2",
			Errors: 0,
		},
		{
			Value:  "259-379",
			Errors: 0,
		},
		{
			Value:  "489-379",
			Errors: 1,
		},
		{
			Value:  "30x-4xx",
			Errors: 1,
		},
		{
			Value:  "7xx",
			Errors: 1,
		},
		{
			Value:  "99",
			Errors: 1,
		},
		{
			Value:  "888",
			Errors: 1,
		},
		{
			Value:  "1111",
			Errors: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Value, func(t *testing.T) {
			_, errors := NetworkConnectionMonitorValidStatusCodeRanges(tc.Value, "valid_status_code_ranges")

			if len(errors) != tc.Errors {
				t.Fatalf("Expected valid_status_code_ranges to return %d error(s) not %d", tc.Errors, len(errors))
			}
		})
	}
}
