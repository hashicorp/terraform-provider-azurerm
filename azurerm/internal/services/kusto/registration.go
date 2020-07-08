package kusto

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Kusto"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Data Explorer",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_kusto_cluster": dataSourceArmKustoCluster(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_kusto_cluster":                      resourceArmKustoCluster(),
		"azurerm_kusto_cluster_principal_assignment": resourceArmKustoClusterPrincipalAssignment(),
		"azurerm_kusto_database":                     resourceArmKustoDatabase(),
		"azurerm_kusto_database_principal":           resourceArmKustoDatabasePrincipal(),
		"azurerm_kusto_eventhub_data_connection":     resourceArmKustoEventHubDataConnection(),
	}
}
