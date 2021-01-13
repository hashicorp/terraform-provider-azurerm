package parse

import (
	"testing"
)

func TestAssignmentID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Error    bool
		Expected *AssignmentId
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "No Resource Groups Segment",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000",
			Error: true,
		},
		{
			Name:  "Invalid scope",
			Input: "/providers/Microsoft.Management/managementGroups/testAccManagementGroup",
			Error: true,
		},
		// We have two valid possibilities to check for
		{
			Name:  "Valid subscription scoped",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Blueprint/blueprintAssignments/assignSimpleBlueprint",
			Error: false,
			Expected: &AssignmentId{
				Scope:        "subscriptions/00000000-0000-0000-0000-000000000000",
				Subscription: "00000000-0000-0000-0000-000000000000",
				Name:         "assignSimpleBlueprint",
			},
		},
		{
			Name:  "Valid managementGroup scoped",
			Input: "/managementGroups/testAccManagementGroup/providers/Microsoft.Blueprint/blueprintAssignments/assignSimpleBlueprint",
			Error: false,
			Expected: &AssignmentId{
				Scope:           "managementGroups/testAccManagementGroup",
				ManagementGroup: "testAccManagementGroup",
				Name:            "assignSimpleBlueprint",
			},
		},
		{
			Name:  "wrong case - subscription scoped",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Blueprint/blueprintassignments/assignSimpleBlueprint",
			Error: true,
		},
		{
			Name:  "wrong case - managementGroup scoped",
			Input: "/managementGroups/testAccManagementGroup/providers/Microsoft.Blueprint/blueprintassignments/assignSimpleBlueprint",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := AssignmentID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Scope != v.Expected.Scope {
			t.Fatalf("Expected %q but got %q for Scope", v.Expected.Scope, actual.Scope)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.ManagementGroup == "" && actual.Subscription != v.Expected.Subscription {
			t.Fatalf("Expected %q but got %q for Subscription", v.Expected.Subscription, actual.Subscription)
		}

		if actual.Subscription == "" && actual.ManagementGroup != v.Expected.ManagementGroup {
			t.Fatalf("Expected %q but got %q for ManagementGroup", v.Expected.ManagementGroup, actual.ManagementGroup)
		}
	}
}
