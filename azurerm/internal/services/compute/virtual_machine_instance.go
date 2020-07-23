package compute

import (
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-12-01/compute"
)

// virtualMachineShouldBeStarted determines if the Virtual Machine should be started after
// the Virtual Machine has been shut down for maintenance. This means that Virtual Machines
// which are already stopped can be updated but will not be started
// nolint: deadcode unused
func virtualMachineShouldBeStarted(instanceView compute.VirtualMachineInstanceView) bool {
	if instanceView.Statuses != nil {
		for _, status := range *instanceView.Statuses {
			if status.Code == nil {
				continue
			}

			// could also be the provisioning state which we're not bothered with here
			state := strings.ToLower(*status.Code)
			if !strings.HasPrefix(state, "powerstate/") {
				continue
			}

			state = strings.TrimPrefix(state, "powerstate/")
			if strings.EqualFold(state, "running") {
				return true
			}
		}
	}

	return false
}
