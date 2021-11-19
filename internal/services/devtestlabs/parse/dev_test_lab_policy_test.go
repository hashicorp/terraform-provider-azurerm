package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceid"
)

var _ resourceid.Formatter = DevTestLabPolicyId{}

func TestDevTestLabPolicyIDFormatter(t *testing.T) {
	actual := NewDevTestLabPolicyID("12345678-1234-9876-4563-123456789012", "group1", "lab1", "policyset1", "policy1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.DevTestLab/labs/lab1/policySets/policyset1/policies/policy1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestDevTestLabPolicyID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *DevTestLabPolicyId
	}{

		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing SubscriptionId
			Input: "/",
			Error: true,
		},

		{
			// missing value for SubscriptionId
			Input: "/subscriptions/",
			Error: true,
		},

		{
			// missing ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Error: true,
		},

		{
			// missing LabName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.DevTestLab/",
			Error: true,
		},

		{
			// missing value for LabName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.DevTestLab/labs/",
			Error: true,
		},

		{
			// missing PolicySetName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.DevTestLab/labs/lab1/",
			Error: true,
		},

		{
			// missing value for PolicySetName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.DevTestLab/labs/lab1/policySets/",
			Error: true,
		},

		{
			// missing PolicyName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.DevTestLab/labs/lab1/policySets/policyset1/",
			Error: true,
		},

		{
			// missing value for PolicyName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.DevTestLab/labs/lab1/policySets/policyset1/policies/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.DevTestLab/labs/lab1/policySets/policyset1/policies/policy1",
			Expected: &DevTestLabPolicyId{
				SubscriptionId: "12345678-1234-9876-4563-123456789012",
				ResourceGroup:  "group1",
				LabName:        "lab1",
				PolicySetName:  "policyset1",
				PolicyName:     "policy1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/GROUP1/PROVIDERS/MICROSOFT.DEVTESTLAB/LABS/LAB1/POLICYSETS/POLICYSET1/POLICIES/POLICY1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := DevTestLabPolicyID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q for SubscriptionId", v.Expected.SubscriptionId, actual.SubscriptionId)
		}
		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
		if actual.LabName != v.Expected.LabName {
			t.Fatalf("Expected %q but got %q for LabName", v.Expected.LabName, actual.LabName)
		}
		if actual.PolicySetName != v.Expected.PolicySetName {
			t.Fatalf("Expected %q but got %q for PolicySetName", v.Expected.PolicySetName, actual.PolicySetName)
		}
		if actual.PolicyName != v.Expected.PolicyName {
			t.Fatalf("Expected %q but got %q for PolicyName", v.Expected.PolicyName, actual.PolicyName)
		}
	}
}
