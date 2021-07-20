package validate

import "testing"

func TestPortalTenantConfigurationID(t *testing.T) {
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
			// missing Name
			Input: "/providers/Microsoft.Portal/",
			Valid: false,
		},

		{
			// missing value for Name
			Input: "/providers/Microsoft.Portal/tenantConfigurations/",
			Valid: false,
		},

		{
			// valid
			Input: "/providers/Microsoft.Portal/tenantConfigurations/default",
			Valid: true,
		},

		{
			// upper-cased
			Input: "/PROVIDERS/MICROSOFT.PORTAL/TENANTCONFIGURATIONS/DEFAULT",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := PortalTenantConfigurationID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
