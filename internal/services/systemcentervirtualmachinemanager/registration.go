// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package systemcentervirtualmachinemanager

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/resources/system_center_virtual_machine_manager_availability_set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/resources/system_center_virtual_machine_manager_cloud"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/resources/system_center_virtual_machine_manager_inventory_items"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/resources/system_center_virtual_machine_manager_server"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/resources/system_center_virtual_machine_manager_virtual_machine_instance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/resources/system_center_virtual_machine_manager_virtual_machine_instance_guest_agent"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/resources/system_center_virtual_machine_manager_virtual_machine_template"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/resources/system_center_virtual_machine_manager_virtual_network"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration = Registration{}
	_ sdk.TypedServiceRegistration     = Registration{}
)

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
		system_center_virtual_machine_manager_inventory_items.DataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		system_center_virtual_machine_manager_availability_set.Resource{},
		system_center_virtual_machine_manager_cloud.Resource{},
		system_center_virtual_machine_manager_server.Resource{},
		system_center_virtual_machine_manager_virtual_machine_instance_guest_agent.Resource{},
		system_center_virtual_machine_manager_virtual_machine_instance.Resource{},
		system_center_virtual_machine_manager_virtual_machine_template.Resource{},
		system_center_virtual_machine_manager_virtual_network.Resource{},
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
		system_center_virtual_machine_manager_availability_set.ListResource{},
		system_center_virtual_machine_manager_cloud.ListResource{},
		system_center_virtual_machine_manager_server.ListResource{},
		system_center_virtual_machine_manager_virtual_machine_template.ListResource{},
		system_center_virtual_machine_manager_virtual_network.ListResource{},
	}
}
