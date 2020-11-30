package web

import "testing"

func TestParseAppServiceHybridConnectionID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *AppServiceHybridConnectionResourceID
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
			Expected: &AppServiceHybridConnectionResourceID{
				ResourceGroup: "testResourceGroup1",
				Name:          "testRelay1",
				AppName:       "testApp1",
				Namespace:     "testNamespace1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := ParseAppServiceHybridConnectionID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %s but got %s for Name", v.Expected.Name, actual.Name)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %s but got %s for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}

		if actual.AppName != v.Expected.AppName {
			t.Fatalf("Expected %s but got %s for AppName", v.Expected.AppName, actual.AppName)
		}

		if actual.Namespace != v.Expected.Namespace {
			t.Fatalf("Expected %s but got %s for Namespace", v.Expected.Namespace, actual.Namespace)
		}
	}
}
