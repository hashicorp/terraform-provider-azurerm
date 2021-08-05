package sql

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "SQL"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Database",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_sql_server":   dataSourceSqlServer(),
		"azurerm_sql_database": dataSourceSqlDatabase(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_sql_active_directory_administrator": resourceSqlAdministrator(),
		"azurerm_sql_database":                       resourceSqlDatabase(),
		"azurerm_sql_elasticpool":                    resourceSqlElasticPool(),
		"azurerm_sql_failover_group":                 resourceSqlFailoverGroup(),
		"azurerm_sql_firewall_rule":                  resourceSqlFirewallRule(),
		"azurerm_sql_server":                         resourceSqlServer(),
		"azurerm_sql_virtual_network_rule":           resourceSqlVirtualNetworkRule(),
	}
}
