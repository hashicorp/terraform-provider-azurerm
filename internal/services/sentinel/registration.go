// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_alert_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_alert_rule_anomaly"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_alert_rule_anomaly_built_in"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_alert_rule_anomaly_duplicate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_alert_rule_fusion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_alert_rule_machine_learning_behavior_analytics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_alert_rule_ms_security_incident"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_alert_rule_nrt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_alert_rule_scheduled"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_alert_rule_template"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_alert_rule_threat_intelligence"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_automation_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_data_connector_aws_cloud_trail"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_data_connector_aws_s3"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_data_connector_azure_active_directory"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_data_connector_azure_advanced_threat_protection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_data_connector_azure_security_center"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_data_connector_dynamics_365"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_data_connector_iot"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_data_connector_microsoft_cloud_app_security"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_data_connector_microsoft_defender_advanced_threat_protection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_data_connector_microsoft_threat_intelligence"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_data_connector_microsoft_threat_protection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_data_connector_office_365"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_data_connector_office_365_project"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_data_connector_office_atp"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_data_connector_office_irm"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_data_connector_office_power_bi"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_data_connector_threat_intelligence"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_data_connector_threat_intelligence_taxii"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_log_analytics_workspace_onboarding"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_metadata"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_threat_intelligence_indicator"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_watchlist"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/resources/sentinel_watchlist_item"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
	_ sdk.TypedServiceRegistration                   = Registration{}
)

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
		"azurerm_sentinel_alert_rule":          sentinel_alert_rule.RegisterDataSource(),
		"azurerm_sentinel_alert_rule_template": sentinel_alert_rule_template.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_sentinel_alert_rule_fusion":                                            sentinel_alert_rule_fusion.RegisterResource(),
		"azurerm_sentinel_alert_rule_machine_learning_behavior_analytics":               sentinel_alert_rule_machine_learning_behavior_analytics.RegisterResource(),
		"azurerm_sentinel_alert_rule_ms_security_incident":                              sentinel_alert_rule_ms_security_incident.RegisterResource(),
		"azurerm_sentinel_alert_rule_nrt":                                               sentinel_alert_rule_nrt.RegisterResource(),
		"azurerm_sentinel_alert_rule_scheduled":                                         sentinel_alert_rule_scheduled.RegisterResource(),
		"azurerm_sentinel_automation_rule":                                              sentinel_automation_rule.RegisterResource(),
		"azurerm_sentinel_data_connector_aws_cloud_trail":                               sentinel_data_connector_aws_cloud_trail.RegisterResource(),
		"azurerm_sentinel_data_connector_azure_active_directory":                        sentinel_data_connector_azure_active_directory.RegisterResource(),
		"azurerm_sentinel_data_connector_azure_advanced_threat_protection":              sentinel_data_connector_azure_advanced_threat_protection.RegisterResource(),
		"azurerm_sentinel_data_connector_azure_security_center":                         sentinel_data_connector_azure_security_center.RegisterResource(),
		"azurerm_sentinel_data_connector_microsoft_cloud_app_security":                  sentinel_data_connector_microsoft_cloud_app_security.RegisterResource(),
		"azurerm_sentinel_data_connector_microsoft_defender_advanced_threat_protection": sentinel_data_connector_microsoft_defender_advanced_threat_protection.RegisterResource(),
		"azurerm_sentinel_data_connector_office_365":                                    sentinel_data_connector_office_365.RegisterResource(),
		"azurerm_sentinel_data_connector_office_atp":                                    sentinel_data_connector_office_atp.RegisterResource(),
		"azurerm_sentinel_data_connector_threat_intelligence":                           sentinel_data_connector_threat_intelligence.RegisterResource(),
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		sentinel_alert_rule_anomaly.DataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		sentinel_alert_rule_anomaly_built_in.Resource{},
		sentinel_alert_rule_anomaly_duplicate.Resource{},
		sentinel_alert_rule_threat_intelligence.Resource{},
		sentinel_data_connector_aws_s3.Resource{},
		sentinel_data_connector_dynamics_365.Resource{},
		sentinel_data_connector_iot.Resource{},
		sentinel_data_connector_microsoft_threat_intelligence.Resource{},
		sentinel_data_connector_microsoft_threat_protection.Resource{},
		sentinel_data_connector_office_365_project.Resource{},
		sentinel_data_connector_office_irm.Resource{},
		sentinel_data_connector_office_power_bi.Resource{},
		sentinel_data_connector_threat_intelligence_taxii.Resource{},
		sentinel_log_analytics_workspace_onboarding.Resource{},
		sentinel_metadata.Resource{},
		sentinel_threat_intelligence_indicator.Resource{},
		sentinel_watchlist_item.Resource{},
		sentinel_watchlist.Resource{},
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
