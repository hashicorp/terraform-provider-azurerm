// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/resources/monitor_aad_diagnostic_setting"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/resources/monitor_action_group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/resources/monitor_activity_log_alert"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/resources/monitor_alert_processing_rule_action_group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/resources/monitor_alert_processing_rule_suppression"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/resources/monitor_alert_prometheus_rule_group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/resources/monitor_autoscale_setting"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/resources/monitor_data_collection_endpoint"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/resources/monitor_data_collection_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/resources/monitor_data_collection_rule_association"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/resources/monitor_diagnostic_categories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/resources/monitor_diagnostic_setting"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/resources/monitor_metric_alert"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/resources/monitor_private_link_scope"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/resources/monitor_private_link_scoped_service"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/resources/monitor_scheduled_query_rules_alert"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/resources/monitor_scheduled_query_rules_alert_v2"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/resources/monitor_scheduled_query_rules_log"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/resources/monitor_smart_detector_alert_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/resources/monitor_workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.TypedServiceRegistration                   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/monitor"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		monitor_data_collection_endpoint.DataSource{},
		monitor_data_collection_rule.DataSource{},
		monitor_workspace.DataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		monitor_alert_processing_rule_action_group.Resource{},
		monitor_alert_processing_rule_suppression.Resource{},
		monitor_alert_prometheus_rule_group.Resource{},
		monitor_data_collection_endpoint.Resource{},
		monitor_data_collection_rule_association.Resource{},
		monitor_data_collection_rule.Resource{},
		monitor_scheduled_query_rules_alert_v2.Resource{},
		monitor_workspace.Resource{},
	}
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Monitor"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Monitor",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_monitor_action_group":                monitor_action_group.RegisterDataSource(),
		"azurerm_monitor_diagnostic_categories":       monitor_diagnostic_categories.RegisterDataSource(),
		"azurerm_monitor_scheduled_query_rules_alert": monitor_scheduled_query_rules_alert.RegisterDataSource(),
		"azurerm_monitor_scheduled_query_rules_log":   monitor_scheduled_query_rules_log.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_monitor_aad_diagnostic_setting":      monitor_aad_diagnostic_setting.RegisterResource(),
		"azurerm_monitor_action_group":                monitor_action_group.RegisterResource(),
		"azurerm_monitor_activity_log_alert":          monitor_activity_log_alert.RegisterResource(),
		"azurerm_monitor_autoscale_setting":           monitor_autoscale_setting.RegisterResource(),
		"azurerm_monitor_diagnostic_setting":          monitor_diagnostic_setting.RegisterResource(),
		"azurerm_monitor_metric_alert":                monitor_metric_alert.RegisterResource(),
		"azurerm_monitor_private_link_scope":          monitor_private_link_scope.RegisterResource(),
		"azurerm_monitor_private_link_scoped_service": monitor_private_link_scoped_service.RegisterResource(),
		"azurerm_monitor_scheduled_query_rules_alert": monitor_scheduled_query_rules_alert.RegisterResource(),
		"azurerm_monitor_scheduled_query_rules_log":   monitor_scheduled_query_rules_log.RegisterResource(),
		"azurerm_monitor_smart_detector_alert_rule":   monitor_smart_detector_alert_rule.RegisterResource(),
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
		monitor_action_group.ListResource{},
		monitor_activity_log_alert.ListResource{},
		monitor_metric_alert.ListResource{},
		monitor_scheduled_query_rules_alert.ListResource{},
		monitor_workspace.ListResource{},
		monitor_scheduled_query_rules_alert_v2.ListResource{},
	}
}
