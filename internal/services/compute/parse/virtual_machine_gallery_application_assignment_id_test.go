// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplicationversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
)

func TestVirtualMachineGalleryApplicationAssignmentId(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Expect *VirtualMachineGalleryApplicationAssignmentId
		Error  bool
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "One Segment",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000001/resourceGroups/group1/providers/Microsoft.Compute/virtualMachines/virtualMachine1",
			Error: true,
		},
		{
			Name:  "Two Segments Invalid Virtual Machine ID",
			Input: "hello|/subscriptions/00000000-0000-0000-0000-000000000002/resourceGroups/group2/providers/Microsoft.Compute/galleries/gallery2/applications/application2/versions/0.0.2",
			Error: true,
		},
		{
			Name:  "Two Segments Invalid Gallery Application Version ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000001/resourceGroups/group1/providers/Microsoft.Compute/virtualMachines/virtualMachine1|world",
			Error: true,
		},
		{
			Name:  "Virtual Machine ID / Gallery Application Version ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000001/resourceGroups/group1/providers/Microsoft.Compute/virtualMachines/virtualMachine1|/subscriptions/00000000-0000-0000-0000-000000000002/resourceGroups/group2/providers/Microsoft.Compute/galleries/gallery2/applications/application2/versions/0.0.2",
			Error: false,
			Expect: &VirtualMachineGalleryApplicationAssignmentId{
				VirtualMachineId: virtualmachines.VirtualMachineId{
					SubscriptionId:     "00000000-0000-0000-0000-000000000001",
					ResourceGroupName:  "group1",
					VirtualMachineName: "virtualMachine1",
				},
				GalleryApplicationVersionId: galleryapplicationversions.ApplicationVersionId{
					SubscriptionId:    "00000000-0000-0000-0000-000000000002",
					ResourceGroupName: "group2",
					GalleryName:       "gallery2",
					ApplicationName:   "application2",
					VersionName:       "0.0.2",
				},
			},
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := VirtualMachineGalleryApplicationAssignmentID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.VirtualMachineId.VirtualMachineName != v.Expect.VirtualMachineId.VirtualMachineName {
			t.Fatalf("Expected %q but got %q for VirtualMachineName", v.Expect.VirtualMachineId.VirtualMachineName, actual.VirtualMachineId.VirtualMachineName)
		}

		if actual.GalleryApplicationVersionId.GalleryName != v.Expect.GalleryApplicationVersionId.GalleryName {
			t.Fatalf("Expected %q but got %q for GalleryName", v.Expect.GalleryApplicationVersionId.GalleryName, actual.GalleryApplicationVersionId.GalleryName)
		}

		if actual.GalleryApplicationVersionId.ApplicationName != v.Expect.GalleryApplicationVersionId.ApplicationName {
			t.Fatalf("Expected %q but got %q for ApplicationName", v.Expect.GalleryApplicationVersionId.ApplicationName, actual.GalleryApplicationVersionId.ApplicationName)
		}

		if actual.GalleryApplicationVersionId.VersionName != v.Expect.GalleryApplicationVersionId.VersionName {
			t.Fatalf("Expected %q but got %q for VersionName", v.Expect.GalleryApplicationVersionId.VersionName, actual.GalleryApplicationVersionId.VersionName)
		}
	}
}
