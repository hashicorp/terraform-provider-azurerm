package compute

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/tombuildsstuff/kermit/sdk/compute/2022-08-01/compute"
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
			} else if strings.EqualFold(state, "starting") {
				return true
			}
		}
	}

	return false
}

func virtualMachinePowerStateRefreshFunc(ctx context.Context, client *compute.VirtualMachinesClient, id parse.VirtualMachineId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instanceView, err := client.InstanceView(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return nil, "", fmt.Errorf("retrieving InstanceView for %q: %+v", id, err)
		}

		if instanceView.Statuses != nil {
			for _, status := range *instanceView.Statuses {
				if status.Code != nil {
					state := strings.ToLower(*status.Code)
					if strings.HasPrefix(state, "powerstate/") {
						state = strings.TrimPrefix(state, "powerstate/")
						return state, strings.ToLower(state), nil
					}
				}
			}
		}

		return nil, "starting", nil
	}
}
