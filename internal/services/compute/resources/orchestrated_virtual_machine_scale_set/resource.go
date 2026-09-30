// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package orchestrated_virtual_machine_scale_set

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
		Create: resourceOrchestratedVirtualMachineScaleSetCreate,
		Read:   resourceOrchestratedVirtualMachineScaleSetRead,
		Update: resourceOrchestratedVirtualMachineScaleSetUpdate,
		Delete: resourceOrchestratedVirtualMachineScaleSetDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := commonids.ParseVirtualMachineScaleSetID(id)
			return err
		}, helpers.ImportOrchestratedVirtualMachineScaleSet),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		// The plan was to remove support the legacy Orchestrated Virtual Machine Scale Set in 3.0.
		// Turns out it's still in use
		// TODO: Revisit for 6.0
		// TODO: exposing requireGuestProvisionSignal once it's available
		// https://github.com/Azure/azure-rest-api-specs/pull/7246

		Schema: orchestratedVirtualMachineScaleSetSchema(),

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

			// SKU Profile validation
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
				skuName, hasSkuName := diff.GetOk("sku_name")
				_, hasSkuProfile := diff.GetOk("sku_profile")

				if hasSkuProfile {
					if !hasSkuName || skuName != SkuNameMix {
						return fmt.Errorf("`sku_profile` can only be configured when `sku_name` is set to `Mix`, got `%s`", skuName)
					}
				} else {
					if hasSkuName && skuName == SkuNameMix {
						return fmt.Errorf("`sku_profile` must be configured when `sku_name` is set to `Mix`")
					}
				}

				// Validate SKU Profile configuration
				if skuProfileRaw := diff.Get("sku_profile").([]any); len(skuProfileRaw) > 0 {
					skuProfile := skuProfileRaw[0].(map[string]any)
					allocationStrategy := skuProfile["allocation_strategy"].(string)

					configRaw := diff.GetRawConfig().AsValueMap()["sku_profile"]
					if configRaw.IsNull() || !configRaw.IsKnown() {
						return nil
					}

					skuProfiles := configRaw.AsValueSlice()
					if len(skuProfiles) == 0 {
						return nil
					}

					if vmSizes := skuProfile["virtual_machine_size"].(*pluginsdk.Set).List(); len(vmSizes) > 0 {
						for _, vmSize := range vmSizes {
							vmSizeMap := vmSize.(map[string]any)
							rank := vmSizeMap["rank"].(int)

							if rank != 0 && allocationStrategy != string(virtualmachinescalesets.AllocationStrategyPrioritized) {
								return fmt.Errorf("`rank` can only be set when `allocation_strategy` is `Prioritized`, got `%s`", allocationStrategy)
							}
						}
					}
				}

				return nil
			}),

			// Force recreation when sku_profile is removed and sku_name changes from mix
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
				oldName, newName := diff.GetChange("sku_name")

				if oldName.(string) == "Mix" && newName.(string) != "Mix" {
					if err := diff.ForceNew("sku_profile"); err != nil {
						return fmt.Errorf("forcing new resource when removing `sku_profile` with `sku_name` change from `Mix`: %+v", err)
					}
				}

				return nil
			}),

			// Upgrade policy validation
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
				upgradeMode := virtualmachinescalesets.UpgradeMode(diff.Get("upgrade_mode").(string))
				rollingUpgradePolicyRaw := diff.Get("rolling_upgrade_policy").([]any)

				if upgradeMode == virtualmachinescalesets.UpgradeModeManual && len(rollingUpgradePolicyRaw) > 0 {
					return fmt.Errorf("`rolling_upgrade_policy` cannot be specified when `upgrade_mode` is set to `%s`", string(upgradeMode))
				}

				if upgradeMode == virtualmachinescalesets.UpgradeModeRolling && len(rollingUpgradePolicyRaw) == 0 {
					return fmt.Errorf("`rolling_upgrade_policy` is required when `upgrade_mode` is set to `%s`", string(upgradeMode))
				}

				return nil
			}),

			// Network interface validation
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
				networkInterfaces := diff.Get("network_interface").([]any)
				for _, networkInterface := range networkInterfaces {
					raw := networkInterface.(map[string]any)
					auxiliaryMode := raw["auxiliary_mode"].(string)
					auxiliarySku := raw["auxiliary_sku"].(string)

					if auxiliaryMode != "" && auxiliarySku == "" {
						return fmt.Errorf("when `auxiliary_mode` is set, `auxiliary_sku` must also be set")
					}

					if auxiliarySku != "" && auxiliaryMode == "" {
						return fmt.Errorf("when `auxiliary_sku` is set, `auxiliary_mode` must also be set")
					}

					if auxiliaryMode != "" {
						networkApiVersion := virtualmachinescalesets.NetworkApiVersion(diff.Get("network_api_version").(string))
						if networkApiVersion == virtualmachinescalesets.NetworkApiVersionTwoZeroTwoZeroNegativeOneOneNegativeZeroOne {
							return fmt.Errorf("`auxiliary_mode` and `auxiliary_sku` can be set only when `network_api_version` is later than `2020-11-01`")
						}
					}
				}

				return nil
			}),
		),
	}
}
