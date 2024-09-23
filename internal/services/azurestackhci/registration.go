// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azurestackhci

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
	_ sdk.TypedServiceRegistrationWithAGitHubLabel   = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/azure-stack-hci"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Azure Stack HCI"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Azure Stack HCI",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_stack_hci_cluster": resourceArmStackHCICluster(),
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		StackHCIClusterDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		StackHCIDeploymentSettingResource{},
		StackHCILogicalNetworkResource{},
		StackHCIStoragePathResource{},
		StackHCIVirtualHardDiskResource{},
	}
}
