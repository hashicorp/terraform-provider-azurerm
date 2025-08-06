// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

func TestMaintenanceAssignmentVirtualMachineID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *MaintenanceAssignmentVirtualMachineId
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
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resGroup1/providers/microsoft.compute/virtualMachines/",
			Error: true,
		},
		{
			Name:  "No Maintenance Assignment Segment",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resGroup1/providers/microsoft.compute/virtualMachines/vm1/providers/Microsoft.Maintenance/",
			Error: true,
		},
		{
			Name:  "No Maintenance Assignment name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resGroup1/providers/microsoft.compute/virtualMachines/vm1/providers/Microsoft.Maintenance/configurationAssignments/",
			Error: true,
		},
		{
			Name:  "ID of Maintenance Assignment to vm",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resGroup1/providers/microsoft.compute/virtualMachines/vm1/providers/Microsoft.Maintenance/configurationAssignments/assign1",
			Error: false,
			Expect: &MaintenanceAssignmentVirtualMachineId{
				VirtualMachineIdRaw: "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resGroup1/providers/microsoft.compute/virtualMachines/vm1",
				VirtualMachineId:    pointer.To(commonids.NewVirtualMachineID("00000000-0000-0000-0000-000000000000", "resGroup1", "vm1")),
				Name:                "assign1",
			},
		},
		{
			Name:  "Wrong Casing",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resGroup1/providers/microsoft.compute/virtualMachines/vm1/providers/Microsoft.Maintenance/ConfigurationAssignments/assign1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.Name)

		actual, err := MaintenanceAssignmentVirtualMachineID(v.Input)
		if err != nil {
			if v.Expect == nil {
				continue
			}
			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expect.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.Name, actual.Name)
		}

		if actual.VirtualMachineIdRaw != v.Expect.VirtualMachineIdRaw {
			t.Fatalf("Expected %q but got %q for VirtualMachineIdRaw", v.Expect.VirtualMachineIdRaw, actual.VirtualMachineIdRaw)
		}

		if !reflect.DeepEqual(v.Expect.VirtualMachineId, actual.VirtualMachineId) {
			t.Fatalf("Expected %+v but got %+v", v.Expect.VirtualMachineId, actual.VirtualMachineId)
		}
	}
}
