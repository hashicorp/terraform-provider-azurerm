// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package firewall_network_rule_collection

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/azurefirewalls"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceFirewallNetworkRuleCollectionDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.AzureFirewalls
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FirewallNetworkRuleCollectionID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.AzureFirewallName, AzureFirewallResourceName)
	defer locks.UnlockByName(id.AzureFirewallName, AzureFirewallResourceName)

	firewallId := azurefirewalls.NewAzureFirewallID(id.SubscriptionId, id.ResourceGroup, id.AzureFirewallName)

	firewall, err := client.Get(ctx, firewallId)
	if err != nil {
		if response.WasNotFound(firewall.HttpResponse) {
			// assume deleted
			return nil
		}

		return fmt.Errorf("making Read request on Azure Firewall %q (Resource Group %q): %+v", id.AzureFirewallName, id.ResourceGroup, err)
	}

	if firewall.Model == nil {
		return fmt.Errorf("retrieving Network Rule Collection %q (Firewall %q / Resource Group %q): `model` was nil", id.NetworkRuleCollectionName, id.AzureFirewallName, id.ResourceGroup)
	}

	props := firewall.Model.Properties
	if props == nil {
		return fmt.Errorf("retrieving Network Rule Collection %q (Firewall %q / Resource Group %q): `props` was nil", id.NetworkRuleCollectionName, id.AzureFirewallName, id.ResourceGroup)
	}
	if props.NetworkRuleCollections == nil {
		return fmt.Errorf("retrieving Network Rule Collection %q (Firewall %q / Resource Group %q): `props.NetworkRuleCollections` was nil", id.NetworkRuleCollectionName, id.AzureFirewallName, id.ResourceGroup)
	}

	networkRules := make([]azurefirewalls.AzureFirewallNetworkRuleCollection, 0)
	for _, rule := range *props.NetworkRuleCollections {
		if rule.Name == nil {
			continue
		}

		if *rule.Name != id.NetworkRuleCollectionName {
			networkRules = append(networkRules, rule)
		}
	}
	props.NetworkRuleCollections = &networkRules

	if err = client.CreateOrUpdateThenPoll(ctx, firewallId, *firewall.Model, azurefirewalls.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("deleting Network Rule Collection %q from Firewall %q (Resource Group %q): %+v", id.NetworkRuleCollectionName, id.AzureFirewallName, id.ResourceGroup, err)
	}

	return nil
}
