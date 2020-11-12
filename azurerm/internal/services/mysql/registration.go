package mysql

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_mysql_server": dataSourceArmMySqlServer(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_mysql_configuration":                  resourceArmMySQLConfiguration(),
		"azurerm_mysql_database":                       resourceArmMySqlDatabase(),
		"azurerm_mysql_firewall_rule":                  resourceArmMySqlFirewallRule(),
		"azurerm_mysql_server":                         resourceArmMySqlServer(),
		"azurerm_mysql_server_key":                     resourceArmMySQLServerKey(),
		"azurerm_mysql_virtual_network_rule":           resourceArmMySQLVirtualNetworkRule(),
		"azurerm_mysql_active_directory_administrator": resourceArmMySQLAdministrator()}
}
