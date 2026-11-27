// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package system_center_virtual_machine_manager_virtual_machine_instance

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.Resource           = Resource{}
	_ sdk.ResourceWithUpdate = Resource{}
)

type Resource struct{}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.SystemCenterVirtualMachineManagerVirtualMachineInstanceID
}

func (r Resource) ResourceType() string {
	return "azurerm_system_center_virtual_machine_manager_virtual_machine_instance"
}
