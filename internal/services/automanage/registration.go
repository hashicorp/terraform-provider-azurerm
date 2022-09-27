package automanage

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

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_automanage_configuration_profile":                dataSourceAutomanageConfigurationProfile(),
		"azurerm_automanage_configuration_profiles_version":       dataSourceAutomanageConfigurationProfilesVersion(),
		"azurerm_automanage_configuration_profile_assignment":     dataSourceAutomanageConfigurationProfileAssignment(),
		"azurerm_automanage_configuration_profile_hcrpassignment": dataSourceAutomanageConfigurationProfileHCRPAssignment(),
		"azurerm_automanage_configuration_profile_hciassignment":  dataSourceAutomanageConfigurationProfileHCIAssignment(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_automanage_configuration_profile":                resourceAutomanageConfigurationProfile(),
		"azurerm_automanage_configuration_profiles_version":       resourceAutomanageConfigurationProfilesVersion(),
		"azurerm_automanage_configuration_profile_assignment":     resourceAutomanageConfigurationProfileAssignment(),
		"azurerm_automanage_configuration_profile_hcrpassignment": resourceAutomanageConfigurationProfileHCRPAssignment(),
		"azurerm_automanage_configuration_profile_hciassignment":  resourceAutomanageConfigurationProfileHCIAssignment(),
	}
}
