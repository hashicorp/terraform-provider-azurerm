package parse

import (
	"testing"
)

func TestNamespaceAuthorizationRuleID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *NamespaceAuthorizationRuleId
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
			Name:  "No Resource Groups Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Error: true,
		},
		{
			Name:  "Resource Group ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/",
			Error: true,
		},
		{
			Name:  "Missing Namespaces Key",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/",
			Error: true,
		},
		{
			Name:  "Missing Namespaces Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1",
			Error: true,
		},
		{
			Name:  "Missing authorizationRules Key",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1/authorizationRules",
			Error: true,
		},
		{
			Name:  "Namespace Authorization Rule ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1/authorizationRules/rule1",
			Error: false,
			Expect: &NamespaceAuthorizationRuleId{
				ResourceGroup:         "group1",
				NamespaceName:         "namespace1",
				AuthorizationRuleName: "rule1",
			},
		},
		{
			Name:  "Wrong Casing",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1/AuthorizationRules/rule1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := NamespaceAuthorizationRuleID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.AuthorizationRuleName != v.Expect.AuthorizationRuleName {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.AuthorizationRuleName, actual.AuthorizationRuleName)
		}

		if actual.NamespaceName != v.Expect.NamespaceName {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.NamespaceName, actual.NamespaceName)
		}

		if actual.ResourceGroup != v.Expect.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expect.ResourceGroup, actual.ResourceGroup)
		}
	}
}
