// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package firewall

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/azurefirewalls"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceFirewallRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.AzureFirewalls
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azurefirewalls.ParseAzureFirewallID(d.Id())
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			log.Printf("[DEBUG] Firewall %q was not found in Resource Group %q - removing from state!", id.AzureFirewallName, id.ResourceGroupName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on Azure Firewall %s : %+v", *id, err)
	}
	return resourceFirewallSetFlatten(d, id, read.Model)
}

func validateFirewallIPConfigurationSettings(configs []any) error {
	if len(configs) == 0 {
		return nil
	}

	subnetNumber := 0

	for _, configRaw := range configs {
		data := configRaw.(map[string]any)
		if subnet, exist := data["subnet_id"].(string); exist && subnet != "" {
			subnetNumber++
		}
	}

	if subnetNumber != 1 {
		return fmt.Errorf(`the "ip_configuration" is invalid, %d "subnet_id" have been set, one "subnet_id" should be set among all "ip_configuration" blocks`, subnetNumber)
	}

	return nil
}
