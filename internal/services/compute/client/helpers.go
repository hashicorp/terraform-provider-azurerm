package client

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"log"
)

func (c *Client) CancelRollingUpgradesBeforeDeletion(ctx context.Context, id parse.VirtualMachineScaleSetId) error {
	resp, err := c.VMScaleSetRollingUpgradesClient.GetLatest(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		// no rolling upgrades are running so skipping attempt to cancel them before deletion
		if utils.ResponseWasNotFound(resp.Response) {
			return nil
		}
		return fmt.Errorf("retrieving rolling updates for %s: %+v", id, err)
	}

	future, err := c.VMScaleSetRollingUpgradesClient.Cancel(ctx, id.ResourceGroup, id.Name)

	if err := future.WaitForCompletionRef(ctx, c.VMScaleSetExtensionsClient.Client); err != nil {
		return fmt.Errorf("waiting for cancelling of rolling upgrades for %s: %+v", id, err)
	}
	log.Printf("[DEBUG] cancelled Virtual Machine Scale Set Rolling Upgrades for %s.", id)
	return nil
}
