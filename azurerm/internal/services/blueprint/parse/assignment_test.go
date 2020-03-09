package parse

import "testing"

func TestBlueprintAssignmentID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Error    bool
		Expected *BlueprintAssignmentId
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "Missing management group ID",
			Input: "/providers/Microsoft.Management/managementGroups/providers/Microsoft.Blueprint/blueprintAssignments/assignment1",
			Error: true,
		},
		{
			Name:  "Missing assignment name",
			Input: "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Blueprint/blueprintAssignments/",
			Error: true,
		},
		{
			Name:  "Blueprint Assignment ID in subscription",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Blueprint/blueprintAssignments/assignment1",
			Expected: &BlueprintAssignmentId{
				Name: "assignment1",
				BlueprintAssignmentScopeId: BlueprintAssignmentScopeId{
					ScopeId:        "/subscriptions/00000000-0000-0000-0000-000000000000",
					SubscriptionId: "00000000-0000-0000-0000-000000000000",
				},
			},
		},
		{
			Name:  "Missing assignment name in subscription",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Blueprint/blueprintAssignments/",
			Error: true,
		},
		{
			Name:  "Assignment ID in resource group",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Blueprint/blueprintAssignments/assignment1",
			Error: true,
		},
		{
			Name:  "missing resource group name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/providers/Microsoft.Blueprint/blueprintAssignments/assignment1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := BlueprintAssignmentID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %+v", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q", v.Expected.Name, actual.Name)
		}

		if actual.ScopeId != v.Expected.ScopeId {
			t.Fatalf("Expected %q but got %q", v.Expected.ScopeId, actual.ScopeId)
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q", v.Expected.SubscriptionId, actual.SubscriptionId)
		}
	}
}

func TestBlueprintAssignmentScopeID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Error    bool
		Expected *BlueprintAssignmentScopeId
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "Management group ID but missing components",
			Input: "/providers/Microsoft.Management/managementGroups/",
			Error: true,
		},
		{
			Name:  "Resource Group ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1",
			Error: true,
		},
		{
			Name:  "Incomplete resource group ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Error: true,
		},
		{
			Name:  "Resource ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/virtualMachines/vm1",
			Error: true,
		},
		{
			Name:  "Incomplete resource ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/virtualMachines/",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := BlueprintAssignmentScopeID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %+v", err)
		}

		if actual.ScopeId != v.Expected.ScopeId {
			t.Fatalf("Expected %q but got %q", v.Expected.ScopeId, actual.ScopeId)
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q", v.Expected.SubscriptionId, actual.SubscriptionId)
		}
	}
}
