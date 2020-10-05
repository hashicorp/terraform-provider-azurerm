package parse

import (
	"reflect"
	"testing"
)

func TestPolicyScopeID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Error    bool
		Expected PolicyScopeId
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
			Name:  "Subscription Resource ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000",
			Expected: ScopeAtSubscription{
				scopeId:        "/subscriptions/00000000-0000-0000-0000-000000000000",
				SubscriptionId: "00000000-0000-0000-0000-000000000000",
			},
		},
		{
			Name:  "Resource group name",
			Input: "group1",
			Error: true,
		},
		{
			Name:  "No resource group name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Error: true,
		},
		{
			Name:  "Resource Group ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1",
			Expected: ScopeAtResourceGroup{
				scopeId:        "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1",
				SubscriptionId: "00000000-0000-0000-0000-000000000000",
				ResourceGroup:  "group1",
			},
		},
		{
			Name:  "Incomplete resource ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/virtualMachines",
			Error: true,
		},
		{
			Name:  "Resource ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/virtualMachines/vm1",
			Expected: ScopeAtResource{
				scopeId: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/virtualMachines/vm1",
			},
		},
		{
			Name:  "Management Group ID",
			Input: "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000",
			Expected: ScopeAtManagementGroup{
				scopeId:             "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000",
				ManagementGroupName: "00000000-0000-0000-0000-000000000000",
			},
		},
		{
			Name:  "Management Group ID with readable id",
			Input: "/providers/Microsoft.Management/managementGroups/group1",
			Expected: ScopeAtManagementGroup{
				scopeId:             "/providers/Microsoft.Management/managementGroups/group1",
				ManagementGroupName: "group1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := PolicyScopeID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %+v", err)
		}

		if !reflect.DeepEqual(v.Expected, actual) {
			t.Fatalf("Expected %+v but got %+v", v.Expected, actual)
		}
	}
}
