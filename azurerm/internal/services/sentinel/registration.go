package sentinel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Sentinel"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Sentinel",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_sentinel_alert_rule":          dataSourceSentinelAlertRule(),
		"azurerm_sentinel_alert_rule_template": dataSourceSentinelAlertRuleTemplate(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_sentinel_alert_rule_fusion":                               resourceSentinelAlertRuleFusion(),
		"azurerm_sentinel_alert_rule_ms_security_incident":                 resourceSentinelAlertRuleMsSecurityIncident(),
		"azurerm_sentinel_alert_rule_scheduled":                            resourceSentinelAlertRuleScheduled(),
		"azurerm_sentinel_data_connector_azure_active_directory":           resourceSentinelDataConnectorAzureActiveDirectory(),
		"azurerm_sentinel_data_connector_azure_advanced_threat_protection": resourceSentinelDataConnectorAzureAdvancedThreatProtection(),
		"azurerm_sentinel_data_connector_azure_security_center":            resourceSentinelDataConnectorAzureSecurityCenter(),
		"azurerm_sentinel_data_connector_microsoft_cloud_app_security":     resourceSentinelDataConnectorMicrosoftCloudAppSecurity(),
		"azurerm_sentinel_data_connector_threat_intelligence":              resourceSentinelDataConnectorThreatIntelligence(),
		"azurerm_sentinel_data_connector_office_365":                       resourceSentinelDataConnectorOffice365(),
		"azurerm_sentinel_data_connector_office_atp":                       resourceSentinelDataConnectorOfficeATP(),
	}
}
