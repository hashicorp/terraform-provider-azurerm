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
			Description: "The name which should be used for this Exadata Database Storage Vault.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 255),
				validate.ExascaleDatabaseStorageVaultName,
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"additional_flash_cache_percentage": {
			Description:  "The size of additional Flash Cache in percentage of High Capacity database storage.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(0, 100),
		},

		"description": {
			Description: "Exadata Database Storage Vault description.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// Note: O+C API use `display_name` value if omitted
			Computed: true,
			ForceNew: true,
		},

		"display_name": {
			Description: "The user-friendly name for the Exadata Database Storage Vault resource. The name does not need to be unique.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 255),
				validate.ExascaleDatabaseStorageVaultName,
			),
		},

		"high_capacity_database_storage": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"total_size_in_gb": {
						Description: "Total capacity in gigabytes.",
						Type:        pluginsdk.TypeInt,
						Required:    true,
						ForceNew:    true,
					},
					"available_size_in_gb": {
						Description: "Available size in gigabytes.",
						Type:        pluginsdk.TypeInt,
						Computed:    true,
					},
				},
			},
		},

		"time_zone": {
			Description: "The time zone that you want to use for the Exadata Database Storage Vault. For details, see [Time Zones](https://docs.oracle.com/en/cloud/paas/base-database/time-zone/).",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     "UTC",
		},

		"tags": commonschema.Tags(),

		"zones": commonschema.ZonesMultipleRequiredForceNew(),
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
