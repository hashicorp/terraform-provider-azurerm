package parse

import (
	"testing"
)

func TestDeviceSecurityGroupID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *DeviceSecurityGroupId
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
			Name:  "No Device Security Group Segment",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/providers/Microsoft.Security/",
			Error: true,
		},
		{
			Name:  "No Device Security Group name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/providers/Microsoft.Security/deviceSecurityGroups/",
			Error: true,
		},
		{
			Name:  "ID of Device Security Group",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/providers/Microsoft.Security/deviceSecurityGroups/group1",
			Error: false,
			Expect: &DeviceSecurityGroupId{
				TargetResourceID: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1",
				Name:             "group1",
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

		actual, err := DeviceSecurityGroupID(v.Input)
		if err != nil {
			if v.Expect == nil {
				continue
			}
			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expect.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.Name, actual.Name)
		}
		if actual.TargetResourceID != v.Expect.TargetResourceID {
			t.Fatalf("Expected %q but got %q for TargetResourceID", v.Expect.TargetResourceID, actual.TargetResourceID)
		}
	}
}
