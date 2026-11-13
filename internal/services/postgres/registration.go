// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgres

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/resources/postgresql_flexible_server"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/resources/postgresql_flexible_server_active_directory_administrator"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/resources/postgresql_flexible_server_backup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/resources/postgresql_flexible_server_configuration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/resources/postgresql_flexible_server_database"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/resources/postgresql_flexible_server_firewall_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/resources/postgresql_flexible_server_virtual_endpoint"
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
		"azurerm_postgresql_flexible_server": postgresql_flexible_server.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_postgresql_flexible_server":                                postgresql_flexible_server.RegisterResource(),
		"azurerm_postgresql_flexible_server_active_directory_administrator": postgresql_flexible_server_active_directory_administrator.RegisterResource(),
		"azurerm_postgresql_flexible_server_configuration":                  postgresql_flexible_server_configuration.RegisterResource(),
		"azurerm_postgresql_flexible_server_database":                       postgresql_flexible_server_database.RegisterResource(),
		"azurerm_postgresql_flexible_server_firewall_rule":                  postgresql_flexible_server_firewall_rule.RegisterResource(),
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		postgresql_flexible_server_backup.Resource{},
		postgresql_flexible_server_virtual_endpoint.Resource{},
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
	return []sdk.FrameworkListWrappedResource{
		postgresql_flexible_server_configuration.ListResource{},
		postgresql_flexible_server_firewall_rule.ListResource{},
	}
}
