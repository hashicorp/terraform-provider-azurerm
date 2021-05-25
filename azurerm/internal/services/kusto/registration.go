package kusto

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
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
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_kusto_cluster": dataSourceKustoCluster(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_kusto_cluster":                         resourceKustoCluster(),
		"azurerm_kusto_cluster_customer_managed_key":    resourceKustoClusterCustomerManagedKey(),
		"azurerm_kusto_cluster_principal_assignment":    resourceKustoClusterPrincipalAssignment(),
		"azurerm_kusto_database":                        resourceKustoDatabase(),
		"azurerm_kusto_database_principal":              resourceKustoDatabasePrincipal(),
		"azurerm_kusto_database_principal_assignment":   resourceKustoDatabasePrincipalAssignment(),
		"azurerm_kusto_eventgrid_data_connection":       resourceKustoEventGridDataConnection(),
		"azurerm_kusto_eventhub_data_connection":        resourceKustoEventHubDataConnection(),
		"azurerm_kusto_iothub_data_connection":          resourceKustoIotHubDataConnection(),
		"azurerm_kusto_attached_database_configuration": resourceKustoAttachedDatabaseConfiguration(),
	}
}
