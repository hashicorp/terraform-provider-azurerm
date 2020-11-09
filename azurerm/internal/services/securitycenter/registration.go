package securitycenter

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_advanced_threat_protection":           resourceArmAdvancedThreatProtection(),
		"azurerm_security_center_contact":              resourceArmSecurityCenterContact(),
		"azurerm_security_center_setting":              resourceArmSecurityCenterSetting(),
		"azurerm_security_center_subscription_pricing": resourceArmSecurityCenterSubscriptionPricing(),
		"azurerm_security_center_workspace":            resourceArmSecurityCenterWorkspace(),
		"azurerm_security_center_auto_provisioning":    resourceArmSecurityCenterAutoProvisioning(),
	}
}
