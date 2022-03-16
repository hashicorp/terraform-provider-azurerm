package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = ManagementGroupPolicyRemediationId{}

func TestManagementGroupPolicyRemediationIDFormatter(t *testing.T) {
	actual := NewManagementGroupPolicyRemediationID("group1", "remediation1").ID()
	expected := "/providers/namespace1/managementGroups/group1/providers/Microsoft.PolicyInsights/remediations/remediation1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestManagementGroupPolicyRemediationID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ManagementGroupPolicyRemediationId
	}{

		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing ManagementGroupName
			Input: "/providers/namespace1/",
			Error: true,
		},

		{
			// missing value for ManagementGroupName
			Input: "/providers/namespace1/managementGroups/",
			Error: true,
		},

		{
			// missing RemediationName
			Input: "/providers/namespace1/managementGroups/group1/providers/Microsoft.PolicyInsights/",
			Error: true,
		},

		{
			// missing value for RemediationName
			Input: "/providers/namespace1/managementGroups/group1/providers/Microsoft.PolicyInsights/remediations/",
			Error: true,
		},

		{
			// valid
			Input: "/providers/namespace1/managementGroups/group1/providers/Microsoft.PolicyInsights/remediations/remediation1",
			Expected: &ManagementGroupPolicyRemediationId{
				ManagementGroupName: "group1",
				RemediationName:     "remediation1",
			},
		},

		{
			// upper-cased
			Input: "/PROVIDERS/NAMESPACE1/MANAGEMENTGROUPS/GROUP1/PROVIDERS/MICROSOFT.POLICYINSIGHTS/REMEDIATIONS/REMEDIATION1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ManagementGroupPolicyRemediationID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.ManagementGroupName != v.Expected.ManagementGroupName {
			t.Fatalf("Expected %q but got %q for ManagementGroupName", v.Expected.ManagementGroupName, actual.ManagementGroupName)
		}
		if actual.RemediationName != v.Expected.RemediationName {
			t.Fatalf("Expected %q but got %q for RemediationName", v.Expected.RemediationName, actual.RemediationName)
		}
	}
}
