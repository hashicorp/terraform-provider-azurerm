// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

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
	return "service/monitor"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		DataCollectionEndpointDataSource{},
		DataCollectionRuleDataSource{},
		WorkspaceDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		AlertProcessingRuleActionGroupResource{},
		AlertProcessingRuleSuppressionResource{},
		DataCollectionEndpointResource{},
		DataCollectionRuleAssociationResource{},
		DataCollectionRuleResource{},
		ScheduledQueryRulesAlertV2Resource{},
		AlertPrometheusRuleGroupResource{},
		WorkspaceResource{},
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
	dataSources := map[string]*pluginsdk.Resource{
		"azurerm_monitor_action_group":                dataSourceMonitorActionGroup(),
		"azurerm_monitor_diagnostic_categories":       dataSourceMonitorDiagnosticCategories(),
		"azurerm_monitor_scheduled_query_rules_alert": dataSourceMonitorScheduledQueryRulesAlert(),
		"azurerm_monitor_scheduled_query_rules_log":   dataSourceMonitorScheduledQueryRulesLog(),
	}

	return dataSources
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	resources := map[string]*pluginsdk.Resource{
		"azurerm_monitor_aad_diagnostic_setting":      resourceMonitorAADDiagnosticSetting(),
		"azurerm_monitor_autoscale_setting":           resourceMonitorAutoScaleSetting(),
		"azurerm_monitor_action_group":                resourceMonitorActionGroup(),
		"azurerm_monitor_activity_log_alert":          resourceMonitorActivityLogAlert(),
		"azurerm_monitor_diagnostic_setting":          resourceMonitorDiagnosticSetting(),
		"azurerm_monitor_metric_alert":                resourceMonitorMetricAlert(),
		"azurerm_monitor_private_link_scope":          resourceMonitorPrivateLinkScope(),
		"azurerm_monitor_private_link_scoped_service": resourceMonitorPrivateLinkScopedService(),
		"azurerm_monitor_scheduled_query_rules_alert": resourceMonitorScheduledQueryRulesAlert(),
		"azurerm_monitor_scheduled_query_rules_log":   resourceMonitorScheduledQueryRulesLog(),
		"azurerm_monitor_smart_detector_alert_rule":   resourceMonitorSmartDetectorAlertRule(),
	}

	return resources
}
