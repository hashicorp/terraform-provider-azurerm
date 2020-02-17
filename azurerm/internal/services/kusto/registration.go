package kusto

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Kusto"
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_kusto_cluster":                  resourceArmKustoCluster(),
		"azurerm_kusto_database":                 resourceArmKustoDatabase(),
		"azurerm_kusto_database_principal":       resourceArmKustoDatabasePrincipal(),
		"azurerm_kusto_eventhub_data_connection": resourceArmKustoEventHubDataConnection(),
	}
}
