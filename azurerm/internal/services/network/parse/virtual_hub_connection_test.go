package parse

import (
	"testing"
)

func TestParseVirtualHubConnection(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *VirtualHubConnectionResourceID
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No Virtual Hubs Segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo",
			Expected: nil,
		},
		{
			Name:     "No Virtual Hubs Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/virtualHubs/",
			Expected: nil,
		},
		{
			Name:     "No Hub Network Connections Segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/virtualHubs/example",
			Expected: nil,
		},
		{
			Name:     "No Virtual Hubs Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/virtualHubs/hubVirtualNetworkConnections/",
			Expected: nil,
		},
		{
			Name:  "Completed",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/virtualHubs/example/hubVirtualNetworkConnections/connection1",
			Expected: &VirtualHubConnectionResourceID{
				Name:           "connection1",
				VirtualHubName: "example",
				ResourceGroup:  "foo",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := ParseVirtualHubConnectionID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
