// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package postgres

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
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
	return map[string]*pluginsdk.Resource{
		"azurerm_postgresql_flexible_server": dataSourcePostgresqlFlexibleServer(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_postgresql_flexible_server":                                resourcePostgresqlFlexibleServer(),
		"azurerm_postgresql_flexible_server_active_directory_administrator": resourcePostgresqlFlexibleServerAdministrator(),
		"azurerm_postgresql_flexible_server_configuration":                  resourcePostgresqlFlexibleServerConfiguration(),
		"azurerm_postgresql_flexible_server_database":                       resourcePostgresqlFlexibleServerDatabase(),
		"azurerm_postgresql_flexible_server_firewall_rule":                  resourcePostgresqlFlexibleServerFirewallRule(),
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		PostgresqlFlexibleServerBackupResource{},
		PostgresqlFlexibleServerVirtualEndpointResource{},
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Actions() []func() action.Action {
	return []func() action.Action{}
}

func (r Registration) FrameworkResources() []sdk.FrameworkWrappedResource {
	return []sdk.FrameworkWrappedResource{}
}

func (r Registration) FrameworkDataSources() []sdk.FrameworkWrappedDataSource {
	return []sdk.FrameworkWrappedDataSource{}
}

func (r Registration) EphemeralResources() []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (r Registration) ListResources() []sdk.FrameworkListWrappedResource {
	return []sdk.FrameworkListWrappedResource{}
}
