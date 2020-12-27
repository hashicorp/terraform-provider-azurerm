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
		"azurerm_api_management_user":            dataSourceApiManagementUser(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_api_management":                             resourceApiManagementService(),
		"azurerm_api_management_api":                         resourceApiManagementApi(),
		"azurerm_api_management_api_diagnostic":              resourceApiManagementApiDiagnostic(),
		"azurerm_api_management_api_operation":               resourceApiManagementApiOperation(),
		"azurerm_api_management_api_operation_policy":        resourceApiManagementApiOperationPolicy(),
		"azurerm_api_management_api_policy":                  resourceApiManagementApiPolicy(),
		"azurerm_api_management_api_schema":                  resourceApiManagementApiSchema(),
		"azurerm_api_management_api_version_set":             resourceApiManagementApiVersionSet(),
		"azurerm_api_management_authorization_server":        resourceApiManagementAuthorizationServer(),
		"azurerm_api_management_backend":                     resourceApiManagementBackend(),
		"azurerm_api_management_certificate":                 resourceApiManagementCertificate(),
		"azurerm_api_management_custom_domain":               resourceApiManagementCustomDomain(),
		"azurerm_api_management_diagnostic":                  resourceApiManagementDiagnostic(),
		"azurerm_api_management_group":                       resourceApiManagementGroup(),
		"azurerm_api_management_group_user":                  resourceApiManagementGroupUser(),
		"azurerm_api_management_identity_provider_aad":       resourceApiManagementIdentityProviderAAD(),
		"azurerm_api_management_identity_provider_facebook":  resourceApiManagementIdentityProviderFacebook(),
		"azurerm_api_management_identity_provider_google":    resourceApiManagementIdentityProviderGoogle(),
		"azurerm_api_management_identity_provider_microsoft": resourceApiManagementIdentityProviderMicrosoft(),
		"azurerm_api_management_identity_provider_twitter":   resourceApiManagementIdentityProviderTwitter(),
		"azurerm_api_management_logger":                      resourceApiManagementLogger(),
		"azurerm_api_management_named_value":                 resourceApiManagementNamedValue(),
		"azurerm_api_management_openid_connect_provider":     resourceApiManagementOpenIDConnectProvider(),
		"azurerm_api_management_policy":                      resourceApiManagementPolicy(),
		"azurerm_api_management_product":                     resourceApiManagementProduct(),
		"azurerm_api_management_product_api":                 resourceApiManagementProductApi(),
		"azurerm_api_management_product_group":               resourceApiManagementProductGroup(),
		"azurerm_api_management_product_policy":              resourceApiManagementProductPolicy(),
		"azurerm_api_management_property":                    resourceApiManagementProperty(),
		"azurerm_api_management_subscription":                resourceApiManagementSubscription(),
		"azurerm_api_management_user":                        resourceApiManagementUser(),
	}
}
