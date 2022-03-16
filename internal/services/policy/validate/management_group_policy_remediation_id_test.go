package validate

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import "testing"

func TestManagementGroupPolicyRemediationID(t *testing.T) {
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
			// missing ManagementGroupName
			Input: "/providers/namespace1/",
			Valid: false,
		},

		{
			// missing value for ManagementGroupName
			Input: "/providers/namespace1/managementGroups/",
			Valid: false,
		},

		{
			// missing RemediationName
			Input: "/providers/namespace1/managementGroups/group1/providers/Microsoft.PolicyInsights/",
			Valid: false,
		},

		{
			// missing value for RemediationName
			Input: "/providers/namespace1/managementGroups/group1/providers/Microsoft.PolicyInsights/remediations/",
			Valid: false,
		},

		{
			// valid
			Input: "/providers/namespace1/managementGroups/group1/providers/Microsoft.PolicyInsights/remediations/remediation1",
			Valid: true,
		},

		{
			// upper-cased
			Input: "/PROVIDERS/NAMESPACE1/MANAGEMENTGROUPS/GROUP1/PROVIDERS/MICROSOFT.POLICYINSIGHTS/REMEDIATIONS/REMEDIATION1",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := ManagementGroupPolicyRemediationID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
