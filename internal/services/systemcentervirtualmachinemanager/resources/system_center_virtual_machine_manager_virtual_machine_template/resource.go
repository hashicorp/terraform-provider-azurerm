// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package system_center_virtual_machine_manager_virtual_machine_template

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/virtualmachinetemplates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -test-sequential

var (
	_ sdk.Resource             = Resource{}
	_ sdk.ResourceWithUpdate   = Resource{}
	_ sdk.ResourceWithIdentity = Resource{}
)

type Resource struct{}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return virtualmachinetemplates.ValidateVirtualMachineTemplateID
}

func (r Resource) ResourceType() string {
	return "azurerm_system_center_virtual_machine_manager_virtual_machine_template"
}

func (r Resource) Identity() resourceids.ResourceId {
	return &virtualmachinetemplates.VirtualMachineTemplateId{}
}
