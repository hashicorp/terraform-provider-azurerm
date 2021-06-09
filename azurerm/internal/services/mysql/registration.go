package mysql

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "MySQL"
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
		"azurerm_mysql_server": dataSourceMySqlServer(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_mysql_configuration":                  resourceMySQLConfiguration(),
		"azurerm_mysql_database":                       resourceMySqlDatabase(),
		"azurerm_mysql_firewall_rule":                  resourceMySqlFirewallRule(),
		"azurerm_mysql_server":                         resourceMySqlServer(),
		"azurerm_mysql_server_key":                     resourceMySQLServerKey(),
		"azurerm_mysql_virtual_network_rule":           resourceMySQLVirtualNetworkRule(),
		"azurerm_mysql_active_directory_administrator": resourceMySQLAdministrator(),
	}
}
