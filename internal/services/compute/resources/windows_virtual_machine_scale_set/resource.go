// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package windows_virtual_machine_scale_set

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachinescalesets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceWindowsVirtualMachineScaleSetCreate,
		Read:   resourceWindowsVirtualMachineScaleSetRead,
		Update: resourceWindowsVirtualMachineScaleSetUpdate,
		Delete: resourceWindowsVirtualMachineScaleSetDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := commonids.ParseVirtualMachineScaleSetID(id)
			return err
		}, helpers.ImportVirtualMachineScaleSet(virtualmachinescalesets.OperatingSystemTypesWindows, "azurerm_windows_virtual_machine_scale_set")),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		// TODO: exposing requireGuestProvisionSignal once it's available
		// https://github.com/Azure/azure-rest-api-specs/pull/7246

		Schema: resourceWindowsVirtualMachineScaleSetSchema(),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			// Removing existing zones is currently not supported for Virtual Machine Scale Sets
			pluginsdk.ForceNewIfChange("zones", func(ctx context.Context, old, new, meta any) bool {
				oldZones := zones.ExpandUntyped(old.(*schema.Set).List())
				newZones := zones.ExpandUntyped(new.(*schema.Set).List())

				for _, ov := range oldZones {
					found := slices.Contains(newZones, ov)

					if !found {
						return true
					}
				}

				return false
			}),

			pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
				networkInterfaces := diff.Get("network_interface").([]any)
				for _, v := range networkInterfaces {
					raw := v.(map[string]any)
					auxiliaryMode := raw["auxiliary_mode"].(string)
					auxiliarySku := raw["auxiliary_sku"].(string)

					if auxiliaryMode != "" && auxiliarySku == "" {
						return fmt.Errorf("when `auxiliary_mode` is set, `auxiliary_sku` must also be set")
					}

					if auxiliarySku != "" && auxiliaryMode == "" {
						return fmt.Errorf("when `auxiliary_sku` is set, `auxiliary_mode` must also be set")
					}
				}

				return nil
			}),
		),
	}
}
