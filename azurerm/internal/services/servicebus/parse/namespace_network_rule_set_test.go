package parse

import "testing"

func TestServiceBusNamespaceNetworkRuleID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Error    bool
		Expected *NamespaceNetworkRuleSetId
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
			Name:  "Missing Service Bus Namespace Name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.ServiceBus/namespaces/",
			Error: true,
		},
		{
			Name:  "Service Bus Namespace ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.ServiceBus/namespaces/namespace1",
			Error: true,
		},
		{
			Name:  "Missing Network Rule Name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.ServiceBus/namespaces/namespace1/networkrulesets/",
			Error: true,
		},
		{
			Name:  "Service Bus Namespace Network Rule ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.ServiceBus/namespaces/namespace1/networkrulesets/default",
			Expected: &NamespaceNetworkRuleSetId{
				NetworkrulesetName: "default",
				NamespaceName:      "namespace1",
				ResourceGroup:      "resGroup1",
			},
		},
		{
			Name:  "Wrong casing",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.ServiceBus/Namespaces/namespace1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := NamespaceNetworkRuleSetID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.NetworkrulesetName != v.Expected.NetworkrulesetName {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.NetworkrulesetName, actual.NetworkrulesetName)
		}

		if actual.NamespaceName != v.Expected.NamespaceName {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.NamespaceName, actual.NamespaceName)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
