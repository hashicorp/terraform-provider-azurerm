package parse

import (
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
				TargetResourceId: &TargetResourceId{
					HasParentResource: false,
					ResourceGroup:     "resGroup1",
					ResourceProvider:  "microsoft.compute",
					ResourceType:      "virtualmachines",
					ResourceName:      "vm1",
				},
				ResourceId: "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resGroup1/providers/microsoft.compute/virtualmachines/vm1",
				Name:       "assign1",
			},
		},
		{
			Name:  "ID of Maintenance Assignment to dedicated host",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resGroup1/providers/microsoft.compute/hostGroups/group1/hosts/host1/providers/Microsoft.Maintenance/configurationAssignments/assign1",
			Error: false,
			Expect: &MaintenanceAssignmentId{
				TargetResourceId: &TargetResourceId{
					HasParentResource:  true,
					ResourceGroup:      "resGroup1",
					ResourceProvider:   "microsoft.compute",
					ResourceParentType: "hostGroups",
					ResourceParentName: "group1",
					ResourceType:       "hosts",
					ResourceName:       "host1",
				},
				ResourceId: "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resGroup1/providers/microsoft.compute/virtualmachines/hostGroups/group1/hosts/host1",
				Name:       "assign1",
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

		if actual.ResourceGroup != v.Expect.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expect.ResourceGroup, actual.ResourceGroup)
		}

		if actual.ResourceProvider != v.Expect.ResourceProvider {
			t.Fatalf("Expected %q but got %q for ResourceProvider", v.Expect.ResourceProvider, actual.ResourceProvider)
		}

		if actual.ResourceType != v.Expect.ResourceType {
			t.Fatalf("Expected %q but got %q for ResourceType", v.Expect.ResourceType, actual.ResourceType)
		}

		if actual.ResourceName != v.Expect.ResourceName {
			t.Fatalf("Expected %q but got %q for ResourceName", v.Expect.ResourceName, actual.ResourceName)
		}

		if actual.Name != v.Expect.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.Name, actual.Name)
		}
	}
}
