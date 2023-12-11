// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sql

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/sql"
}

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
		"azurerm_sql_server":           dataSourceSqlServer(),
		"azurerm_sql_database":         dataSourceSqlDatabase(),
		"azurerm_sql_managed_instance": dataSourceArmSqlMiServer(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_sql_active_directory_administrator":                  resourceSqlAdministrator(),
		"azurerm_sql_database":                                        resourceSqlDatabase(),
		"azurerm_sql_elasticpool":                                     resourceSqlElasticPool(),
		"azurerm_sql_failover_group":                                  resourceSqlFailoverGroup(),
		"azurerm_sql_firewall_rule":                                   resourceSqlFirewallRule(),
		"azurerm_sql_managed_database":                                resourceArmSqlManagedDatabase(),
		"azurerm_sql_managed_instance":                                resourceArmSqlMiServer(),
		"azurerm_sql_managed_instance_failover_group":                 resourceSqlInstanceFailoverGroup(),
		"azurerm_sql_managed_instance_active_directory_administrator": resourceSqlManagedInstanceAdministrator(),
		"azurerm_sql_server":                                          resourceSqlServer(),
		"azurerm_sql_virtual_network_rule":                            resourceSqlVirtualNetworkRule(),
	}
}
