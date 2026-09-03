// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package desktopvirtualization

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.TypedServiceRegistration                   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/virtual-desktops"
}

func (r Registration) Name() string {
	return "Desktop Virtualization"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Desktop Virtualization",
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		DesktopVirtualizationWorkspaceDataSource{},
		DesktopVirtualizationApplicationGroupDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_virtual_desktop_host_pool": dataSourceVirtualDesktopHostPool(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_virtual_desktop_workspace":                               resourceArmDesktopVirtualizationWorkspace(),
		"azurerm_virtual_desktop_host_pool":                               resourceVirtualDesktopHostPool(),
		"azurerm_virtual_desktop_scaling_plan":                            resourceVirtualDesktopScalingPlan(),
		"azurerm_virtual_desktop_application_group":                       resourceVirtualDesktopApplicationGroup(),
		"azurerm_virtual_desktop_application":                             resourceVirtualDesktopApplication(),
		"azurerm_virtual_desktop_workspace_application_group_association": resourceVirtualDesktopWorkspaceApplicationGroupAssociation(),
		"azurerm_virtual_desktop_scaling_plan_host_pool_association":      resourceVirtualDesktopScalingPlanHostPoolAssociation(),
		"azurerm_virtual_desktop_host_pool_registration_info":             resourceVirtualDesktopHostPoolRegistrationInfo(),
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
