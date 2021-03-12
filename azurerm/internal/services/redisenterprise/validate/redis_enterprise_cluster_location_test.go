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
			Input: "East US 3",
			Valid: false,
		},
		{
			// Unsupported location all upper with space
			Input: "EAST US 3",
			Valid: false,
		},
		{
			// Unsupported location all lower with space
			Input: "east us 3",
			Valid: false,
		},
		{
			// Unsupported location all upper without space
			Input: "EASTUS3",
			Valid: false,
		},
		{
			// Unsupported location all lower without space
			Input: "eastus3",
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
