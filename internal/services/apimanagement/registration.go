// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/api-management"
}

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
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_api_management":                                 dataSourceApiManagementService(),
		"azurerm_api_management_api":                             dataSourceApiManagementApi(),
		"azurerm_api_management_api_version_set":                 dataSourceApiManagementApiVersionSet(),
		"azurerm_api_management_gateway":                         dataSourceApiManagementGateway(),
		"azurerm_api_management_gateway_host_name_configuration": dataSourceApiManagementGatewayHostNameConfiguration(),
		"azurerm_api_management_group":                           dataSourceApiManagementGroup(),
		"azurerm_api_management_product":                         dataSourceApiManagementProduct(),
		"azurerm_api_management_user":                            dataSourceApiManagementUser(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_api_management":                                 resourceApiManagementService(),
		"azurerm_api_management_api":                             resourceApiManagementApi(),
		"azurerm_api_management_api_diagnostic":                  resourceApiManagementApiDiagnostic(),
		"azurerm_api_management_api_operation":                   resourceApiManagementApiOperation(),
		"azurerm_api_management_api_operation_policy":            resourceApiManagementApiOperationPolicy(),
		"azurerm_api_management_api_operation_tag":               resourceApiManagementApiOperationTag(),
		"azurerm_api_management_api_policy":                      resourceApiManagementApiPolicy(),
		"azurerm_api_management_api_release":                     resourceApiManagementApiRelease(),
		"azurerm_api_management_api_schema":                      resourceApiManagementApiSchema(),
		"azurerm_api_management_api_tag":                         resourceApiManagementApiTag(),
		"azurerm_api_management_api_tag_description":             resourceApiManagementApiTagDescription(),
		"azurerm_api_management_api_version_set":                 resourceApiManagementApiVersionSet(),
		"azurerm_api_management_authorization_server":            resourceApiManagementAuthorizationServer(),
		"azurerm_api_management_backend":                         resourceApiManagementBackend(),
		"azurerm_api_management_certificate":                     resourceApiManagementCertificate(),
		"azurerm_api_management_custom_domain":                   resourceApiManagementCustomDomain(),
		"azurerm_api_management_diagnostic":                      resourceApiManagementDiagnostic(),
		"azurerm_api_management_email_template":                  resourceApiManagementEmailTemplate(),
		"azurerm_api_management_gateway":                         resourceApiManagementGateway(),
		"azurerm_api_management_gateway_api":                     resourceApiManagementGatewayApi(),
		"azurerm_api_management_gateway_certificate_authority":   resourceApiManagementGatewayCertificateAuthority(),
		"azurerm_api_management_gateway_host_name_configuration": resourceApiManagementGatewayHostNameConfiguration(),
		"azurerm_api_management_global_schema":                   resourceApiManagementGlobalSchema(),
		"azurerm_api_management_group":                           resourceApiManagementGroup(),
		"azurerm_api_management_group_user":                      resourceApiManagementGroupUser(),
		"azurerm_api_management_identity_provider_aad":           resourceApiManagementIdentityProviderAAD(),
		"azurerm_api_management_identity_provider_aadb2c":        resourceArmApiManagementIdentityProviderAADB2C(),
		"azurerm_api_management_identity_provider_facebook":      resourceApiManagementIdentityProviderFacebook(),
		"azurerm_api_management_identity_provider_google":        resourceApiManagementIdentityProviderGoogle(),
		"azurerm_api_management_identity_provider_microsoft":     resourceApiManagementIdentityProviderMicrosoft(),
		"azurerm_api_management_identity_provider_twitter":       resourceApiManagementIdentityProviderTwitter(),
		"azurerm_api_management_logger":                          resourceApiManagementLogger(),
		"azurerm_api_management_named_value":                     resourceApiManagementNamedValue(),
		"azurerm_api_management_openid_connect_provider":         resourceApiManagementOpenIDConnectProvider(),
		"azurerm_api_management_policy":                          resourceApiManagementPolicy(),
		"azurerm_api_management_policy_fragment":                 resourceApiManagementPolicyFragment(),
		"azurerm_api_management_product":                         resourceApiManagementProduct(),
		"azurerm_api_management_product_api":                     resourceApiManagementProductApi(),
		"azurerm_api_management_product_group":                   resourceApiManagementProductGroup(),
		"azurerm_api_management_product_policy":                  resourceApiManagementProductPolicy(),
		"azurerm_api_management_product_tag":                     resourceApiManagementProductTag(),
		"azurerm_api_management_redis_cache":                     resourceApiManagementRedisCache(),
		"azurerm_api_management_subscription":                    resourceApiManagementSubscription(),
		"azurerm_api_management_tag":                             resourceApiManagementTag(),
		"azurerm_api_management_user":                            resourceApiManagementUser(),
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		ApiManagementNotificationRecipientEmailResource{},
		ApiManagementNotificationRecipientUserResource{},
	}
}
