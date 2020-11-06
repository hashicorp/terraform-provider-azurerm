package validate

import "testing"

func TestNetworkConnectionMonitorHttpPath(t *testing.T) {
	cases := []struct {
		Value  string
		Errors int
	}{
		{
			Value:  "",
			Errors: 1,
		},
		{
			Value:  "a/b",
			Errors: 1,
		},
		{
			Value:  "/ab/b1/",
			Errors: 0,
		},
		{
			Value:  "/a/b",
			Errors: 0,
		},
		{
			Value:  "http://www.terraform.io",
			Errors: 1,
		},
		{
			Value:  "/a/b/",
			Errors: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Value, func(t *testing.T) {
			_, errors := NetworkConnectionMonitorHttpPath(tc.Value, "path")

			if len(errors) != tc.Errors {
				t.Fatalf("Expected Path to return %d error(s) not %d", tc.Errors, len(errors))
			}
		})
	}
}

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

func TestNetworkConnectionMonitorEndpointAddressWithDomainName(t *testing.T) {
	cases := []struct {
		Value  string
		Errors int
	}{
		{
			Value:  "",
			Errors: 1,
		},
		{
			Value:  "a-b",
			Errors: 0,
		},
		{
			Value:  "terraform.io",
			Errors: 0,
		},
		{
			Value:  "www.google.com",
			Errors: 0,
		},
		{
			Value:  "http://www.google.com",
			Errors: 1,
		},
		{
			Value:  "www.google.com/a/b?a=1",
			Errors: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Value, func(t *testing.T) {
			_, errors := NetworkConnectionMonitorEndpointAddress(tc.Value, "address")

			if len(errors) != tc.Errors {
				t.Fatalf("Expected address to return %d error(s) not %d", tc.Errors, len(errors))
			}
		})
	}
}
