package validate

import "testing"

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
