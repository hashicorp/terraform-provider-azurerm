// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_vault

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

func RegisterResource() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceDataProtectionBackupVaultCreateUpdate,
		Read:   resourceDataProtectionBackupVaultRead,
		Update: resourceDataProtectionBackupVaultCreateUpdate,
		Delete: resourceDataProtectionBackupVaultDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.DataProtectionBackupVaultV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&backupvaultresources.BackupVaultId{}),
		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&backupvaultresources.BackupVaultId{}),
		},

		Schema: dataProtectionBackupVaultSchema(),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(

			// Once `cross_region_restore_enabled` is enabled it cannot be disabled.
			pluginsdk.ForceNewIfChange("cross_region_restore_enabled", func(ctx context.Context, old, new, meta any) bool {
				return old.(bool) && new.(bool) != old.(bool)
			}),

			// Once `immutability` is enabled it cannot be disabled.
			pluginsdk.ForceNewIfChange("immutability", func(ctx context.Context, old, new, meta any) bool {
				return old.(string) == string(backupvaultresources.ImmutabilityStateLocked) && new.(string) != string(backupvaultresources.ImmutabilityStateLocked)
			}),

			pluginsdk.ForceNewIfChange("soft_delete", func(ctx context.Context, old, new, meta any) bool {
				return old.(string) == string(backupvaultresources.SoftDeleteStateAlwaysOn) && new.(string) != string(backupvaultresources.SoftDeleteStateAlwaysOn)
			}),

			pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v any) error {
				crossRegionRestore := d.GetRawConfig().AsValueMap()["cross_region_restore_enabled"]
				if !features.SixPointOh() {
					if !crossRegionRestore.IsNull() && d.Get("redundancy").(string) != string(backupvaultresources.StorageSettingTypesGeoRedundant) {
						return fmt.Errorf("`cross_region_restore_enabled` can only be specified when `redundancy` is specified for `GeoRedundant`")
					}
					return nil
				}

				if !d.GetRawConfig().GetAttr("storage_setting").IsWhollyKnown() {
					return nil
				}
				var redundancy string
				hasArchiveStore, hasVaultStore := false, false
				for _, raw := range d.Get("storage_setting").([]any) {
					setting := raw.(map[string]any)
					datastoreType := setting["datastore_type"].(string)
					if strings.EqualFold(datastoreType, string(backupvaultresources.StorageSettingStoreTypesArchiveStore)) {
						hasArchiveStore = true
					}
					if strings.EqualFold(datastoreType, string(backupvaultresources.StorageSettingStoreTypesVaultStore)) {
						hasVaultStore = true
						redundancy = setting["redundancy"].(string)
					} else if redundancy == "" {
						redundancy = setting["redundancy"].(string)
					}
				}
				if hasArchiveStore && !hasVaultStore {
					return fmt.Errorf("`storage_setting` must include `VaultStore` when `ArchiveStore` is specified")
				}
				if !crossRegionRestore.IsNull() && redundancy != string(backupvaultresources.StorageSettingTypesGeoRedundant) {
					// Cross region restore is only allowed on `GeoRedundant` vault.
					return fmt.Errorf("`cross_region_restore_enabled` can only be specified when `redundancy` is specified for `GeoRedundant`")
				}
				return nil
			}),
		),
	}

	if !features.SixPointOh() {
		resource.SchemaVersion = 0
		resource.StateUpgraders = nil
	}

	return resource
}
