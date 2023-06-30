package validate

import "testing"

func TestRoleAssignmentMarketplaceID(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{

		{
			Input: "",
			Valid: false,
		},

		{
			Input: "/",
			Valid: false,
		},

		{
			Input: "/providers/Microsoft.Marketplace/providers/Microsoft.Authorization/roleAssignments/",
			Valid: false,
		},

		{
			Input: "/providers/Microsoft.Subscription/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121",
			Valid: false,
		},

		{
			Input: "/providers/Microsoft.Marketplace/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121",
			Valid: true,
		},

		{
			Input: "/providers/Microsoft.Marketplace/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121|12345678-1234-5678-1234-567890123456",
			Valid: true,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := RoleAssignmentMarketplaceID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
