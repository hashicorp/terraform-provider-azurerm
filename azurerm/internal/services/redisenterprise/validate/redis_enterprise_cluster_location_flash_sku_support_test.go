package validate

import "testing"

func TestRedisEnterpriseClusterLocationFlashSkuSupport(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			// Invalid location
			Input: "UK West",
			Valid: false,
		},
		{
			// Unsupported location all upper with space
			Input: "UK WEST",
			Valid: false,
		},
		{
			// Unsupported location all lower with space
			Input: "uk west",
			Valid: false,
		},
		{
			// Unsupported location all upper without space
			Input: "UKWEST",
			Valid: false,
		},
		{
			// Unsupported location all lower without space
			Input: "ukwest",
			Valid: false,
		},
		{
			// empty
			Input: "",
			Valid: true,
		},
		{
			// Random text
			Input: "Lorem ipsum dolor sit amet",
			Valid: true,
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
		var valid bool
		if err := RedisEnterpriseClusterLocationFlashSkuSupport(tc.Input); err == nil {
			valid = true
		}

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
