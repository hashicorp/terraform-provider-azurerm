// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_exascale_database_storage_vault

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),

		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 255),
				validate.ExascaleDatabaseStorageVaultName,
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"additional_flash_cache_percentage": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(0, 100),
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// Note: O+C API use `display_name` value if omitted
			Computed: true,
			ForceNew: true,
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 255),
				validate.ExascaleDatabaseStorageVaultName,
			),
		},

		"high_capacity_database_storage": {
			Type:     pluginsdk.TypeList,
			MinItems: 1,
			MaxItems: 1,
			Required: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"total_size_in_gb": {
						Type:     pluginsdk.TypeInt,
						Required: true,
						ForceNew: true,
					},
					"available_size_in_gb": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},

		"time_zone": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  "UTC",
		},

		"tags": commonschema.Tags(),

		"zones": commonschema.ZonesMultipleRequiredForceNew(),
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
