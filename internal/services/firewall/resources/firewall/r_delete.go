// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package firewall

import (
	"fmt"
	"log"
	"slices"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/azurefirewalls"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/firewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceFirewallDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.AzureFirewalls
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azurefirewalls.ParseAzureFirewallID(d.Id())
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			// deleted outside of TF
			log.Printf("[DEBUG] Firewall %q was not found in Resource Group %q - assuming removed!", id.AzureFirewallName, id.ResourceGroupName)
			return nil
		}

		return fmt.Errorf("retrieving Firewall %s : %+v", *id, err)
	}

	subnetNamesToLock := make([]string, 0)
	virtualNetworkNamesToLock := make([]string, 0)
	if model := read.Model; model != nil {
		if props := model.Properties; props != nil {
			if configs := props.IPConfigurations; configs != nil {
				for _, config := range *configs {
					if config.Properties == nil || config.Properties.Subnet == nil || config.Properties.Subnet.Id == nil {
						continue
					}

					parsedSubnetID, err2 := commonids.ParseSubnetID(*config.Properties.Subnet.Id)
					if err2 != nil {
						return err2
					}

					if !slices.Contains(subnetNamesToLock, parsedSubnetID.SubnetName) {
						subnetNamesToLock = append(subnetNamesToLock, parsedSubnetID.SubnetName)
					}

					if !slices.Contains(virtualNetworkNamesToLock, parsedSubnetID.VirtualNetworkName) {
						virtualNetworkNamesToLock = append(virtualNetworkNamesToLock, parsedSubnetID.VirtualNetworkName)
					}
				}
			}

			if mconfig := props.ManagementIPConfiguration; mconfig != nil {
				if mconfig.Properties != nil && mconfig.Properties.Subnet != nil && mconfig.Properties.Subnet.Id != nil {
					parsedSubnetID, err2 := commonids.ParseSubnetID(*mconfig.Properties.Subnet.Id)
					if err2 != nil {
						return err2
					}

					if !slices.Contains(subnetNamesToLock, parsedSubnetID.SubnetName) {
						subnetNamesToLock = append(subnetNamesToLock, parsedSubnetID.SubnetName)
					}

					if !slices.Contains(virtualNetworkNamesToLock, parsedSubnetID.VirtualNetworkName) {
						virtualNetworkNamesToLock = append(virtualNetworkNamesToLock, parsedSubnetID.VirtualNetworkName)
					}
				}
			}
		}

		if read.Model.Properties != nil && read.Model.Properties.FirewallPolicy != nil && read.Model.Properties.FirewallPolicy.Id != nil {
			id, err := firewallpolicies.ParseFirewallPolicyIDInsensitively(*read.Model.Properties.FirewallPolicy.Id)
			if err != nil {
				return err
			}
			locks.ByName(id.FirewallPolicyName, AzureFirewallPolicyResourceName)
			defer locks.UnlockByName(id.FirewallPolicyName, AzureFirewallPolicyResourceName)
		}

		locks.ByName(id.AzureFirewallName, AzureFirewallResourceName)
		defer locks.UnlockByName(id.AzureFirewallName, AzureFirewallResourceName)

		locks.MultipleByName(&virtualNetworkNamesToLock, VirtualNetworkResourceName)
		defer locks.UnlockMultipleByName(&virtualNetworkNamesToLock, VirtualNetworkResourceName)

		locks.MultipleByName(&subnetNamesToLock, SubnetResourceName)
		defer locks.UnlockMultipleByName(&subnetNamesToLock, SubnetResourceName)

		// todo see if this is still needed this way
		/*
			// Change this back to using the SDK method once https://github.com/Azure/azure-sdk-for-go/issues/17013 is addressed.
			future, err := azuresdkhacks.DeleteFirewall(ctx, client, id.ResourceGroup, id.AzureFirewallName)
			if err != nil {
				return fmt.Errorf("deleting Azure Firewall %s : %+v", *id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for the deletion of Azure Firewall %s : %+v", *id, err)
			}
		*/

		return client.DeleteThenPoll(ctx, *id)
	}

	return err
}
