package parse

import (
	"testing"
)

func TestNetworkPacketCaptureID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *NetworkPacketCaptureId
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
			Name:  "Missing Network Watcher Key",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/networkWatchers/",
			Error: true,
		},
		{
			Name:  "Missing Network Watcher Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/networkWatchers/watcher1",
			Error: true,
		},
		{
			Name:  "Missing Network Packet Capture Key",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/networkWatchers/watcher1/packetCaptures",
			Error: true,
		},
		{
			Name:  "Namespace Network Packet Capture Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/networkWatchers/watcher1/packetCaptures/packetCapture1",
			Error: false,
			Expect: &NetworkPacketCaptureId{
				ResourceGroup: "group1",
				WatcherName:   "watcher1",
				Name:          "packetCapture1",
			},
		},
		{
			Name:  "Wrong Segment",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/networkWatchers/watcher1/NetworkPacketCaptures/packetCapture1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := NetworkPacketCaptureID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expect.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.Name, actual.Name)
		}

		if actual.WatcherName != v.Expect.WatcherName {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.WatcherName, actual.WatcherName)
		}

		if actual.ResourceGroup != v.Expect.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expect.ResourceGroup, actual.ResourceGroup)
		}
	}
}
