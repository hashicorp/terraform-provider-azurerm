package firewall

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

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
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_firewall":        dataSourceArmFirewall(),
		"azurerm_firewall_policy": dataSourceArmFirewallPolicy(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_firewall_application_rule_collection":  resourceArmFirewallApplicationRuleCollection(),
		"azurerm_firewall_policy":                       resourceArmFirewallPolicy(),
		"azurerm_firewall_policy_rule_collection_group": resourceArmFirewallPolicyRuleCollectionGroup(),
		"azurerm_firewall_nat_rule_collection":          resourceArmFirewallNatRuleCollection(),
		"azurerm_firewall_network_rule_collection":      resourceArmFirewallNetworkRuleCollection(),
		"azurerm_firewall":                              resourceArmFirewall(),
	}
}
