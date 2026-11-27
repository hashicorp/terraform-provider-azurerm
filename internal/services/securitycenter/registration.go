// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package securitycenter

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/resources/advanced_threat_protection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/resources/iot_security_device_group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/resources/iot_security_solution"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/resources/security_center_assessment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/resources/security_center_assessment_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/resources/security_center_automation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/resources/security_center_contact"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/resources/security_center_server_vulnerability_assessment_virtual_machine"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/resources/security_center_server_vulnerability_assessments_setting"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/resources/security_center_setting"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/resources/security_center_storage_defender"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/resources/security_center_subscription_pricing"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/resources/security_center_workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
	_ sdk.TypedServiceRegistrationWithAGitHubLabel   = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/security-center"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Security Center"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Security Center",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_advanced_threat_protection":                                      advanced_threat_protection.RegisterResource(),
		"azurerm_iot_security_device_group":                                       iot_security_device_group.RegisterResource(),
		"azurerm_iot_security_solution":                                           iot_security_solution.RegisterResource(),
		"azurerm_security_center_assessment":                                      security_center_assessment.RegisterResource(),
		"azurerm_security_center_assessment_policy":                               security_center_assessment_policy.RegisterResource(),
		"azurerm_security_center_automation":                                      security_center_automation.RegisterResource(),
		"azurerm_security_center_contact":                                         security_center_contact.RegisterResource(),
		"azurerm_security_center_server_vulnerability_assessment_virtual_machine": security_center_server_vulnerability_assessment_virtual_machine.RegisterResource(),
		"azurerm_security_center_server_vulnerability_assessments_setting":        security_center_server_vulnerability_assessments_setting.RegisterResource(),
		"azurerm_security_center_setting":                                         security_center_setting.RegisterResource(),
		"azurerm_security_center_subscription_pricing":                            security_center_subscription_pricing.RegisterResource(),
		"azurerm_security_center_workspace":                                       security_center_workspace.RegisterResource(),
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		security_center_storage_defender.Resource{},
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
