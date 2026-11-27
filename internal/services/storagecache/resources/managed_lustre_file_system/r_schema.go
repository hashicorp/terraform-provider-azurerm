// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package managed_lustre_file_system

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2024-07-01/amlfilesystems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storagecache/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ManagedLustreFileSystemName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"maintenance_window": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"day_of_week": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(amlfilesystems.PossibleValuesForMaintenanceDayOfWeekType(), false),
					},

					"time_of_day_in_utc": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.TimeOfDayInUTC,
					},
				},
			},
		},

		"root_squash": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"mode": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(amlfilesystems.AmlFilesystemSquashModeAll),
							string(amlfilesystems.AmlFilesystemSquashModeRootOnly),
						}, false),
					},
					"no_squash_nids": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"squash_gid": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      0,
						ValidateFunc: validation.IntAtLeast(1),
					},
					"squash_uid": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      0,
						ValidateFunc: validation.IntAtLeast(1),
					},
				},
			},
		},

		"sku_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(PossibleValuesForSkuName(), false),
		},

		"storage_capacity_in_tb": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(4),
		},

		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"zones": commonschema.ZonesMultipleRequiredForceNew(),

		"hsm_setting": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"container_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateStorageContainerID,
					},

					"logging_container_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateStorageContainerID,
					},

					"import_prefix": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validate.ImportPrefix,
					},
				},
			},
		},

		"identity": commonschema.UserAssignedIdentityOptionalForceNew(),

		"encryption_key": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_url": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey),
					},

					"source_vault_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateKeyVaultID,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"mgs_address": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if oldVal, newVal := metadata.ResourceDiff.GetChange("encryption_key"); len(oldVal.([]any)) > 0 && len(newVal.([]any)) == 0 {
				if err := metadata.ResourceDiff.ForceNew("encryption_key"); err != nil {
					return err
				}
			}

			configSku := metadata.ResourceDiff.Get("sku_name")
			configCapacity := metadata.ResourceDiff.Get("storage_capacity_in_tb")
			skuProperties := GetSkuPropertiesByName(configSku.(string))

			if skuProperties != nil {
				if configCapacity.(int) < skuProperties.MinimumStorage {
					return fmt.Errorf("'storage_capacity_in_tb' field cannot be less than '%d' for the %q sku, got '%d'", skuProperties.MinimumStorage, configSku, configCapacity)
				}

				if configCapacity.(int)%skuProperties.MinimumStorage != 0 {
					return fmt.Errorf("'storage_capacity_in_tb' field must be defined in increments of '%d' for the %q sku, got '%d'", skuProperties.MinimumStorage, configSku, configCapacity)
				}
			}

			return nil
		},
	}
}
