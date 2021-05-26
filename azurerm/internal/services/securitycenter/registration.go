package securitycenter

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

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
		"azurerm_advanced_threat_protection":                      resourceAdvancedThreatProtection(),
		"azurerm_iot_security_device_group":                       resourceIotSecurityDeviceGroup(),
		"azurerm_iot_security_solution":                           resourceIotSecuritySolution(),
		"azurerm_security_center_assessment":                      resourceSecurityCenterAssessment(),
		"azurerm_security_center_assessment_metadata":             resourceArmSecurityCenterAssessmentMetadata(),
		"azurerm_security_center_assessment_policy":               resourceArmSecurityCenterAssessmentPolicy(),
		"azurerm_security_center_contact":                         resourceSecurityCenterContact(),
		"azurerm_security_center_setting":                         resourceSecurityCenterSetting(),
		"azurerm_security_center_subscription_pricing":            resourceSecurityCenterSubscriptionPricing(),
		"azurerm_security_center_workspace":                       resourceSecurityCenterWorkspace(),
		"azurerm_security_center_automation":                      resourceSecurityCenterAutomation(),
		"azurerm_security_center_auto_provisioning":               resourceSecurityCenterAutoProvisioning(),
		"azurerm_security_center_server_vulnerability_assessment": resourceServerVulnerabilityAssessment(),
	}
}
