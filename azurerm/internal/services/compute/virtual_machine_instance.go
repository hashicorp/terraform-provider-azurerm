package compute

import (
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
)

// shouldBootVirtualMachine determines if the Virtual Machine should be started after
// the Virtual Machine has been shut down for maintenance. This means that Virtual Machines
// which are already stopped can be updated but will not be started
// nolint: deadcode unused
func shouldBootVirtualMachine(instanceView compute.VirtualMachineInstanceView) bool {
	if instanceView.Statuses != nil {
		for _, status := range *instanceView.Statuses {
			if status.Code == nil {
				continue
			}

			// could also be the provisioning state which we're not bothered with here
			state := strings.ToLower(*status.Code)
			if !strings.HasPrefix(state, "PowerState/") {
				continue
			}

			state = strings.TrimPrefix(state, "powerstate/")
			switch strings.ToLower(state) {
			case "deallocating":
			case "deallocated":
			case "stopped":
			case "stopping":
				return true
			}
		}
	}

	return false
}
