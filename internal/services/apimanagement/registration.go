// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_api"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_api_diagnostic"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_api_operation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_api_operation_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_api_operation_tag"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_api_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_api_release"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_api_schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_api_tag"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_api_tag_description"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_api_version_set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_authorization_server"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_backend"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_certificate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_custom_domain"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_diagnostic"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_email_template"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_gateway"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_gateway_api"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_gateway_certificate_authority"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_gateway_host_name_configuration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_global_schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_group_user"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_identity_provider_aad"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_identity_provider_aadb2c"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_identity_provider_facebook"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_identity_provider_google"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_identity_provider_microsoft"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_identity_provider_twitter"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_logger"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_named_value"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_notification_recipient_email"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_notification_recipient_user"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_openid_connect_provider"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_policy_fragment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_product"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_product_api"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_product_group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_product_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_product_tag"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_redis_cache"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_standalone_gateway"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_subscription"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_tag"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_user"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_workspace_api_version_set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_workspace_certificate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_workspace_named_value"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_workspace_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_workspace_policy_fragment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.TypedServiceRegistration                   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

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
		"azurerm_api_management":                                 api_management.RegisterDataSource(),
		"azurerm_api_management_api":                             api_management_api.RegisterDataSource(),
		"azurerm_api_management_api_version_set":                 api_management_api_version_set.RegisterDataSource(),
		"azurerm_api_management_gateway":                         api_management_gateway.RegisterDataSource(),
		"azurerm_api_management_gateway_host_name_configuration": api_management_gateway_host_name_configuration.RegisterDataSource(),
		"azurerm_api_management_group":                           api_management_group.RegisterDataSource(),
		"azurerm_api_management_product":                         api_management_product.RegisterDataSource(),
		"azurerm_api_management_user":                            api_management_user.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_api_management":                                 api_management.RegisterResource(),
		"azurerm_api_management_api":                             api_management_api.RegisterResource(),
		"azurerm_api_management_api_diagnostic":                  api_management_api_diagnostic.RegisterResource(),
		"azurerm_api_management_api_operation":                   api_management_api_operation.RegisterResource(),
		"azurerm_api_management_api_operation_policy":            api_management_api_operation_policy.RegisterResource(),
		"azurerm_api_management_api_operation_tag":               api_management_api_operation_tag.RegisterResource(),
		"azurerm_api_management_api_policy":                      api_management_api_policy.RegisterResource(),
		"azurerm_api_management_api_release":                     api_management_api_release.RegisterResource(),
		"azurerm_api_management_api_schema":                      api_management_api_schema.RegisterResource(),
		"azurerm_api_management_api_tag":                         api_management_api_tag.RegisterResource(),
		"azurerm_api_management_api_tag_description":             api_management_api_tag_description.RegisterResource(),
		"azurerm_api_management_api_version_set":                 api_management_api_version_set.RegisterResource(),
		"azurerm_api_management_authorization_server":            api_management_authorization_server.RegisterResource(),
		"azurerm_api_management_backend":                         api_management_backend.RegisterResource(),
		"azurerm_api_management_certificate":                     api_management_certificate.RegisterResource(),
		"azurerm_api_management_custom_domain":                   api_management_custom_domain.RegisterResource(),
		"azurerm_api_management_diagnostic":                      api_management_diagnostic.RegisterResource(),
		"azurerm_api_management_email_template":                  api_management_email_template.RegisterResource(),
		"azurerm_api_management_gateway":                         api_management_gateway.RegisterResource(),
		"azurerm_api_management_gateway_api":                     api_management_gateway_api.RegisterResource(),
		"azurerm_api_management_gateway_certificate_authority":   api_management_gateway_certificate_authority.RegisterResource(),
		"azurerm_api_management_gateway_host_name_configuration": api_management_gateway_host_name_configuration.RegisterResource(),
		"azurerm_api_management_global_schema":                   api_management_global_schema.RegisterResource(),
		"azurerm_api_management_group":                           api_management_group.RegisterResource(),
		"azurerm_api_management_group_user":                      api_management_group_user.RegisterResource(),
		"azurerm_api_management_identity_provider_aad":           api_management_identity_provider_aad.RegisterResource(),
		"azurerm_api_management_identity_provider_aadb2c":        api_management_identity_provider_aadb2c.RegisterResource(),
		"azurerm_api_management_identity_provider_facebook":      api_management_identity_provider_facebook.RegisterResource(),
		"azurerm_api_management_identity_provider_google":        api_management_identity_provider_google.RegisterResource(),
		"azurerm_api_management_identity_provider_microsoft":     api_management_identity_provider_microsoft.RegisterResource(),
		"azurerm_api_management_identity_provider_twitter":       api_management_identity_provider_twitter.RegisterResource(),
		"azurerm_api_management_logger":                          api_management_logger.RegisterResource(),
		"azurerm_api_management_named_value":                     api_management_named_value.RegisterResource(),
		"azurerm_api_management_openid_connect_provider":         api_management_openid_connect_provider.RegisterResource(),
		"azurerm_api_management_policy":                          api_management_policy.RegisterResource(),
		"azurerm_api_management_policy_fragment":                 api_management_policy_fragment.RegisterResource(),
		"azurerm_api_management_product":                         api_management_product.RegisterResource(),
		"azurerm_api_management_product_api":                     api_management_product_api.RegisterResource(),
		"azurerm_api_management_product_group":                   api_management_product_group.RegisterResource(),
		"azurerm_api_management_product_policy":                  api_management_product_policy.RegisterResource(),
		"azurerm_api_management_product_tag":                     api_management_product_tag.RegisterResource(),
		"azurerm_api_management_redis_cache":                     api_management_redis_cache.RegisterResource(),
		"azurerm_api_management_subscription":                    api_management_subscription.RegisterResource(),
		"azurerm_api_management_tag":                             api_management_tag.RegisterResource(),
		"azurerm_api_management_user":                            api_management_user.RegisterResource(),
	}
}

// DataSources returns a list of Data Sources supported by this Service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		api_management_subscription.DataSource{},
		api_management_workspace.DataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		api_management_notification_recipient_email.Resource{},
		api_management_notification_recipient_user.Resource{},
		api_management_standalone_gateway.Resource{},
		api_management_workspace_api_version_set.Resource{},
		api_management_workspace_certificate.Resource{},
		api_management_workspace_named_value.Resource{},
		api_management_workspace_policy_fragment.Resource{},
		api_management_workspace_policy.Resource{},
		api_management_workspace.Resource{},
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
