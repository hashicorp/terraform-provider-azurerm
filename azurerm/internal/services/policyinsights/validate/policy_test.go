package validate

import "testing"

func TestValidatePolicyAssignmentID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected bool
	}{
		{
			Name:     "empty",
			Input:    "",
			Expected: false,
		},
		{
			Name:     "policy assignment in resource group",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Authorization/policyAssignments/assignment1",
			Expected: true,
		},
		{
			Name:     "policy assignment in resource group but no name",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Authorization/policyAssignments/",
			Expected: false,
		},
		{
			Name:     "the returned value of policy assignment id may not keep its casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.authorization/policyassignments/assignment1",
			Expected: true,
		},
		{
			Name:     "policy assignment in subscription",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policyAssignments/assignment1",
			Expected: true,
		},
		{
			Name:     "policy assignment in subscription but no name",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policyAssignments/",
			Expected: false,
		},
		{
			Name:     "policy assignment in management group",
			Input:    "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policyAssignments/assignment1",
			Expected: true,
		},
		{
			Name:     "policy assignment in management group but no name",
			Input:    "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policyAssignments/",
			Expected: false,
		},
		{
			Name:     "policy assignment in resource but no name",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Compute/virtualMachines/vm1/providers/Microsoft.Authorization/policyAssignments/assignment1",
			Expected: true,
		},
		{
			Name:     "policy assignment in resource but no name",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Compute/virtualMachines/vm1/providers/Microsoft.Authorization/policyAssignments/",
			Expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		_, errors := PolicyAssignmentID(v.Input, "name")
		actual := len(errors) == 0
		if v.Expected != actual {
			t.Fatalf("Expected %t but got %t", v.Expected, actual)
		}
	}
}

func TestValidatePolicyDefinitionID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected bool
	}{
		{
			Name:     "empty",
			Input:    "",
			Expected: false,
		},
		{
			Name:     "regular policy definition",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policyDefinitions/def1",
			Expected: true,
		},
		{
			Name:     "regular policy definition but no name",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policyDefinitions/",
			Expected: false,
		},
		{
			Name:     "policy definition in management group",
			Input:    "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policyDefinitions/def1",
			Expected: true,
		},
		{
			Name:     "policy definition in management group (without keep casing)",
			Input:    "/providers/Microsoft.Management/managementgroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policyDefinitions/def1",
			Expected: true,
		},
		{
			Name:     "policy definition in management group but no name",
			Input:    "/providers/Microsoft.Management/managementgroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policyDefinitions/",
			Expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		_, errors := PolicyDefinitionID(v.Input, "name")
		actual := len(errors) == 0
		if v.Expected != actual {
			t.Fatalf("Expected %t but got %t", v.Expected, actual)
		}
	}
}
