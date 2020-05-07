package sql

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_sql_server":   dataSourceSqlServer(),
		"azurerm_sql_database": dataSourceSqlDatabase(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_sql_active_directory_administrator": resourceArmSqlAdministrator(),
		"azurerm_sql_database":                       resourceArmSqlDatabase(),
		"azurerm_sql_elasticpool":                    resourceArmSqlElasticPool(),
		"azurerm_sql_failover_group":                 resourceArmSqlFailoverGroup(),
		"azurerm_sql_firewall_rule":                  resourceArmSqlFirewallRule(),
		"azurerm_sql_server":                         resourceArmSqlServer(),
		"azurerm_sql_virtual_network_rule":           resourceArmSqlVirtualNetworkRule(),
	}
}
