// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_machine

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// NOTE: the `azurerm_virtual_machine` resource has been superseded by the `azurerm_linux_virtual_machine` and
//
//	`azurerm_windows_virtual_machine` resources - as such this resource is feature-frozen and new
//	functionality will be added to these new resources instead.
func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualMachineCreateUpdate,
		Read:   resourceVirtualMachineRead,
		Update: resourceVirtualMachineCreateUpdate,
		Delete: resourceVirtualMachineDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := virtualmachines.ParseVirtualMachineID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: virtualMachineSchema(),
	}
}
