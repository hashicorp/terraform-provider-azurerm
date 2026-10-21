// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package ip_group

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/azurefirewalls"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/firewallpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/ipgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIpGroupDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.IPGroups
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := ipgroups.ParseIPGroupID(d.Id())
	if err != nil {
		return err
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	resp, err := client.Get(ctx, *id, ipgroups.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			return nil
		}
		return fmt.Errorf("retrieving %s : %+v", *id, err)
	}

	if resp.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if resp.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	for _, fw := range *resp.Model.Properties.Firewalls {
		fwID, err := azurefirewalls.ParseAzureFirewallID(pointer.From(fw.Id))
		if err != nil {
			return fmt.Errorf("parsing Azure Firewall ID %q: %+v", pointer.From(fw.Id), err)
		}
		locks.ByName(fwID.AzureFirewallName, firewall.AzureFirewallResourceName)
		defer locks.UnlockByName(fwID.AzureFirewallName, firewall.AzureFirewallResourceName)
	}

	for _, fwpol := range *resp.Model.Properties.FirewallPolicies {
		polID, err := firewallpolicies.ParseFirewallPolicyID(pointer.From(fwpol.Id))
		if err != nil {
			return fmt.Errorf("parsing Azure Firewall Policy ID %q: %+v", *fwpol.Id, err)
		}
		locks.ByName(polID.FirewallPolicyName, firewall.AzureFirewallPolicyResourceName)
		defer locks.UnlockByName(polID.FirewallPolicyName, firewall.AzureFirewallPolicyResourceName)
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return err
}
