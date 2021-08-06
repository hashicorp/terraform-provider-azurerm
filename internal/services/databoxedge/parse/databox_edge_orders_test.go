package parse

import (
	"testing"
)

func TestDataboxEdgeOrderIDFormatter(t *testing.T) {
	subscriptionId := "12345678-1234-5678-1234-123456789012"
	id := NewDataboxEdgeOrderID("resourceGroup1", "device1")
	actual := id.ID(subscriptionId)
	expected := "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.DataBoxEdge/dataBoxEdgeDevices/device1/orders/default"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestDataboxEdgeOrderID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *DataboxEdgeOrderId
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
			Name:     "No Resource Groups Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Expected: nil,
		},
		{
			Name:     "Resource Group ID",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/",
			Expected: nil,
		},
		{
			Name:     "Missing Order Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.DataBoxEdge/dataBoxEdgeDevices/device1/orders",
			Expected: nil,
		},
		{
			Name:  "databoxedge Order ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.DataBoxEdge/dataBoxEdgeDevices/device1/orders/default",
			Expected: &DataboxEdgeOrderId{
				ResourceGroup: "resourceGroup1",
				DeviceName:    "device1",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.DataBoxEdge/dataBoxEdgeDevices/device1/Orders/default",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.Name)

		actual, err := DataboxEdgeOrderID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}
			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.DeviceName != v.Expected.DeviceName {
			t.Fatalf("Expected %q but got %q for DeviceName", v.Expected.DeviceName, actual.DeviceName)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
