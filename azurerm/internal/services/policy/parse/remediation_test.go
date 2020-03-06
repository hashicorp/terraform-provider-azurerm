package parse

import (
	"reflect"
	"testing"
)

func TestRemediationID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Error    bool
		Expected *RemediationId
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "Subscription ID or Management Group ID itself",
			Input: "00000000-0000-0000-0000-000000000000",
			Error: true,
		},
		{
			Name:  "Policy Remediation ID at Subscription",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Policyinsights/remediations/test",
			Expected: &RemediationId{
				Name: "test",
				PolicyScopeId: ScopeAtSubscription{
					scopeId:        "/subscriptions/00000000-0000-0000-0000-000000000000",
					SubscriptionId: "00000000-0000-0000-0000-000000000000",
				},
			},
		},
		{
			Name:  "Policy Remediation ID at Subscription with wrong casing",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/microsoft.policyinsights/remediations/test",
			Expected: &RemediationId{
				Name: "test",
				PolicyScopeId: ScopeAtSubscription{
					scopeId:        "/subscriptions/00000000-0000-0000-0000-000000000000",
					SubscriptionId: "00000000-0000-0000-0000-000000000000",
				},
			},
		},
		{
			Name:  "Policy Remediation ID at Subscription but missing name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Policyinsights/remediations/",
			Error: true,
		},
		{
			Name:  "No resource group name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/providers/Microsoft.Policyinsights/remediations/test",
			Error: true,
		},
		{
			Name:  "Policy Remediation ID at Resource Group",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Policyinsights/remediations/test",
			Expected: &RemediationId{
				Name: "test",
				PolicyScopeId: ScopeAtResourceGroup{
					scopeId:        "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1",
					SubscriptionId: "00000000-0000-0000-0000-000000000000",
					ResourceGroup:  "group1",
				},
			},
		},
		{
			Name:  "Policy Remediation ID at Resource Group with wrong casing",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/group1/providers/microsoft.policyinsights/remediations/test",
			Expected: &RemediationId{
				Name: "test",
				PolicyScopeId: ScopeAtResourceGroup{
					scopeId:        "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/group1",
					SubscriptionId: "00000000-0000-0000-0000-000000000000",
					ResourceGroup:  "group1",
				},
			},
		},
		{
			Name:  "Policy Remediation ID at Resource Group but missing name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Policyinsights/remediations/",
			Error: true,
		},
		{
			Name:  "Missing scope resource name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/virtualMachines/providers/Microsoft.Policyinsights/remediations/test",
			Error: true,
		},
		{
			Name:  "Policy Remediation ID at Resource",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/virtualMachines/vm1/providers/Microsoft.Policyinsights/remediations/test",
			Expected: &RemediationId{
				Name: "test",
				PolicyScopeId: ScopeAtResource{
					scopeId: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/virtualMachines/vm1",
				},
			},
		},
		{
			Name:  "Policy Remediation ID at Resource with wrong casing",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/group1/providers/microsoft.compute/virtualmachines/vm1/providers/microsoft.policyinsights/remediations/test",
			Expected: &RemediationId{
				Name: "test",
				PolicyScopeId: ScopeAtResource{
					scopeId: "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/group1/providers/microsoft.compute/virtualmachines/vm1",
				},
			},
		},
		{
			Name:  "Policy Remediation ID at Resource but missing name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/virtualMachines/vm1/providers/Microsoft.Policyinsights/remediations/",
			Error: true,
		},
		{
			Name:  "Policy Remediation ID at Management Group",
			Input: "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.PolicyInsights/remediations/test",
			Expected: &RemediationId{
				Name: "test",
				PolicyScopeId: ScopeAtManagementGroup{
					scopeId:           "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000",
					ManagementGroupId: "00000000-0000-0000-0000-000000000000",
				},
			},
		},
		{
			Name:  "Policy Remediation ID at Management Group with readable id",
			Input: "/providers/Microsoft.Management/managementGroups/group1/providers/Microsoft.PolicyInsights/remediations/test",
			Expected: &RemediationId{
				Name: "test",
				PolicyScopeId: ScopeAtManagementGroup{
					scopeId:           "/providers/Microsoft.Management/managementGroups/group1",
					ManagementGroupId: "group1",
				},
			},
		},
		{
			Name:  "Policy Remediation ID at Management Group with wrong casing",
			Input: "/providers/microsoft.management/managementgroups/00000000-0000-0000-0000-000000000000/providers/microsoft.policyinsights/remediations/test",
			Expected: &RemediationId{
				Name: "test",
				PolicyScopeId: ScopeAtManagementGroup{
					scopeId:           "/providers/microsoft.management/managementgroups/00000000-0000-0000-0000-000000000000",
					ManagementGroupId: "00000000-0000-0000-0000-000000000000",
				},
			},
		},
		{
			Name:  "Policy Remediation ID at Management Group but missing name",
			Input: "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.PolicyInsights/remediations/",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := RemediationID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %+v", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q", v.Expected.Name, actual.Name)
		}

		if !reflect.DeepEqual(v.Expected.PolicyScopeId, actual.PolicyScopeId) {
			t.Fatalf("Expected %+v but got %+v", v.Expected.PolicyScopeId, actual.PolicyScopeId)
		}

		//if actual.Type != v.Expected.Type {
		//	t.Fatalf("Expected type %q but got type %q", v.Expected.Type, actual.Type)
		//}
		//
		//if actual.ScopeId != v.Expected.ScopeId {
		//	t.Fatalf("Expected %q but got %q", v.Expected.ScopeId, actual.ScopeId)
		//}
		//
		//if actual.SubscriptionId != v.Expected.SubscriptionId {
		//	t.Fatalf("Expected %q but got %q", v.Expected.SubscriptionId, actual.SubscriptionId)
		//}
		//
		//if actual.ResourceGroup != v.Expected.ResourceGroup {
		//	t.Fatalf("Expected %q but got %q", v.Expected.ResourceGroup, actual.ResourceGroup)
		//}
		//
		//if actual.ManagementGroupId != v.Expected.ManagementGroupId {
		//	t.Fatalf("Expected %q but got %q", v.Expected.ManagementGroupId, actual.ManagementGroupId)
		//}
	}
}
