// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/resources/management_group_template_deployment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/resources/management_lock"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/resources/resource_deployment_script_azure_cli"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/resources/resource_deployment_script_azure_power_shell"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/resources/resource_group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/resources/resource_group_template_deployment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/resources/resource_management_private_link"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/resources/resource_management_private_link_association"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/resources/resource_provider_feature_registration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/resources/resource_provider_registration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/resources/resources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/resources/subscription_template_deployment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/resources/template_spec_version"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/resources/tenant_template_deployment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.TypedServiceRegistrationWithAGitHubLabel   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
	_ sdk.FrameworkServiceRegistration               = Registration{}
)

type Registration struct{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/resource"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Resources"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Base",
		"Management",
		"Template",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_management_group_template_deployment": management_group_template_deployment.RegisterDataSource(),
		"azurerm_resource_group":                       resource_group.RegisterDataSource(),
		"azurerm_resource_group_template_deployment":   resource_group_template_deployment.RegisterDataSource(),
		"azurerm_resources":                            resources.RegisterDataSource(),
		"azurerm_subscription_template_deployment":     subscription_template_deployment.RegisterDataSource(),
		"azurerm_template_spec_version":                template_spec_version.RegisterDataSource(),
		"azurerm_tenant_template_deployment":           tenant_template_deployment.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_management_group_template_deployment": management_group_template_deployment.RegisterResource(),
		"azurerm_management_lock":                      management_lock.RegisterResource(),
		"azurerm_resource_group":                       resource_group.RegisterResource(),
		"azurerm_resource_group_template_deployment":   resource_group_template_deployment.RegisterResource(),
		"azurerm_subscription_template_deployment":     subscription_template_deployment.RegisterResource(),
		"azurerm_tenant_template_deployment":           tenant_template_deployment.RegisterResource(),
	}
}

// DataSources returns a list of Data Sources supported by this Service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

// Resources returns a list of Resources supported by this Service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		resource_deployment_script_azure_cli.Resource{},
		resource_deployment_script_azure_power_shell.Resource{},
		resource_management_private_link_association.Resource{},
		resource_management_private_link.Resource{},
		resource_provider_feature_registration.Resource{},
		resource_provider_registration.Resource{},
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
	return []sdk.FrameworkListWrappedResource{
		resource_group.ListResource{},
	}
}
