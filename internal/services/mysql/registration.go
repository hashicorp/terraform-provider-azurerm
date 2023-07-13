// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mysql

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.TypedServiceRegistration                   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/mysql"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		MySQLFlexibleServerAdministratorResource{},
	}
}

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
		"azurerm_mysql_server":          dataSourceMySqlServer(),
		"azurerm_mysql_flexible_server": dataSourceMysqlFlexibleServer(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_mysql_configuration":                  resourceMySQLConfiguration(),
		"azurerm_mysql_database":                       resourceMySqlDatabase(),
		"azurerm_mysql_firewall_rule":                  resourceMySqlFirewallRule(),
		"azurerm_mysql_flexible_server":                resourceMysqlFlexibleServer(),
		"azurerm_mysql_flexible_database":              resourceMySqlFlexibleDatabase(),
		"azurerm_mysql_flexible_server_configuration":  resourceMySQLFlexibleServerConfiguration(),
		"azurerm_mysql_flexible_server_firewall_rule":  resourceMySqlFlexibleServerFirewallRule(),
		"azurerm_mysql_server":                         resourceMySqlServer(),
		"azurerm_mysql_server_key":                     resourceMySQLServerKey(),
		"azurerm_mysql_virtual_network_rule":           resourceMySQLVirtualNetworkRule(),
		"azurerm_mysql_active_directory_administrator": resourceMySQLAdministrator(),
	}
}
