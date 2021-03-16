package postgres

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "PostgreSQL"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Database",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_postgresql_server": dataSourcePostgreSqlServer(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_postgresql_configuration":                  resourcePostgreSQLConfiguration(),
		"azurerm_postgresql_database":                       resourcePostgreSQLDatabase(),
		"azurerm_postgresql_firewall_rule":                  resourcePostgreSQLFirewallRule(),
		"azurerm_postgresql_server":                         resourcePostgreSQLServer(),
		"azurerm_postgresql_server_key":                     resourcePostgreSQLServerKey(),
		"azurerm_postgresql_virtual_network_rule":           resourcePostgreSQLVirtualNetworkRule(),
		"azurerm_postgresql_active_directory_administrator": resourcePostgreSQLAdministrator(),
	}
}
