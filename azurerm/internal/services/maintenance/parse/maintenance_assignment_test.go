package parse

import (
	"reflect"
	"testing"
)

func TestMaintenanceAssignmentID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *MaintenanceAssignmentId
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
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resGroup1/providers/microsoft.compute/virtualmachines/",
			Error: true,
		},
		{
			Name:  "No Maintenance Assignment Segment",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resGroup1/providers/microsoft.compute/virtualmachines/vm1/providers/Microsoft.Maintenance/",
			Error: true,
		},
		{
			Name:  "No Maintenance Assignment name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resGroup1/providers/microsoft.compute/virtualmachines/vm1/providers/Microsoft.Maintenance/configurationAssignments/",
			Error: true,
		},
		{
			Name:  "ID of Maintenance Assignment to vm",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resGroup1/providers/microsoft.compute/virtualmachines/vm1/providers/Microsoft.Maintenance/configurationAssignments/assign1",
			Error: false,
			Expect: &MaintenanceAssignmentId{
				TargetResourceId: ScopeResource{
					id:               "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resGroup1/providers/microsoft.compute/virtualmachines/vm1",
					ResourceGroup:    "resGroup1",
					ResourceProvider: "microsoft.compute",
					ResourceType:     "virtualmachines",
					ResourceName:     "vm1",
				},
				Name: "assign1",
			},
		},
		{
			Name:  "ID of Maintenance Assignment to dedicated host",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resGroup1/providers/microsoft.compute/hostGroups/group1/hosts/host1/providers/Microsoft.Maintenance/configurationAssignments/assign1",
			Error: false,
			Expect: &MaintenanceAssignmentId{
				TargetResourceId: ScopeInResource{
					id:                 "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resGroup1/providers/microsoft.compute/hostGroups/group1/hosts/host1",
					ResourceGroup:      "resGroup1",
					ResourceProvider:   "microsoft.compute",
					ResourceParentType: "hostGroups",
					ResourceParentName: "group1",
					ResourceType:       "hosts",
					ResourceName:       "host1",
				},
				Name: "assign1",
			},
		},
		{
			Name:  "Wrong Casing",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resGroup1/providers/microsoft.compute/virtualmachines/vm1/providers/Microsoft.Maintenance/ConfigurationAssignments/assign1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.Name)

		actual, err := MaintenanceAssignmentID(v.Input)
		if err != nil {
			if v.Expect == nil {
				continue
			}
			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expect.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.Name, actual.Name)
		}

		if !reflect.DeepEqual(v.Expect.TargetResourceId, actual.TargetResourceId) {
			t.Fatalf("Expected %+v but got %+v", v.Expect.TargetResourceId, actual.TargetResourceId)
		}
	}
}
