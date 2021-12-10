package parse

import (
	"reflect"
	"testing"
)

func TestPolicySetDefinitionID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Error    bool
		Expected *PolicySetDefinitionId
	}{
		{
			Name:  "empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "regular policy set definition",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policySetDefinitions/set1",
			Expected: &PolicySetDefinitionId{
				Name: "set1",
				PolicyScopeId: ScopeAtSubscription{
					scopeId:        "/subscriptions/00000000-0000-0000-0000-000000000000",
					SubscriptionId: "00000000-0000-0000-0000-000000000000",
				},
			},
		},
		{
			Name:  "regular policy set definition but no name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policySetDefinitions/",
			Error: true,
		},
		{
			Name:  "policy set definition in management group",
			Input: "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policySetDefinitions/set1",
			Expected: &PolicySetDefinitionId{
				Name: "set1",
				PolicyScopeId: ScopeAtManagementGroup{
					scopeId:             "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000",
					ManagementGroupName: "00000000-0000-0000-0000-000000000000",
				},
			},
		},
		{
			Name:  "policy set definition in management group with inconsistent casing",
			Input: "/providers/Microsoft.Management/managementgroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policySetDefinitions/set1",
			Expected: &PolicySetDefinitionId{
				Name: "set1",
				PolicyScopeId: ScopeAtManagementGroup{
					scopeId:             "/providers/Microsoft.Management/managementgroups/00000000-0000-0000-0000-000000000000",
					ManagementGroupName: "00000000-0000-0000-0000-000000000000",
				},
			},
		},
		{
			Name:  "policy set definition in management group but no name",
			Input: "/providers/Microsoft.Management/managementgroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policySetDefinitions/",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := PolicySetDefinitionID(v.Input)
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
