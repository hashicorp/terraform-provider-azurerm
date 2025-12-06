// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package postgres

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.TypedServiceRegistrationWithAGitHubLabel   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/postgresql"
}

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
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	dataSources := map[string]*pluginsdk.Resource{
		"azurerm_postgresql_flexible_server": dataSourcePostgresqlFlexibleServer(),
	}

	if !features.FivePointOh() {
		dataSources["azurerm_postgresql_server"] = dataSourcePostgreSqlServer()
	}

	return dataSources
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	resources := map[string]*pluginsdk.Resource{
		"azurerm_postgresql_flexible_server":                                resourcePostgresqlFlexibleServer(),
		"azurerm_postgresql_flexible_server_firewall_rule":                  resourcePostgresqlFlexibleServerFirewallRule(),
		"azurerm_postgresql_flexible_server_configuration":                  resourcePostgresqlFlexibleServerConfiguration(),
		"azurerm_postgresql_flexible_server_database":                       resourcePostgresqlFlexibleServerDatabase(),
		"azurerm_postgresql_flexible_server_active_directory_administrator": resourcePostgresqlFlexibleServerAdministrator(),
	}

	if !features.FivePointOh() {
		resources["azurerm_postgresql_server"] = resourcePostgreSQLServer()
		resources["azurerm_postgresql_active_directory_administrator"] = resourcePostgreSQLAdministrator()
		resources["azurerm_postgresql_configuration"] = resourcePostgreSQLConfiguration()
		resources["azurerm_postgresql_database"] = resourcePostgreSQLDatabase()
		resources["azurerm_postgresql_firewall_rule"] = resourcePostgreSQLFirewallRule()
		resources["azurerm_postgresql_server_key"] = resourcePostgreSQLServerKey()
		resources["azurerm_postgresql_virtual_network_rule"] = resourcePostgreSQLVirtualNetworkRule()
	}

	return resources
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		PostgresqlFlexibleServerBackupResource{},
		PostgresqlFlexibleServerVirtualEndpointResource{},
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		PostgresqlFlexibleServerActiveDirectoryAdministratorDataSource{},
	}
}
