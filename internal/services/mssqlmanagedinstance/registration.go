// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssqlmanagedinstance

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/resources/mssql_managed_database"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/resources/mssql_managed_instance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/resources/mssql_managed_instance_active_directory_administrator"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/resources/mssql_managed_instance_failover_group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/resources/mssql_managed_instance_security_alert_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/resources/mssql_managed_instance_start_stop_schedule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/resources/mssql_managed_instance_transparent_data_encryption"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/resources/mssql_managed_instance_vulnerability_assessment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.TypedServiceRegistration                   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/mssqlmanagedinstance"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Microsoft SQL Server Managed Instances"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Database",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_mssql_managed_instance_security_alert_policy":       mssql_managed_instance_security_alert_policy.RegisterResource(),
		"azurerm_mssql_managed_instance_transparent_data_encryption": mssql_managed_instance_transparent_data_encryption.RegisterResource(),
		"azurerm_mssql_managed_instance_vulnerability_assessment":    mssql_managed_instance_vulnerability_assessment.RegisterResource(),
	}
}

// DataSources returns the typed DataSources supported by this service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		mssql_managed_database.DataSource{},
		mssql_managed_instance.DataSource{},
	}
}

// Resources returns the typed Resources supported by this service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		mssql_managed_database.Resource{},
		mssql_managed_instance_active_directory_administrator.Resource{},
		mssql_managed_instance_failover_group.Resource{},
		mssql_managed_instance.Resource{},
		mssql_managed_instance_start_stop_schedule.Resource{},
	}
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
