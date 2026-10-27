// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"maps"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var postgresqlFlexibleServerResourceName = "azurerm_postgresql_flexible_server"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePostgresqlFlexibleServerCreate,
		Read:   resourcePostgresqlFlexibleServerRead,
		Update: resourcePostgresqlFlexibleServerUpdate,
		Delete: resourcePostgresqlFlexibleServerDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(1 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(1 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(1 * time.Hour),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := servers.ParseFlexibleServerID(id)
			return err
		}),

		Schema: postgresqlFlexibleServerSchema(),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			func(ctx context.Context, d *pluginsdk.ResourceDiff, v any) error {
				if d.HasChange("version") {
					oldVersionVal, newVersionVal := d.GetChange("version")
					// `version` value has been validated already, ignore the parse errors is safe
					oldVersion, _ := strconv.ParseInt(oldVersionVal.(string), 10, 32)
					newVersion, _ := strconv.ParseInt(newVersionVal.(string), 10, 32)

					if oldVersion > newVersion {
						d.ForceNew("version")
					}
					return nil
				}
				return nil
			}, func(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
				oldLoginName, _ := diff.GetChange("administrator_login")
				if oldLoginName != "" {
					diff.ForceNew("administrator_login")
				}
				return nil
			}, func(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
				if diff.Get("storage_type").(string) == string(servers.StorageTypePremiumVTwoLRS) {
					return nil
				}

				storageTierMappings := validate.InitializeFlexibleServerStorageTierDefaults()
				var newTier string
				var newMb int
				var isValid bool

				oldStorageMbRaw, newStorageMbRaw := diff.GetChange("storage_mb")
				oldTierRaw, newTierRaw := diff.GetChange("storage_tier")

				if oldStorageMbRaw.(int) == 0 && oldTierRaw.(string) == "" && newStorageMbRaw.(int) == 0 && newTierRaw.(string) == "" {
					// This is a new resource without any values in the state
					// or config, default values will be set in create...
					return nil
				}

				newMb = newStorageMbRaw.(int)
				newTier = newTierRaw.(string)

				// if newMb is smaller than oldStorageMb, it's a downgrade need to trigger a force new
				if newMb > 0 && oldStorageMbRaw.(int) > newMb {
					diff.ForceNew("storage_mb")
				}

				// if newMb or newTier values are empty,
				// assign the default values that will
				// be assigned in the create func...
				if newMb == 0 {
					newMb = 32768
				}

				// get the valid mappings for the passed
				// storage_mb size...
				storageTiers, ok := storageTierMappings[newMb]
				if !ok {
					return nil
				}

				if newTier == "" {
					newTier = string(storageTiers.DefaultTier)
				}

				// verify that the storage_tier is valid
				// for the given storage_mb...
				if slices.Contains(*storageTiers.ValidTiers, newTier) {
					isValid = true
				}

				if !isValid {
					if strings.EqualFold(oldTierRaw.(string), newTier) {
						// The tier value did not change, so we need to determine if they are
						// using the default value for the tier, or they actually defined the
						// tier in the config or not... If they did not define
						// the tier in the config we need to assign a new valid default
						// tier for the newMb value. However, if the tier is in the config
						// this is a valid error and should be returned...
						if v := diff.GetRawConfig().AsValueMap()["storage_tier"]; v.IsNull() {
							diff.SetNew("storage_tier", string(storageTiers.DefaultTier))
							log.Printf("[DEBUG]: 'storage_tier' was not valid and was not in the config assigning new default 'storage_tier' %q -> %q\n", newTier, storageTiers.DefaultTier)
							return nil
						}
					}

					return fmt.Errorf("invalid 'storage_tier' %q for defined 'storage_mb' size '%d', expected one of %q", newTier, newMb, *storageTiers.ValidTiers)
				}

				return nil
			}, func(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
				oldIdentityRaw, newIdentityRaw := diff.GetChange("identity")
				oldIdentity := oldIdentityRaw.([]any)
				oldIdentityType := string(identity.TypeNone)
				if len(oldIdentity) > 0 {
					oldIdentityBlock := oldIdentity[0].(map[string]any)
					oldIdentityType = oldIdentityBlock["type"].(string)
				}

				newIdentity := newIdentityRaw.([]any)
				newIdentityType := string(identity.TypeNone)
				if len(newIdentity) > 0 {
					newIdentityBlock := newIdentity[0].(map[string]any)
					newIdentityType = newIdentityBlock["type"].(string)
				}

				if (oldIdentityType == string(identity.TypeUserAssigned) && newIdentityType == string(identity.TypeSystemAssigned)) || (oldIdentityType == string(identity.TypeUserAssigned) && newIdentityType == string(identity.TypeNone)) || (oldIdentityType == string(identity.TypeSystemAssignedUserAssigned) && newIdentityType == string(identity.TypeSystemAssigned)) || (oldIdentityType == string(identity.TypeSystemAssignedUserAssigned) && newIdentityType == string(identity.TypeNone)) {
					diff.ForceNew("identity.0.type")
				}

				return nil
			}, func(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
				configMap := diff.GetRawConfig().AsValueMap()
				createMode := diff.Get("create_mode").(string)

				if diff.Get("storage_type").(string) == string(servers.StorageTypePremiumVTwoLRS) {
					version := diff.Get("version").(string)
					if version == string(servers.PostgresMajorVersionOneOne) || version == string(servers.PostgresMajorVersionOneTwo) || version == string(servers.PostgresMajorVersionOneThree) {
						return fmt.Errorf("PostgreSQL version `%s` is not supported when `storage_type` is `PremiumV2_LRS`", version)
					}

					if skuName, ok := diff.GetOk("sku_name"); ok {
						if strings.HasPrefix(skuName.(string), "B_") {
							return errors.New("burstable compute tier is not supported when `storage_type` is `PremiumV2_LRS`")
						}
					}

					if v := configMap["storage_tier"]; !v.IsNull() {
						return errors.New("`storage_tier` is not supported when `storage_type` is `PremiumV2_LRS`")
					}

					if diff.Get("auto_grow_enabled").(bool) {
						return errors.New("`auto_grow_enabled` is not supported when `storage_type` is `PremiumV2_LRS`")
					}

					if diff.Get("geo_redundant_backup_enabled").(bool) {
						if _, ok := diff.GetOk("customer_managed_key"); ok {
							return errors.New("`geo_redundant_backup_enabled` with `customer_managed_key` is not supported when `storage_type` is `PremiumV2_LRS`")
						}
					}

					if createMode == "" || createMode == string(servers.CreateModeDefault) {
						if v := configMap["storage_iops"]; v.IsNull() {
							return errors.New("`storage_iops` is required when `storage_type` is `PremiumV2_LRS` and `create_mode` is `Default`")
						}

						if v := configMap["storage_throughput"]; v.IsNull() {
							return errors.New("`storage_throughput` is required when `storage_type` is `PremiumV2_LRS` and `create_mode` is `Default`")
						}
					}
				}

				if diff.Get("storage_type").(string) == string(servers.StorageTypePremiumLRS) {
					if v := configMap["storage_iops"]; !v.IsNull() {
						return errors.New("`storage_iops` is only supported when `storage_type` is `PremiumV2_LRS`")
					}

					if v := configMap["storage_throughput"]; !v.IsNull() {
						return errors.New("`storage_throughput` is only supported when `storage_type` is `PremiumV2_LRS`")
					}
				}

				return nil
			}, func(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
				// only pg 17+ support cluster, and cluster does not support major version upgrade
				if _, ok := diff.GetOk("cluster"); !ok {
					return nil
				}

				// can't downgrade cluster size
				if diff.HasChange("cluster") {
					oldCluster, newCluster := diff.GetChange("cluster")
					oldClusterObj := oldCluster.([]any)
					newClusterObj := newCluster.([]any)
					if len(oldClusterObj) > 0 && len(newClusterObj) > 0 {
						var oldClusterSize, newClusterSize int
						if tmpObj, ok := oldClusterObj[0].(map[string]any); ok {
							oldClusterSize, _ = tmpObj["size"].(int)
						}
						if tempObj, ok := newClusterObj[0].(map[string]any); ok {
							newClusterSize, _ = tempObj["size"].(int)
						}

						if newClusterSize < oldClusterSize {
							diff.ForceNew("cluster.0.size")
						}
					}
				}

				create_mode := servers.CreateMode(diff.Get("create_mode").(string))
				if string(create_mode) != "" && create_mode != servers.CreateModeDefault {
					return errors.New("`cluster` is only supported when `create_mode` is Default")
				}

				versionOld, versionNew := diff.GetChange("version")
				// version should be set since cluster only support pg17+
				if version, err := strconv.ParseInt(versionNew.(string), 10, 32); err != nil {
					return fmt.Errorf("parsing `version` %q: %v", versionNew.(string), err)
				} else if version < 17 {
					return errors.New("`cluster` is only supported for PostgreSQL major version 17 or above")
				}

				// major version upgrade is not allowed when cluster is enabled
				if versionOld.(string) != "" && versionOld.(string) != versionNew.(string) {
					return errors.New("major version upgrade not supported when `cluster` is set")
				}

				return nil
			}, func(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
				if diff.HasChange("sku_name") {
					skuOld, skuNew := diff.GetChange("sku_name")
					return validate.FlexibleServerSkuNameChange(skuOld.(string), skuNew.(string))
				}

				return nil
			}, func(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
				if diff.GetRawConfig().AsValueMap()["storage_mb"].IsNull() {
					return nil
				}

				storageMb := diff.Get("storage_mb").(int)

				if diff.Get("storage_type").(string) == string(servers.StorageTypePremiumVTwoLRS) {
					if storageMb%1024 != 0 {
						return fmt.Errorf("`storage_mb` must be a multiple of `1024` when `storage_type` is `PremiumV2_LRS`, got `%d`", storageMb)
					}

					return nil
				}

				validStorageMb := slices.Sorted(maps.Keys(validate.InitializeFlexibleServerStorageTierDefaults()))
				if !slices.Contains(validStorageMb, storageMb) {
					return fmt.Errorf("`storage_mb` must be one of %v when `storage_type` is `Premium_LRS`, got `%d`", validStorageMb, storageMb)
				}

				return nil
			},
		),
	}
}
