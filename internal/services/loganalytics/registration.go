// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.TypedServiceRegistration                   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/log-analytics"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		LogAnalyticsWorkspaceTableDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		LogAnalyticsClusterResource{},
		LogAnalyticsQueryPackResource{},
		LogAnalyticsQueryPackQueryResource{},
		LogAnalyticsSolutionResource{},
		LogAnalyticsWorkspaceTableResource{},
		WorkspaceTableCustomLogResource{},
	}
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Log Analytics"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Log Analytics",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_log_analytics_workspace": dataSourceLogAnalyticsWorkspace(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_log_analytics_cluster_customer_managed_key":           resourceLogAnalyticsClusterCustomerManagedKey(),
		"azurerm_log_analytics_datasource_windows_event":               resourceLogAnalyticsDataSourceWindowsEvent(),
		"azurerm_log_analytics_datasource_windows_performance_counter": resourceLogAnalyticsDataSourceWindowsPerformanceCounter(),
		"azurerm_log_analytics_data_export_rule":                       resourceLogAnalyticsDataExport(),
		"azurerm_log_analytics_linked_service":                         resourceLogAnalyticsLinkedService(),
		"azurerm_log_analytics_linked_storage_account":                 resourceLogAnalyticsLinkedStorageAccount(),
		"azurerm_log_analytics_saved_search":                           resourceLogAnalyticsSavedSearch(),
		"azurerm_log_analytics_storage_insights":                       resourceLogAnalyticsStorageInsights(),
		"azurerm_log_analytics_workspace":                              resourceLogAnalyticsWorkspace(),
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
