package validate

import "testing"

func TestTenantTemplateDeploymentID(t *testing.T) {
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
			// missing DeploymentName
			Input: "/providers/Microsoft.Resources/",
			Valid: false,
		},

		{
			// missing value for DeploymentName
			Input: "/providers/Microsoft.Resources/deployments/",
			Valid: false,
		},

		{
			// valid
			Input: "/providers/Microsoft.Resources/deployments/deploy1",
			Valid: true,
		},

		{
			// upper-cased
			Input: "/PROVIDERS/MICROSOFT.RESOURCES/DEPLOYMENTS/DEPLOY1",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := TenantTemplateDeploymentID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
