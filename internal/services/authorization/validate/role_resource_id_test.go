package validate

import "testing"

func TestRoleResourceID(t *testing.T) {
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
			Input: "/providers/Microsoft.Authorization/roleDefinitions/23456781-2349-8764-5631-234567890121",
			Valid: true,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := RoleResourceID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
