// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/resources/log_analytics_cluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/resources/log_analytics_cluster_customer_managed_key"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/resources/log_analytics_data_export_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/resources/log_analytics_datasource_windows_event"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/resources/log_analytics_datasource_windows_performance_counter"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/resources/log_analytics_linked_service"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/resources/log_analytics_linked_storage_account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/resources/log_analytics_query_pack"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/resources/log_analytics_query_pack_query"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/resources/log_analytics_saved_search"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/resources/log_analytics_solution"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/resources/log_analytics_storage_insights"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/resources/log_analytics_workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/resources/log_analytics_workspace_table"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/resources/log_analytics_workspace_table_custom_log"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/resources/log_analytics_workspace_table_microsoft"
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
		log_analytics_workspace_table.DataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		log_analytics_cluster.Resource{},
		log_analytics_query_pack_query.Resource{},
		log_analytics_query_pack.Resource{},
		log_analytics_solution.Resource{},
		log_analytics_workspace_table.Resource{},
		log_analytics_workspace_table_custom_log.Resource{},
		log_analytics_workspace_table_microsoft.Resource{},
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
		"azurerm_log_analytics_workspace": log_analytics_workspace.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_log_analytics_cluster_customer_managed_key":           log_analytics_cluster_customer_managed_key.RegisterResource(),
		"azurerm_log_analytics_data_export_rule":                       log_analytics_data_export_rule.RegisterResource(),
		"azurerm_log_analytics_datasource_windows_event":               log_analytics_datasource_windows_event.RegisterResource(),
		"azurerm_log_analytics_datasource_windows_performance_counter": log_analytics_datasource_windows_performance_counter.RegisterResource(),
		"azurerm_log_analytics_linked_service":                         log_analytics_linked_service.RegisterResource(),
		"azurerm_log_analytics_linked_storage_account":                 log_analytics_linked_storage_account.RegisterResource(),
		"azurerm_log_analytics_saved_search":                           log_analytics_saved_search.RegisterResource(),
		"azurerm_log_analytics_storage_insights":                       log_analytics_storage_insights.RegisterResource(),
		"azurerm_log_analytics_workspace":                              log_analytics_workspace.RegisterResource(),
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
	return []sdk.FrameworkListWrappedResource{
		log_analytics_workspace.ListResource{},
	}
}
