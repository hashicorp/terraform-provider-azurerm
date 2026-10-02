// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_scale_set_extension

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachinescalesetextensions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// NOTE (also in the docs): this is not intended to be used with the `azurerm_virtual_machine_scale_set` resource

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualMachineScaleSetExtensionCreate,
		Read:   resourceVirtualMachineScaleSetExtensionRead,
		Update: resourceVirtualMachineScaleSetExtensionUpdate,
		Delete: resourceVirtualMachineScaleSetExtensionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := virtualmachinescalesetextensions.ParseVirtualMachineScaleSetExtensionID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: virtualMachineScaleSetExtensionSchema(),
	}
}
