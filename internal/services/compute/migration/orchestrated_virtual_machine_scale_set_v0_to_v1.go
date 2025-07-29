// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = OrchestratedVirtualMachineScaleSetUpgradeV0ToV1{}

type OrchestratedVirtualMachineScaleSetUpgradeV0ToV1 struct{}

func (OrchestratedVirtualMachineScaleSetUpgradeV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return orchestratedVirtualMachineScaleSetSchemaForV0()
}

func (OrchestratedVirtualMachineScaleSetUpgradeV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		log.Printf("[DEBUG] Migrating OVMSS from v0 to v1 - upgrading vm_sizes from string list to object list")

		// Find the sku_profile in the state
		if skuProfileRaw, exists := rawState["sku_profile"]; exists && skuProfileRaw != nil {
			skuProfileList := skuProfileRaw.([]interface{})
			if len(skuProfileList) > 0 && skuProfileList[0] != nil {
				skuProfile := skuProfileList[0].(map[string]interface{})

				// Check if vm_sizes exists and needs migration
				if vmSizesRaw, exists := skuProfile["vm_sizes"]; exists && vmSizesRaw != nil {
					vmSizesList := vmSizesRaw.([]interface{})
					migratedVMSizes := make([]interface{}, 0, len(vmSizesList))

					// Convert each string VM size to an object with name field
					for _, vmSizeRaw := range vmSizesList {
						if vmSizeStr, ok := vmSizeRaw.(string); ok {
							migratedVMSize := map[string]interface{}{
								"name": vmSizeStr,
								// rank is optional and not set in v0, so we don't include it
							}
							migratedVMSizes = append(migratedVMSizes, migratedVMSize)
						}
					}

					// Update the sku_profile with migrated vm_sizes
					skuProfile["vm_sizes"] = migratedVMSizes
					skuProfileList[0] = skuProfile
					rawState["sku_profile"] = skuProfileList

					log.Printf("[DEBUG] Successfully migrated %d vm_sizes from string list to object list", len(migratedVMSizes))
				}
			}
		}

		return rawState, nil
	}
}

// orchestratedVirtualMachineScaleSetSchemaForV0 returns the schema used in v0
func orchestratedVirtualMachineScaleSetSchemaForV0() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"sku_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"allocation_strategy": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"vm_sizes": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						MinItems: 1,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString, // v0 schema - simple string list
						},
					},
				},
			},
		},

		// Include other fields that might exist in state but are not relevant to migration
		"platform_fault_domain_count": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"proximity_placement_group_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"single_placement_group": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"zones": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"instances": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"network_api_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		// Add other fields that exist in the actual schema but aren't critical for migration
		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}
