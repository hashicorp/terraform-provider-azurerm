package parse

import (
	"testing"
)

func TestIotSecurityDeviceGroupID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *IotSecurityDeviceGroupId
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "No Security resource provider",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1",
			Error: true,
		},

		{
			Name:  "No target resource Segment",
			Input: "/providers/Microsoft.Security/deviceSecurityGroups/group1",
			Error: true,
		},
		{
			Name:  "No Iot Security Device Group Segment",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/providers/Microsoft.Security/",
			Error: true,
		},
		{
			Name:  "No Iot Security Device Group name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/providers/Microsoft.Security/deviceSecurityGroups/",
			Error: true,
		},
		{
			Name:  "ID of Iot Security Device Group",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/providers/Microsoft.Security/deviceSecurityGroups/group1",
			Error: false,
			Expect: &IotSecurityDeviceGroupId{
				IotHubID: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1",
				Name:     "group1",
			},
		},
		{
			Name:  "Wrong Casing",
			Input: "/SUBSCRIPTIONS/00000000-0000-0000-0000-000000000000/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.DEVICES/IOTHUBS/HUB1/PROVIDERS/MICROSOFT.SECURITY/DEVICESECURITYGROUPS/GROUP1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.Name)

		actual, err := IotSecurityDeviceGroupID(v.Input)
		if err != nil {
			if v.Expect == nil {
				continue
			}
			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expect.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.Name, actual.Name)
		}
		if actual.IotHubID != v.Expect.IotHubID {
			t.Fatalf("Expected %q but got %q for IotHubID", v.Expect.IotHubID, actual.IotHubID)
		}
	}
}
