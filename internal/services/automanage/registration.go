package automanage

import "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Automanage"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Automanage",
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_automanage_configuration_profile":                 resourceAutomanageConfigurationProfile(),
		"azurerm_automanage_configuration_profiles_version":        resourceAutomanageConfigurationProfilesVersion(),
		"azurerm_automanage_configuration_profile_assignment":      resourceAutomanageConfigurationProfileAssignment(),
		"azurerm_automanage_configuration_profile_hcrp_assignment": resourceAutomanageConfigurationProfileHCRPAssignment(),
		"azurerm_automanage_configuration_profile_hci_assignment":  resourceAutomanageConfigurationProfileHCIAssignment(),
	}
}
