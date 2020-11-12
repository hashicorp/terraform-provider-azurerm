package validate

import (
	"testing"
)

func TestVirtualHubRouteTableName(t *testing.T) {
	cases := []struct {
		Input       string
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "hello",
			ExpectError: false,
		},
		{
			Input:       "hello-world",
			ExpectError: false,
		},
		{
			Input:       "test<",
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		_, errors := VirtualHubRouteTableName(tc.Input, "name")

		hasError := len(errors) > 0
		if tc.ExpectError && !hasError {
			t.Fatalf("Expected the Virtual Hub Route Table Name to trigger a validation error for '%s'", tc.Input)
		}
	}
}
