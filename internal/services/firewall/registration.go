// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package firewall

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/firewall"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Firewall"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Network",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_firewall":        firewallDataSource(),
		"azurerm_firewall_policy": FirewallDataSourcePolicy(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_firewall_application_rule_collection":  resourceFirewallApplicationRuleCollection(),
		"azurerm_firewall_policy":                       resourceFirewallPolicy(),
		"azurerm_firewall_policy_rule_collection_group": resourceFirewallPolicyRuleCollectionGroup(),
		"azurerm_firewall_nat_rule_collection":          resourceFirewallNatRuleCollection(),
		"azurerm_firewall_network_rule_collection":      resourceFirewallNetworkRuleCollection(),
		"azurerm_firewall":                              resourceFirewall(),
	}
}
