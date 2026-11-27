// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_firewall_rule

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/ipfirewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/resources/synapse_firewall_rule/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSynapseFirewallRuleCreateUpdate,
		Read:   resourceSynapseFirewallRuleRead,
		Update: resourceSynapseFirewallRuleCreateUpdate,
		Delete: resourceSynapseFirewallRuleDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := ipfirewallrules.ParseFirewallRuleID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SynapseFirewallRuleV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: synapseFirewallRuleSchema(),
	}
}
