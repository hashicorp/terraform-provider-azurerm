package monitor

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Monitor"
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_monitor_action_group":          dataSourceArmMonitorActionGroup(),
		"azurerm_monitor_diagnostic_categories": dataSourceArmMonitorDiagnosticCategories(),
		"azurerm_monitor_log_profile":           dataSourceArmMonitorLogProfile()}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_autoscale_setting":          resourceArmAutoScaleSetting(),
		"azurerm_metric_alertrule":           resourceArmMetricAlertRule(),
		"azurerm_monitor_autoscale_setting":  resourceArmMonitorAutoScaleSetting(),
		"azurerm_monitor_action_group":       resourceArmMonitorActionGroup(),
		"azurerm_monitor_activity_log_alert": resourceArmMonitorActivityLogAlert(),
		"azurerm_monitor_diagnostic_setting": resourceArmMonitorDiagnosticSetting(),
		"azurerm_monitor_log_profile":        resourceArmMonitorLogProfile(),
		"azurerm_monitor_metric_alert":       resourceArmMonitorMetricAlert(),
		"azurerm_monitor_metric_alertrule":   resourceArmMonitorMetricAlertRule()}
}
