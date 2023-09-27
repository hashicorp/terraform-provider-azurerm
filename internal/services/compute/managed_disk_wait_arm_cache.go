package compute

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// waitForManagedDiskArmCache repeatedly calls the resource group level list API to query the ARM cache for the specified managed disk,
// until it (dis)appear in the ARM cache, based on "shouldExist".
func waitForManagedDiskArmCache(ctx context.Context, client *resources.Client, diskId disks.DiskId, shouldExist bool) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context was missing a deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Done"},
		Refresh: func() (interface{}, string, error) {
			filter := fmt.Sprintf("name eq '%s' and resourceType eq 'Microsoft.Compute/disks'", diskId.DiskName)
			resp, err := client.ListByResourceGroup(ctx, diskId.ResourceGroupName, filter, "", nil)
			if err != nil {
				return nil, "", fmt.Errorf("polling for %s within the resource group: %+v", diskId, err)
			}
			if shouldExist {
				if len(resp.Values()) == 0 {
					return resp, "Pending", nil
				} else {
					return resp, "Done", nil
				}
			} else {
				if len(resp.Values()) == 0 {
					return resp, "Done", nil
				} else {
					return resp, "Pending", nil
				}
			}
		},
		ContinuousTargetOccurence: 5,
		PollInterval:              5 * time.Second,
		Timeout:                   time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return err
	}

	return nil
}
