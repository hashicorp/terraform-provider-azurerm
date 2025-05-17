package validate

import (
	"testing"
)

func TestNumberOfIpAddresses(t *testing.T) {
	cases := []struct {
		Input       string
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "a",
			ExpectError: true,
		},
		{
			Input:       "-1",
			ExpectError: true,
		},
		{
			Input:       "0",
			ExpectError: true,
		},
		{
			Input:       "1",
			ExpectError: false,
		},
		{
			Input:       "2",
			ExpectError: false,
		},
		{
			Input:       "123",
			ExpectError: true,
		},
		{
			Input:       "1024",
			ExpectError: false,
		},
		{
			Input:       "170141183460469231731687303715884105728", // 2^127
			ExpectError: false,
		},
		{
			Input:       "340282366920938463463374607431768211456", // 2^128
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		_, errors := NumberOfIpAddresses(tc.Input, "number_of_ip_addresses")

		hasError := len(errors) > 0
		if tc.ExpectError && !hasError {
			t.Fatalf("Expected NumberOfIpAddresses to trigger a validation error for '%s'", tc.Input)
		}
	}
}
