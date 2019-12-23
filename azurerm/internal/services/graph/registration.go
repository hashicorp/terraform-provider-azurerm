package graph

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Graph"
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_azuread_application":       dataSourceArmAzureADApplication(),
		"azurerm_azuread_service_principal": dataSourceArmActiveDirectoryServicePrincipal()}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_azuread_application":                resourceArmActiveDirectoryApplication(),
		"azurerm_azuread_service_principal_password": resourceArmActiveDirectoryServicePrincipalPassword(),
		"azurerm_azuread_service_principal":          resourceArmActiveDirectoryServicePrincipal()}
}
