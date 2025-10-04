// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package desktopvirtualization

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
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
		DesktopVirtualizationHostPoolDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		DesktopVirtualizationApplicationGroupResource{},
		DesktopVirtualizationApplicationResource{},
		DesktopVirtualizationHostPoolResource{},
		DesktopVirtualizationHostPoolRegistrationInfoResource{},
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_virtual_desktop_workspace":                               resourceArmDesktopVirtualizationWorkspace(),
		"azurerm_virtual_desktop_scaling_plan":                            resourceVirtualDesktopScalingPlan(),
		"azurerm_virtual_desktop_workspace_application_group_association": resourceVirtualDesktopWorkspaceApplicationGroupAssociation(),
		"azurerm_virtual_desktop_scaling_plan_host_pool_association":      resourceVirtualDesktopScalingPlanHostPoolAssociation(),
	}
}
