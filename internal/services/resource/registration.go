// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
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
		"azurerm_resources":                            dataSourceResources(),
		"azurerm_resource_group":                       dataSourceResourceGroup(),
		"azurerm_template_spec_version":                dataSourceTemplateSpecVersion(),
		"azurerm_management_group_template_deployment": dataSourceManagementGroupTemplateDeployment(),
		"azurerm_resource_group_template_deployment":   dataSourceResourceGroupTemplateDeployment(),
		"azurerm_subscription_template_deployment":     dataSourceSubscriptionTemplateDeployment(),
		"azurerm_tenant_template_deployment":           dataSourceTenantTemplateDeployment(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	resources := map[string]*pluginsdk.Resource{
		"azurerm_management_lock":                      resourceManagementLock(),
		"azurerm_management_group_template_deployment": managementGroupTemplateDeploymentResource(),
		"azurerm_resource_group":                       resourceResourceGroup(),
		"azurerm_resource_group_template_deployment":   resourceGroupTemplateDeploymentResource(),
		"azurerm_subscription_template_deployment":     subscriptionTemplateDeploymentResource(),
		"azurerm_tenant_template_deployment":           tenantTemplateDeploymentResource(),
	}

	return resources
}

// DataSources returns a list of Data Sources supported by this Service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

// Resources returns a list of Resources supported by this Service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		ResourceManagementPrivateLinkAssociationResource{},
		ResourceProviderRegistrationResource{},
		ResourceManagementPrivateLinkResource{},
		ResourceDeploymentScriptAzurePowerShellResource{},
		ResourceDeploymentScriptAzureCliResource{},
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
		ResourceGroupListResource{},
	}
}
