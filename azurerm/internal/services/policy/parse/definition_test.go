package parse

import (
	"reflect"
	"testing"
)

func TestValidatePolicyDefinitionID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Error    bool
		Expected *PolicyDefinitionId
	}{
		{
			Name:  "empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "built-in policy definition ID",
			Input: "/providers/Microsoft.Authorization/policyDefinitions/00000000-0000-0000-0000-000000000000",
			Expected: &PolicyDefinitionId{
				Name: "00000000-0000-0000-0000-000000000000",
			},
		},
		{
			Name:  "regular policy definition",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policyDefinitions/def1",
			Expected: &PolicyDefinitionId{
				Name: "def1",
				PolicyScopeId: ScopeAtSubscription{
					scopeId:        "/subscriptions/00000000-0000-0000-0000-000000000000",
					SubscriptionId: "00000000-0000-0000-0000-000000000000",
				},
			},
		},
		{
			Name:  "regular policy definition but no name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policyDefinitions/",
			Error: true,
		},
		{
			Name:  "policy definition in management group",
			Input: "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policyDefinitions/def1",
			Expected: &PolicyDefinitionId{
				Name: "def1",
				PolicyScopeId: ScopeAtManagementGroup{
					scopeId:             "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000",
					ManagementGroupName: "00000000-0000-0000-0000-000000000000",
				},
			},
		},
		{
			Name:  "policy definition in management group with wrong casing",
			Input: "/providers/Microsoft.Management/managementgroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policyDefinitions/def1",
			Expected: &PolicyDefinitionId{
				Name: "def1",
				PolicyScopeId: ScopeAtManagementGroup{
					scopeId:             "/providers/Microsoft.Management/managementgroups/00000000-0000-0000-0000-000000000000",
					ManagementGroupName: "00000000-0000-0000-0000-000000000000",
				},
			},
		},
		{
			Name:  "policy definition in management group but no name",
			Input: "/providers/Microsoft.Management/managementgroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policyDefinitions/",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := PolicyDefinitionID(v.Input)
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
	}
}
