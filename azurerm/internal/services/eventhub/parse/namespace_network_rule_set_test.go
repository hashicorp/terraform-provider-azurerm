package parse

import "testing"

func TestEventHubNamespaceNetworkRuleID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Error    bool
		Expected *EventHubNamespaceNetworkRuleSetId
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
			Name:  "Missing Event Hub Namespace Name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.EventHub/namespaces/",
			Error: true,
		},
		{
			Name:  "Event Hub Namespace ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.EventHub/namespaces/namespace1",
			Error: true,
		},
		{
			Name:  "Missing Network Rule Name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.EventHub/namespaces/namespace1/networkrulesets/",
			Error: true,
		},
		{
			Name:  "Event Hub Namespace Network Rule ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.EventHub/namespaces/namespace1/networkrulesets/default",
			Expected: &EventHubNamespaceNetworkRuleSetId{
				Name:          "default",
				NamespaceName: "namespace1",
				ResourceGroup: "resGroup1",
			},
		},
		{
			Name:  "Wrong casing",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.EventHub/Namespaces/namespace1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := EventHubNamespaceNetworkRuleSetID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.NamespaceName != v.Expected.NamespaceName {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.NamespaceName, actual.NamespaceName)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
