package datalake

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/data-lake"
}

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
		"azurerm_data_lake_store": dataSourceDataLakeStoreAccount(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	out := map[string]*pluginsdk.Resource{
		"azurerm_data_lake_analytics_account":       resourceDataLakeAnalyticsAccount(),
		"azurerm_data_lake_analytics_firewall_rule": resourceDataLakeAnalyticsFirewallRule(),
	}

	if !features.ThreePointOhBeta() {
		out["azurerm_data_lake_store_file"] = resourceDataLakeStoreFile()
		out["azurerm_data_lake_store_firewall_rule"] = resourceDataLakeStoreFirewallRule()
		out["azurerm_data_lake_store"] = resourceDataLakeStore()
		out["azurerm_data_lake_store_virtual_network_rule"] = resourceDataLakeStoreVirtualNetworkRule()
	}

	return out
}
