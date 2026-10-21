// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package ip_group

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/azurefirewalls"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/firewallpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/ipgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIpGroupUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.IPGroups
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	for _, fw := range d.Get("firewall_ids").([]any) {
		id, err := azurefirewalls.ParseAzureFirewallID(fw.(string))
		if err != nil {
			return fmt.Errorf("parsing Azure Firewall ID %q: %+v", fw, err)
		}
		locks.ByName(id.AzureFirewallName, firewall.AzureFirewallResourceName)
		defer locks.UnlockByName(id.AzureFirewallName, firewall.AzureFirewallResourceName)
	}

	for _, fwpol := range d.Get("firewall_policy_ids").([]any) {
		id, err := firewallpolicies.ParseFirewallPolicyID(fwpol.(string))
		if err != nil {
			return fmt.Errorf("parsing Azure Firewall Policy ID %q: %+v", fwpol, err)
		}
		locks.ByName(id.FirewallPolicyName, firewall.AzureFirewallPolicyResourceName)
		defer locks.UnlockByName(id.FirewallPolicyName, firewall.AzureFirewallPolicyResourceName)
	}

	id, err := ipgroups.ParseIPGroupID(d.Id())
	if err != nil {
		return err
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	existing, err := client.Get(ctx, *id, ipgroups.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("cidrs") {
		payload.Properties.IPAddresses = pluginsdk.ExpandStringSlice(d.Get("cidrs").(*pluginsdk.Set).List())
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceIpGroupRead(d, meta)
}
