package datalake

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
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
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_data_lake_store":                   dataSourceDataLakeStoreAccount(),
		"azurerm_storage_data_lake_gen1_filesystem": dataSourceStorageDataLakeGen1Filesystem(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_data_lake_analytics_account":       resourceDataLakeAnalyticsAccount(),
		"azurerm_data_lake_analytics_firewall_rule": resourceDataLakeAnalyticsFirewallRule(),

		// TODO Remove in 3.0 as per https://github.com/hashicorp/terraform-provider-azurerm/issues/13910
		"azurerm_data_lake_store":                      resourceDataLakeStore(),
		"azurerm_data_lake_store_file":                 resourceDataLakeStoreFile(),
		"azurerm_data_lake_store_firewall_rule":        resourceDataLakeStoreFirewallRule(),
		"azurerm_data_lake_store_virtual_network_rule": resourceDataLakeStoreVirtualNetworkRule(),

		"azurerm_storage_data_lake_gen1_filesystem":           resourceStorageDataLakeGen1Filesystem(),
		"azurerm_storage_data_lake_gen1_file":                 resourceStorageDataLakeGen1File(),
		"azurerm_storage_data_lake_gen1_firewall_rule":        resourceStorageDataLakeGen1FirewallRule(),
		"azurerm_storage_data_lake_gen1_virtual_network_rule": resourceStorageDataLakeGen1VirtualNetworkRule(),
	}
}
