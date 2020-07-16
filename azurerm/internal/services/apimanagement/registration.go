package apimanagement

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "API Management"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"API Management",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_api_management":                 dataSourceApiManagementService(),
		"azurerm_api_management_api":             dataSourceApiManagementApi(),
		"azurerm_api_management_api_version_set": dataSourceApiManagementApiVersionSet(),
		"azurerm_api_management_group":           dataSourceApiManagementGroup(),
		"azurerm_api_management_product":         dataSourceApiManagementProduct(),
		"azurerm_api_management_user":            dataSourceArmApiManagementUser(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_api_management":                             resourceArmApiManagementService(),
		"azurerm_api_management_api":                         resourceArmApiManagementApi(),
		"azurerm_api_management_api_diagnostic":              resourceArmApiManagementApiDiagnostic(),
		"azurerm_api_management_api_operation":               resourceArmApiManagementApiOperation(),
		"azurerm_api_management_api_operation_policy":        resourceArmApiManagementApiOperationPolicy(),
		"azurerm_api_management_api_policy":                  resourceArmApiManagementApiPolicy(),
		"azurerm_api_management_api_schema":                  resourceArmApiManagementApiSchema(),
		"azurerm_api_management_api_version_set":             resourceArmApiManagementApiVersionSet(),
		"azurerm_api_management_authorization_server":        resourceArmApiManagementAuthorizationServer(),
		"azurerm_api_management_backend":                     resourceArmApiManagementBackend(),
		"azurerm_api_management_certificate":                 resourceArmApiManagementCertificate(),
		"azurerm_api_management_custom_domain":               resourceArmApiManagementCustomDomain(),
		"azurerm_api_management_diagnostic":                  resourceArmApiManagementDiagnostic(),
		"azurerm_api_management_group":                       resourceArmApiManagementGroup(),
		"azurerm_api_management_group_user":                  resourceArmApiManagementGroupUser(),
		"azurerm_api_management_identity_provider_aad":       resourceArmApiManagementIdentityProviderAAD(),
		"azurerm_api_management_identity_provider_aadb2c":    resourceArmApiManagementIdentityProviderAADB2C(),
		"azurerm_api_management_identity_provider_facebook":  resourceArmApiManagementIdentityProviderFacebook(),
		"azurerm_api_management_identity_provider_google":    resourceArmApiManagementIdentityProviderGoogle(),
		"azurerm_api_management_identity_provider_microsoft": resourceArmApiManagementIdentityProviderMicrosoft(),
		"azurerm_api_management_identity_provider_twitter":   resourceArmApiManagementIdentityProviderTwitter(),
		"azurerm_api_management_logger":                      resourceArmApiManagementLogger(),
		"azurerm_api_management_named_value":                 resourceArmApiManagementNamedValue(),
		"azurerm_api_management_openid_connect_provider":     resourceArmApiManagementOpenIDConnectProvider(),
		"azurerm_api_management_product":                     resourceArmApiManagementProduct(),
		"azurerm_api_management_product_api":                 resourceArmApiManagementProductApi(),
		"azurerm_api_management_product_group":               resourceArmApiManagementProductGroup(),
		"azurerm_api_management_product_policy":              resourceArmApiManagementProductPolicy(),
		"azurerm_api_management_property":                    resourceArmApiManagementProperty(),
		"azurerm_api_management_subscription":                resourceArmApiManagementSubscription(),
		"azurerm_api_management_user":                        resourceArmApiManagementUser(),
	}
}
