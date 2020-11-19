package mariadb

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "MariaDB"
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
		"azurerm_mariadb_server": dataSourceMariaDbServer(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_mariadb_configuration":        resourceArmMariaDbConfiguration(),
		"azurerm_mariadb_database":             resourceArmMariaDbDatabase(),
		"azurerm_mariadb_firewall_rule":        resourceArmMariaDBFirewallRule(),
		"azurerm_mariadb_server":               resourceArmMariaDbServer(),
		"azurerm_mariadb_virtual_network_rule": resourceArmMariaDbVirtualNetworkRule(),
	}
}
