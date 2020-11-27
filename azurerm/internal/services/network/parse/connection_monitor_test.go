package parse

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = NetworkWatcherId{}

func TestNetworkWatcherIDFormatter(t *testing.T) {
	subscriptionId := "12345678-1234-5678-1234-123456789012"
	actual := NewNetworkWatcherID("group1", "watcher1").ID(subscriptionId)
	expected := "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Network/networkWatchers/watcher1"

	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestNetworkConnectionMonitorID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *ConnectionMonitorId
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
			Name:  "Missing Network Connection Monitor Key",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/networkWatchers/watcher1/connectionMonitors",
			Error: true,
		},
		{
			Name:  "Namespace Network Connection Monitor Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/networkWatchers/watcher1/connectionMonitors/connectionMonitor1",
			Error: false,
			Expect: &ConnectionMonitorId{
				ResourceGroup:      "group1",
				NetworkWatcherName: "watcher1",
				Name:               "connectionMonitor1",
			},
		},
		{
			Name:  "Wrong Segment",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/networkWatchers/watcher1/NetworkConnectionMonitors/connectionMonitor1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := ConnectionMonitorID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expect.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.Name, actual.Name)
		}

		if actual.NetworkWatcherName != v.Expect.NetworkWatcherName {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.NetworkWatcherName, actual.NetworkWatcherName)
		}

		if actual.ResourceGroup != v.Expect.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expect.ResourceGroup, actual.ResourceGroup)
		}
	}
}
