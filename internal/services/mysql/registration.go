// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/resources/mysql_flexible_database"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/resources/mysql_flexible_server"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/resources/mysql_flexible_server_active_directory_administrator"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/resources/mysql_flexible_server_configuration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/resources/mysql_flexible_server_firewall_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.TypedServiceRegistration                   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
	_ sdk.FrameworkServiceRegistration               = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/mysql"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		mysql_flexible_server_active_directory_administrator.Resource{},
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
		"azurerm_mysql_flexible_server": mysql_flexible_server.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_mysql_flexible_database":             mysql_flexible_database.RegisterResource(),
		"azurerm_mysql_flexible_server":               mysql_flexible_server.RegisterResource(),
		"azurerm_mysql_flexible_server_configuration": mysql_flexible_server_configuration.RegisterResource(),
		"azurerm_mysql_flexible_server_firewall_rule": mysql_flexible_server_firewall_rule.RegisterResource(),
	}
}

// Actions implements [sdk.FrameworkServiceRegistration].
func (r Registration) Actions() []func() action.Action {
	return []func() action.Action{}
}

// EphemeralResources implements [sdk.FrameworkServiceRegistration].
func (r Registration) EphemeralResources() []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

// FrameworkDataSources implements [sdk.FrameworkServiceRegistration].
func (r Registration) FrameworkDataSources() []sdk.FrameworkWrappedDataSource {
	return []sdk.FrameworkWrappedDataSource{}
}

// FrameworkResources implements [sdk.FrameworkServiceRegistration].
func (r Registration) FrameworkResources() []sdk.FrameworkWrappedResource {
	return []sdk.FrameworkWrappedResource{}
}

// ListResources returns a list of List Resources supported by this Service
func (r Registration) ListResources() []sdk.FrameworkListWrappedResource {
	return []sdk.FrameworkListWrappedResource{
		mysql_flexible_database.ListResource{},
		mysql_flexible_server_configuration.ListResource{},
		mysql_flexible_server_firewall_rule.ListResource{},
		mysql_flexible_server.ListResource{},
	}
}
