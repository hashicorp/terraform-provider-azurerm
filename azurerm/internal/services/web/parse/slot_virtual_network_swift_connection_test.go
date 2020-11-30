package parse

import (
	"testing"
)

func TestSlotVirtualNetworkSwiftConnectionID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *SlotVirtualNetworkSwiftConnectionId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Segemt",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000",
			Expected: nil,
		},
		{
			Name:     "No Sites Segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web",
			Expected: nil,
		},
		{
			Name:     "No Slot Segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/instance1",
			Expected: nil,
		},
		{
			Name:     "Valid Network Association But No Slot Segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/instance1/networkconfig/virtualNetwork",
			Expected: nil,
		},
		{
			Name:  "Slot Virtual Network Association",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/instance1/slots/stageing/networkconfig/virtualNetwork",
			Expected: &SlotVirtualNetworkSwiftConnectionId{
				SiteName:      "instance1",
				ResourceGroup: "mygroup1",
				SlotName:      "stageing",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := SlotVirtualNetworkSwiftConnectionID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.SiteName != v.Expected.SiteName {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.SiteName, actual.SiteName)
		}
		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
		if actual.SlotName != v.Expected.SlotName {
			t.Fatalf("Expected %q but got %q for SlotName", v.Expected.SlotName, actual.SlotName)
		}
	}
}
