// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package systemcentervirtualmachinemanager

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type Registration struct{}

var _ sdk.TypedServiceRegistration = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/systemcentervirtualmachinemanager"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "System Center Virtual Machine Manager"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"System Center Virtual Machine Manager",
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		SystemCenterVirtualMachineManagerInventoryItemsDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		SystemCenterVirtualMachineManagerCloudResource{},
		SystemCenterVirtualMachineManagerServerResource{},
		SystemCenterVirtualMachineManagerAvailabilitySetResource{},
		SystemCenterVirtualMachineManagerVirtualMachineInstanceResource{},
		SystemCenterVirtualMachineManagerVirtualNetworkResource{},
		SystemCenterVirtualMachineManagerVirtualMachineTemplateResource{},
		SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource{},
	}
}
