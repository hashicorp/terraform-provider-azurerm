// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package azurestackhci

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/azurestackhci/resources/stack_hci_cluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/azurestackhci/resources/stack_hci_deployment_setting"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/azurestackhci/resources/stack_hci_extension"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/azurestackhci/resources/stack_hci_logical_network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/azurestackhci/resources/stack_hci_marketplace_gallery_image"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/azurestackhci/resources/stack_hci_network_interface"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/azurestackhci/resources/stack_hci_storage_path"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/azurestackhci/resources/stack_hci_virtual_hard_disk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
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
		"azurerm_stack_hci_cluster": stack_hci_cluster.RegisterResource(),
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		stack_hci_cluster.DataSource{},
		stack_hci_storage_path.DataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		stack_hci_deployment_setting.Resource{},
		stack_hci_extension.Resource{},
		stack_hci_logical_network.Resource{},
		stack_hci_marketplace_gallery_image.Resource{},
		stack_hci_network_interface.Resource{},
		stack_hci_storage_path.Resource{},
		stack_hci_virtual_hard_disk.Resource{},
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
