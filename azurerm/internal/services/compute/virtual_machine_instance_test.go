package compute

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestVirtualMachineShouldBeStarted(t *testing.T) {
	buildInstanceViewStatus := func(statuses ...string) *[]compute.InstanceViewStatus {
		results := make([]compute.InstanceViewStatus, 0)

		for _, v := range statuses {
			results = append(results, compute.InstanceViewStatus{
				Code: utils.String(v),
			})
		}

		return &results
	}

	testCases := []struct {
		Name     string
		Input    *[]compute.InstanceViewStatus
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

		instanceView := compute.VirtualMachineInstanceView{
			Statuses: testCase.Input,
		}
		result := virtualMachineShouldBeStarted(instanceView)
		if result != testCase.Expected {
			t.Fatalf("Expected %t but got %t", testCase.Expected, result)
		}
	}
}
