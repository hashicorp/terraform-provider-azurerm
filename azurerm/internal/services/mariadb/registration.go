package mariadb

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
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
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_mariadb_server": dataSourceMariaDbServer(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_mariadb_configuration":        resourceMariaDbConfiguration(),
		"azurerm_mariadb_database":             resourceMariaDbDatabase(),
		"azurerm_mariadb_firewall_rule":        resourceArmMariaDBFirewallRule(),
		"azurerm_mariadb_server":               resourceMariaDbServer(),
		"azurerm_mariadb_virtual_network_rule": resourceMariaDbVirtualNetworkRule(),
	}
}
