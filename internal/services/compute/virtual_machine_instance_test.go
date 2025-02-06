// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
)

func TestVirtualMachineShouldBeStarted(t *testing.T) {
	buildInstanceViewStatus := func(statuses ...string) *[]virtualmachines.InstanceViewStatus {
		results := make([]virtualmachines.InstanceViewStatus, 0)

		for _, v := range statuses {
			results = append(results, virtualmachines.InstanceViewStatus{
				Code: pointer.To(v),
			})
		}

		return &results
	}

	testCases := []struct {
		Name     string
		Input    *[]virtualmachines.InstanceViewStatus
		Expected bool
	}{
		{
			Name:     "None",
			Expected: false,
			Input:    nil,
		},
		{
			Name:     "No Power State",
			Input:    buildInstanceViewStatus("ProvisioningStatus/Creating"),
			Expected: false,
		},
		{
			Name:     "Running",
			Input:    buildInstanceViewStatus("ProvisioningStatus/succeeded", "PowerState/running"),
			Expected: true,
		},
		{
			Name:     "Deallocated",
			Input:    buildInstanceViewStatus("ProvisioningStatus/succeeded", "PowerState/deallocated"),
			Expected: false,
		},
		{
			Name:     "Deallocating",
			Input:    buildInstanceViewStatus("ProvisioningStatus/updating", "PowerState/deallocating"),
			Expected: false,
		},
		{
			Name:     "Stopped",
			Input:    buildInstanceViewStatus("ProvisioningStatus/updating", "PowerState/stopped"),
			Expected: false,
		},
		{
			Name:     "Stopping",
			Input:    buildInstanceViewStatus("ProvisioningStatus/updating", "PowerState/stopping"),
			Expected: false,
		},
		{
			Name:     "Failed",
			Input:    buildInstanceViewStatus("ProvisioningStatus/failed", "PowerState/failed"),
			Expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running %q..", testCase.Name)

		instanceView := virtualmachines.VirtualMachineInstanceView{
			Statuses: testCase.Input,
		}
		result := virtualMachineShouldBeStarted(&instanceView)
		if result != testCase.Expected {
			t.Fatalf("Expected %t but got %t", testCase.Expected, result)
		}
	}
}
