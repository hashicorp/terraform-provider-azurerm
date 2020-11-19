package datalake

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Data Lake"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Data Lake",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_data_lake_store": dataSourceArmDataLakeStoreAccount(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_data_lake_analytics_account":       resourceArmDataLakeAnalyticsAccount(),
		"azurerm_data_lake_analytics_firewall_rule": resourceArmDataLakeAnalyticsFirewallRule(),
		"azurerm_data_lake_store_file":              resourceArmDataLakeStoreFile(),
		"azurerm_data_lake_store_firewall_rule":     resourceArmDataLakeStoreFirewallRule(),
		"azurerm_data_lake_store":                   resourceArmDataLakeStore(),
	}
}
