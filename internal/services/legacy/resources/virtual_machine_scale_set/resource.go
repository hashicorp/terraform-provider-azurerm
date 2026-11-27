// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_scale_set

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachinescalesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/legacy/resources/virtual_machine_scale_set/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// NOTE: the `azurerm_virtual_machine_scale_set` resource has been superseded by the
//
//	`azurerm_linux_virtual_machine_scale_set` and `azurerm_windows_virtual_machine_scale_set` resources
//	and as such this resource is feature-frozen and new functionality will be added to these new resources instead.
func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualMachineScaleSetCreateUpdate,
		Read:   resourceVirtualMachineScaleSetRead,
		Update: resourceVirtualMachineScaleSetCreateUpdate,
		Delete: resourceVirtualMachineScaleSetDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.LegacyVMSSV0ToV1{},
		}),

		// NOTE: Don't remove with 4.0
		DeprecationMessage: `The 'azurerm_virtual_machine_scale_set' resource has been superseded by the 'azurerm_linux_virtual_machine_scale_set' and 'azurerm_windows_virtual_machine_scale_set' resources. Whilst this resource will continue to be available in the 2.x, 3.x and 4.x releases it is feature-frozen for compatibility purposes, will no longer receive any updates and will be removed in a future major release of the Azure Provider.`,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := virtualmachinescalesets.ParseVirtualMachineScaleSetID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: virtualMachineScaleSetSchema(),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(azureRmVirtualMachineScaleSetCustomizeDiff),
	}
}
