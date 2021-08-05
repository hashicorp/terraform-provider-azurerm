package monitor

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

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
		"azurerm_monitor_action_group":                dataSourceMonitorActionGroup(),
		"azurerm_monitor_diagnostic_categories":       dataSourceMonitorDiagnosticCategories(),
		"azurerm_monitor_log_profile":                 dataSourceMonitorLogProfile(),
		"azurerm_monitor_scheduled_query_rules_alert": dataSourceMonitorScheduledQueryRulesAlert(),
		"azurerm_monitor_scheduled_query_rules_log":   dataSourceMonitorScheduledQueryRulesLog(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_monitor_aad_diagnostic_setting":      resourceMonitorAADDiagnosticSetting(),
		"azurerm_monitor_autoscale_setting":           resourceMonitorAutoScaleSetting(),
		"azurerm_monitor_action_group":                resourceMonitorActionGroup(),
		"azurerm_monitor_action_rule_action_group":    resourceMonitorActionRuleActionGroup(),
		"azurerm_monitor_action_rule_suppression":     resourceMonitorActionRuleSuppression(),
		"azurerm_monitor_activity_log_alert":          resourceMonitorActivityLogAlert(),
		"azurerm_monitor_diagnostic_setting":          resourceMonitorDiagnosticSetting(),
		"azurerm_monitor_log_profile":                 resourceMonitorLogProfile(),
		"azurerm_monitor_metric_alert":                resourceMonitorMetricAlert(),
		"azurerm_monitor_scheduled_query_rules_alert": resourceMonitorScheduledQueryRulesAlert(),
		"azurerm_monitor_scheduled_query_rules_log":   resourceMonitorScheduledQueryRulesLog(),
		"azurerm_monitor_smart_detector_alert_rule":   resourceMonitorSmartDetectorAlertRule(),
	}
}
