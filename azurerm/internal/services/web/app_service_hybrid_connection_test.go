package web

import "testing"

func TestParseAppServiceHybridConnectionID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *HybridConnectionId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Expected: nil,
		},
		{
			Name:     "Missing Sites Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testResourceGroup1/providers/Microsoft.Web/sites/",
			Expected: nil,
		},
		{
			Name:     "Missing Namespace Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-00000000000/resourceGroups/testResourceGroup1/providers/Microsoft.Web/sites/testApp1/hybridConnectionNamespaces/",
			Expected: nil,
		},
		{
			Name:     "Missing Relay Name value",
			Input:    "/subscriptions/00000000-0000-0000-0000-00000000000/resourceGroups/testResourceGroup1/providers/Microsoft.Web/sites/testApp1/hybridConnectionNamespaces/testNamespace/relays/",
			Expected: nil,
		},
		{
			Name:     "Incorrect casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-00000000000/resourceGroups/testResourceGroup1/providers/Microsoft.Web/sites/testApp1/hybridConnectionNamespaces/testNamespace1/Relays/testRelay1",
			Expected: nil,
		},

		{
			Name:  "App Service Hybrid Connection Resource ID",
			Input: "/subscriptions/00000000-0000-0000-0000-00000000000/resourceGroups/testResourceGroup1/providers/Microsoft.Web/sites/testApp1/hybridConnectionNamespaces/testNamespace1/relays/testRelay1",
			Expected: &HybridConnectionId{
				ResourceGroup:                 "testResourceGroup1",
				RelayName:                     "testRelay1",
				SiteName:                      "testApp1",
				HybridConnectionNamespaceName: "testNamespace1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := HybridConnectionID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.RelayName != v.Expected.RelayName {
			t.Fatalf("Expected %s but got %s for Name", v.Expected.RelayName, actual.RelayName)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %s but got %s for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}

		if actual.SiteName != v.Expected.SiteName {
			t.Fatalf("Expected %s but got %s for SiteName", v.Expected.SiteName, actual.SiteName)
		}

		if actual.HybridConnectionNamespaceName != v.Expected.HybridConnectionNamespaceName {
			t.Fatalf("Expected %s but got %s for Namespace", v.Expected.HybridConnectionNamespaceName, actual.HybridConnectionNamespaceName)
		}
	}
}
