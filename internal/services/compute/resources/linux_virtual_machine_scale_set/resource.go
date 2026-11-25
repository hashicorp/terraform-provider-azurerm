// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package linux_virtual_machine_scale_set

import (
	"context"
	"errors"
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

type Properties struct {
	Output Output `tfschema:"output"`
}

type Output struct {
	AccessSAS string `tfschema:"accessSAS"`
}

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages a Linux Virtual Machine Scale Set.",
		Create:      resourceLinuxVirtualMachineScaleSetCreate,
		Read:        resourceLinuxVirtualMachineScaleSetRead,
		Update:      resourceLinuxVirtualMachineScaleSetUpdate,
		Delete:      resourceLinuxVirtualMachineScaleSetDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := commonids.ParseVirtualMachineScaleSetID(id)
			return err
		}, helpers.ImportVirtualMachineScaleSet(virtualmachinescalesets.OperatingSystemTypesLinux, "azurerm_linux_virtual_machine_scale_set")),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(time.Minute * 60),
			Read:   pluginsdk.DefaultTimeout(time.Minute * 5),
			Update: pluginsdk.DefaultTimeout(time.Minute * 60),
			Delete: pluginsdk.DefaultTimeout(time.Minute * 60),
		},

		// TODO: exposing requireGuestProvisionSignal once it's available
		// https://github.com/Azure/azure-rest-api-specs/pull/7246

		Schema: resourceLinuxVirtualMachineScaleSetSchema(),

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

				if !diff.GetRawConfig().GetAttr("capacity_reservation_group_id").IsNull() {
					if diff.Get("single_placement_group").(bool) {
						return errors.New("`single_placement_group` must be set to `false` when `capacity_reservation_group_id` is specified")
					}
				}

				return nil
			}),

			// `zone_balance` can only be enabled when availability zones are specified
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
				// `zone_balance` and `zones` can be unknown during plan (e.g. a reference to another resource), in which
				// case `Get` looks empty, so the raw config is used here
				rawConfig := diff.GetRawConfig()
				if rawConfig.IsNull() || !rawConfig.IsKnown() {
					return nil
				}

				rawValues := rawConfig.AsValueMap()
				zoneBalance, rawZones := rawValues["zone_balance"], rawValues["zones"]
				if zoneBalance.IsNull() || !zoneBalance.IsKnown() || !rawZones.IsKnown() {
					return nil
				}

				if zoneBalance.True() && (rawZones.IsNull() || rawZones.LengthInt() == 0) {
					return fmt.Errorf("`zone_balance` can only be set to `true` when zones are specified")
				}

				return nil
			}),

			pluginsdk.CustomizeDiffShim(helpers.VirtualMachineScaleSetAutomaticZoneRebalancingCustomizeDiff),

			pluginsdk.ForceNewIf("capacity_reservation_group_id", func(ctx context.Context, d *pluginsdk.ResourceDiff, _ any) bool {
				oldCRG, _ := d.GetChange("capacity_reservation_group_id")
				oldZones, _ := d.GetChange("zones")
				// In-place association is only supported when adding a CRG to an already-zonal scale set. Updating or de-associating still requires a force new (actually requires deallocation).
				return oldZones.(*schema.Set).Len() == 0 || oldCRG.(string) != ""
			}),
		),
	}
}
