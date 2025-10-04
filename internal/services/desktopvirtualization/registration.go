// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package desktopvirtualization

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type Registration struct{}

var _ sdk.TypedServiceRegistration = Registration{}

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
		DesktopVirtualizationScalingPlanHostPoolAssociationResource{},
		DesktopVirtualizationScalingPlanResource{},
		DesktopVirtualizationWorkspaceApplicationGroupAssociationResource{},
		DesktopVirtualizationWorkspaceResource{},
	}
}
