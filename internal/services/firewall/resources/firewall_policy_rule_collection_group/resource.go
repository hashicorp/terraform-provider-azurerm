// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package firewall_policy_rule_collection_group

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/firewallpolicyrulecollectiongroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name firewall_policy_rule_collection_group -properties "name" -compare-values "subscription_id:firewall_policy_id,resource_group_name:firewall_policy_id,firewall_policy_name:firewall_policy_id"

const AzureFirewallPolicyResourceName = "azurerm_firewall_policy"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFirewallPolicyRuleCollectionGroupCreateUpdate,
		Read:   resourceFirewallPolicyRuleCollectionGroupRead,
		Update: resourceFirewallPolicyRuleCollectionGroupCreateUpdate,
		Delete: resourceFirewallPolicyRuleCollectionGroupDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&firewallpolicyrulecollectiongroups.RuleCollectionGroupId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&firewallpolicyrulecollectiongroups.RuleCollectionGroupId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: firewallPolicyRuleCollectionGroupSchema(),
	}
}
