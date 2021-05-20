package parse

import (
	"reflect"
	"testing"

	parseCompute "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
)

func TestMaintenanceAssignmentDedicatedHostID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *MaintenanceAssignmentDedicatedHostId
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
			Name:  "No target resource type",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resGroup1/providers/microsoft.compute/",
			Error: true,
		},
		{
			Name:  "No target resource name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resGroup1/providers/microsoft.compute/hostGroups/group1/hosts",
			Error: true,
		},
		{
			Name:  "No Maintenance Assignment Segment",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resGroup1/providers/microsoft.compute/hostGroups/group1/hosts/host1/providers/Microsoft.Maintenance/",
			Error: true,
		},
		{
			Name:  "No Maintenance Assignment name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resGroup1/providers/microsoft.compute/hostGroups/group1/hosts/host1/providers/Microsoft.Maintenance/configurationAssignments/",
			Error: true,
		},
		{
			Name:  "ID of Maintenance Assignment to dedicated host",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resGroup1/providers/microsoft.compute/hostGroups/group1/hosts/host1/providers/Microsoft.Maintenance/configurationAssignments/assign1",
			Error: false,
			Expect: &MaintenanceAssignmentDedicatedHostId{
				DedicatedHostIdRaw: "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resGroup1/providers/microsoft.compute/hostGroups/group1/hosts/host1",
				DedicatedHostId: &parseCompute.DedicatedHostId{
					SubscriptionId: "00000000-0000-0000-0000-000000000000",
					ResourceGroup:  "resGroup1",
					HostGroupName:  "group1",
					HostName:       "host1",
				},
				Name: "assign1",
			},
		},
		{
			Name:  "Wrong Casing",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resGroup1/providers/microsoft.compute/hostGroups/group1/hosts/host1/providers/Microsoft.Maintenance/ConfigurationAssignments/assign1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.Name)

		actual, err := MaintenanceAssignmentDedicatedHostID(v.Input)
		if err != nil {
			if v.Expect == nil {
				continue
			}
			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expect.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.Name, actual.Name)
		}

		if actual.DedicatedHostIdRaw != v.Expect.DedicatedHostIdRaw {
			t.Fatalf("Expected %q but got %q for DedicatedHostIdRaw", v.Expect.DedicatedHostIdRaw, actual.DedicatedHostIdRaw)
		}

		if !reflect.DeepEqual(v.Expect.DedicatedHostId, actual.DedicatedHostId) {
			t.Fatalf("Expected %+v but got %+v", v.Expect.DedicatedHostId, actual.DedicatedHostId)
		}
	}
}
