package validate

import "testing"

func TestRedisEnterpriseClusterLocation(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			// empty
			Input: "",
			Valid: false,
		},
		{
			// Unsupported location
			Input: "West US",
			Valid: false,
		},
		{
			// Unsupported location all upper with space
			Input: "WEST US",
			Valid: false,
		},
		{
			// Unsupported location all lower with space
			Input: "west us",
			Valid: false,
		},
		{
			// Unsupported location all upper without space
			Input: "WESTUS",
			Valid: false,
		},
		{
			// Unsupported location all lower without space
			Input: "westus",
			Valid: false,
		},
		{
			// Random text
			Input: "Lorem ipsum dolor sit amet",
			Valid: false,
		},
		{
			// Expected input
			Input: "Australia East",
			Valid: true,
		},
		{
			// No space
			Input: "AustraliaEast",
			Valid: true,
		},
		{
			// All Upper no space
			Input: "AUSTRALIAEAST",
			Valid: true,
		},
		{
			// All lower no space
			Input: "australiaeast",
			Valid: true,
		},
		{
			// All Upper with space
			Input: "AUSTRALIA EAST",
			Valid: true,
		},
		{
			// All lower with space
			Input: "australia east",
			Valid: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := RedisEnterpriseClusterLocation(tc.Input, "location")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
