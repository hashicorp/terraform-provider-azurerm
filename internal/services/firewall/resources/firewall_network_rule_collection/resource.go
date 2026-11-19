// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package firewall_network_rule_collection

import (
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var AzureFirewallResourceName = "azurerm_firewall"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFirewallNetworkRuleCollectionCreateUpdate,
		Read:   resourceFirewallNetworkRuleCollectionRead,
		Update: resourceFirewallNetworkRuleCollectionCreateUpdate,
		Delete: resourceFirewallNetworkRuleCollectionDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FirewallNetworkRuleCollectionID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: firewallNetworkRuleCollectionSchema(),
	}
}
