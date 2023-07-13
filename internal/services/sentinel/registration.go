// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/sentinel"
}

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
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_sentinel_alert_rule":          dataSourceSentinelAlertRule(),
		"azurerm_sentinel_alert_rule_template": dataSourceSentinelAlertRuleTemplate(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_sentinel_alert_rule_fusion":                                            resourceSentinelAlertRuleFusion(),
		"azurerm_sentinel_alert_rule_machine_learning_behavior_analytics":               resourceSentinelAlertRuleMLBehaviorAnalytics(),
		"azurerm_sentinel_alert_rule_ms_security_incident":                              resourceSentinelAlertRuleMsSecurityIncident(),
		"azurerm_sentinel_alert_rule_scheduled":                                         resourceSentinelAlertRuleScheduled(),
		"azurerm_sentinel_alert_rule_nrt":                                               resourceSentinelAlertRuleNrt(),
		"azurerm_sentinel_data_connector_aws_cloud_trail":                               resourceSentinelDataConnectorAwsCloudTrail(),
		"azurerm_sentinel_data_connector_azure_active_directory":                        resourceSentinelDataConnectorAzureActiveDirectory(),
		"azurerm_sentinel_data_connector_azure_advanced_threat_protection":              resourceSentinelDataConnectorAzureAdvancedThreatProtection(),
		"azurerm_sentinel_data_connector_azure_security_center":                         resourceSentinelDataConnectorAzureSecurityCenter(),
		"azurerm_sentinel_data_connector_microsoft_cloud_app_security":                  resourceSentinelDataConnectorMicrosoftCloudAppSecurity(),
		"azurerm_sentinel_data_connector_microsoft_defender_advanced_threat_protection": resourceSentinelDataConnectorMicrosoftDefenderAdvancedThreatProtection(),
		"azurerm_sentinel_data_connector_office_365":                                    resourceSentinelDataConnectorOffice365(),
		"azurerm_sentinel_data_connector_office_atp":                                    resourceSentinelDataConnectorOfficeATP(),
		"azurerm_sentinel_data_connector_threat_intelligence":                           resourceSentinelDataConnectorThreatIntelligence(),
		"azurerm_sentinel_automation_rule":                                              resourceSentinelAutomationRule(),
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		AlertRuleAnomalyDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		AlertRuleThreatIntelligenceResource{},
		WatchlistResource{},
		WatchlistItemResource{},
		DataConnectorAwsS3Resource{},
		DataConnectorMicrosoftThreatProtectionResource{},
		DataConnectorIOTResource{},
		DataConnectorDynamics365Resource{},
		DataConnectorOffice365ProjectResource{},
		DataConnectorOfficePowerBIResource{},
		DataConnectorOfficeIRMResource{},
		LogAnalyticsWorkspaceOnboardResource{},
		DataConnectorThreatIntelligenceTAXIIResource{},
		DataConnectorMicrosoftThreatIntelligenceResource{},
		AlertRuleAnomalyBuiltInResource{},
		MetadataResource{},
		AlertRuleAnomalyDuplicateResource{},
		ThreatIntelligenceIndicator{},
	}
}
