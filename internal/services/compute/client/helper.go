package client

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func (c *Client) CancelRollingUpgradesBeforeDeletion(ctx context.Context, resourceGroupName string, vmScaleSetName string) error {
	future, err := c.VMScaleSetRollingUpgradesClient.Cancel(ctx, resourceGroupName, vmScaleSetName)

	// If rolling upgrades haven't been run (when VMSS are just provisioned with rolling upgrades but no extensions, auto-scaling are run )
	// we can not cancel rolling upgrades
	// API call :: GET https://management.azure.com/subscriptions/{subId}/resourceGroups/{rgName}/providers/Microsoft.Compute/virtualMachineScaleSets/{vmSSName}/rollingUpgrades/latest?api-version=2021-07-01
	// Azure API throws 409 conflict error saying "The entity was not found in this Azure location."
	// If the above error message matches, we identify and move forward to delete the VMSS
	// in all other cases, it just cancels the rolling upgrades and move ahead to delete the VMSS
	if err != nil && !(future.Response().StatusCode == http.StatusConflict && strings.Contains(err.Error(), "There is no ongoing Rolling Upgrade to cancel.")) {
		return fmt.Errorf("error cancelling rolling upgrades of Virtual Machine Scale Set %q (Resource Group %q): %+v", vmScaleSetName, resourceGroupName, err)
	}

	log.Printf("[DEBUG] Waiting for deletion of Virtual Machine Scale Set %q (Resource Group %q)..", vmScaleSetName, resourceGroupName)
	if err := future.WaitForCompletionRef(ctx, c.VMScaleSetExtensionsClient.Client); err != nil && !(future.Response().StatusCode == http.StatusConflict && strings.Contains(err.Error(), "There is no ongoing Rolling Upgrade to cancel.")) {
		return fmt.Errorf("waiting for cancelling rolling upgrades of Virtual Machine Scale Set %q (Resource Group %q): %+v", vmScaleSetName, resourceGroupName, err)
	}
	log.Printf("[DEBUG] cancelled Virtual Machine Scale Set Rolling Upgrades %q (Resource Group %q).", vmScaleSetName, resourceGroupName)
	return nil
}
