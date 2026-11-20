// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storageaccounts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_account/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

var (
	storageAccountResourceName  = "azurerm_storage_account"
	storageKindsSupportsSkuTier = map[storageaccounts.Kind]struct{}{
		storageaccounts.KindBlobStorage: {},
		storageaccounts.KindFileStorage: {},
		storageaccounts.KindStorageVTwo: {},
	}
	storageKindsSupportHns = map[storageaccounts.Kind]struct{}{
		storageaccounts.KindBlobStorage:      {},
		storageaccounts.KindBlockBlobStorage: {},
		storageaccounts.KindStorageVTwo:      {},
	}
	storageKindsSupportLargeFileShares = map[storageaccounts.Kind]struct{}{
		storageaccounts.KindFileStorage: {},
		storageaccounts.KindStorageVTwo: {},
	}
	storageReplicationZonalMigrationPairs = map[string]string{
		"LRS":    "ZRS",
		"GRS":    "GZRS",
		"RAGRS":  "RAGZRS",
		"ZRS":    "LRS",
		"GZRS":   "GRS",
		"RAGZRS": "RAGRS",
	}
)

func RegisterResource() *pluginsdk.Resource {
	r := &pluginsdk.Resource{
		Create: resourceStorageAccountCreate,
		Read:   resourceStorageAccountRead,
		Update: resourceStorageAccountUpdate,
		Delete: resourceStorageAccountDelete,

		SchemaVersion: 4,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.AccountV0ToV1{},
			1: migration.AccountV1ToV2{},
			2: migration.AccountV2ToV3{},
			3: migration.AccountV3ToV4{},
		}),

		Importer: pluginsdk.ImporterValidatingIdentity(&commonids.StorageAccountId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&commonids.StorageAccountId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: storageAccountSchema(),
		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v any) error {
				if d.HasChange("account_kind") {
					accountKind, changedKind := d.GetChange("account_kind")
					// Don't do the check for create case.
					if accountKind != "" {
						if !features.SixPointOh() {
							// Prevent ForceNew on `StorageV2` -> `Storage`, this would fail as GPv1 accounts can no longer be provisioned
							// Prevent ForceNew on `StorageV2` -> `BlobStorage`, this would fail as legacy Blob accounts can no longer be provisioned
							// Azure is automatically migrating accounts in October due to this retirement.
							if accountKind == string(storageaccounts.KindStorageVTwo) && (changedKind == string(storageaccounts.KindStorage) || changedKind == string(storageaccounts.KindBlobStorage)) {
								url := "https://learn.microsoft.com/azure/storage/common/general-purpose-version-1-account-migration-overview#retirement-timeline-and-key-milestones"
								if changedKind == string(storageaccounts.KindBlobStorage) {
									url = "https://learn.microsoft.com/azure/storage/common/legacy-blob-storage-account-migration-overview#retirement-timeline-and-key-milestones"
								}

								return fmt.Errorf("`account_kind` of type `%[1]s` has been retired by Azure, changing from `%[2]s` to `%[1]s` is no longer possible. For additional information, see %[3]s", changedKind, storageaccounts.KindStorageVTwo, url)
							}
						}

						if accountKind != string(storageaccounts.KindStorage) && changedKind != string(storageaccounts.KindStorageVTwo) {
							log.Printf("[DEBUG] recreate storage account, couldn't be migrated from %q to %q", accountKind, changedKind)
							d.ForceNew("account_kind")
						} else {
							log.Printf("[DEBUG] storage account can be upgraded from %q to %q", accountKind, changedKind)
						}
					}
				}

				if d.HasChange("large_file_share_enabled") {
					lfsEnabled, changedEnabled := d.GetChange("large_file_share_enabled")
					if lfsEnabled.(bool) && !changedEnabled.(bool) {
						d.ForceNew("large_file_share_enabled")
					}
				}

				if d.Get("access_tier") != "" {
					accountKind := storageaccounts.Kind(d.Get("account_kind").(string))
					if _, ok := storageKindsSupportsSkuTier[accountKind]; !ok {
						keys := sortedKeysFromSlice(storageKindsSupportsSkuTier)
						return fmt.Errorf("`access_tier` is only available for accounts where `kind` is set to one of: %+v", strings.Join(keys, " / "))
					}
				}

				if d.HasChange("immutability_policy.0.state") {
					old, new := d.GetChange("immutability_policy.0.state")

					// The initial value can be either "Disabled" or "Unlocked".
					// API error: InvalidStorageAccountImmutabilityPolicy: Storage account immutability policy state can be set to "Unlocked, Disabled".
					//      Note that the immutability state can either be Disabled or Unlocked during account creation
					if old == "" && (new.(string) != string(storageaccounts.AccountImmutabilityPolicyStateDisabled) && new.(string) != string(storageaccounts.AccountImmutabilityPolicyStateUnlocked)) {
						return fmt.Errorf("initial value of `immutability_policy.0.state` can be either Disabled or Unlocked, got=%s", new)
					}

					// Only "Unlocked" state can be updated to "Locked"
					// API error: InvalidStorageAccountImmutabilityPolicy: Storage account immutability policy state can be set to "Unlocked, Disabled".
					//      Note that the immutability state can either be Disabled or Unlocked during account creation
					if new.(string) == string(storageaccounts.AccountImmutabilityPolicyStateLocked) && old.(string) != string(storageaccounts.AccountImmutabilityPolicyStateUnlocked) {
						return fmt.Errorf("`immutability_policy.0.state` can only be set to Locked from Unlocked, got=%s", old)
					}

					// Once "Locked", can't be changed.
					// API error: InvalidStorageAccountImmutabilityPolicy: The immutability policy state cannot be changed once Locked.
					if old.(string) == string(storageaccounts.AccountImmutabilityPolicyStateLocked) {
						d.ForceNew("immutability_policy.0.state")
					}
				}

				isHnsEnabled := d.Get("is_hns_enabled").(bool)
				accountKind := storageaccounts.Kind(d.Get("account_kind").(string))
				if _, supportsHns := storageKindsSupportHns[accountKind]; !supportsHns && isHnsEnabled {
					keys := sortedKeysFromSlice(storageKindsSupportHns)
					return fmt.Errorf("`is_hns_enabled` can only be used for accounts with `account_kind` set to one of: %+v", strings.Join(keys, " / "))
				}
				return nil
			}),
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v any) error {
				if d.Get("account_replication_type_migration_in_progress").(bool) {
					o, n := d.GetChange("account_replication_type")
					migratingTo := d.Get("account_replication_type_migrating_to").(string)

					oldType, newType := o.(string), n.(string)
					if newType != migratingTo {
						return fmt.Errorf("a migration from `%s` to `%s` is in progress, please wait until this operation has completed before changing `account_replication_type` (configured type: `%s`)", oldType, migratingTo, newType)
					}
				}

				return nil
			}),
			pluginsdk.ForceNewIfChange("account_replication_type", func(ctx context.Context, old, new, meta any) bool {
				n := strings.ToUpper(new.(string))

				switch o := strings.ToUpper(old.(string)); o {
				case "LRS", "GRS", "RAGRS":
					if n == "GZRS" || n == "RAGZRS" || n == "ZRS" {
						return storageReplicationZonalMigrationPairs[o] != n
					}
				case "ZRS", "GZRS", "RAGZRS":
					if n == "LRS" || n == "GRS" || n == "RAGRS" {
						return storageReplicationZonalMigrationPairs[o] != n
					}
				}
				return false
			}),
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v any) error {
				if !features.SixPointOh() {
					// If both are `null`/unset, set diff to return default of `true` / `Enabled`
					// to ensure removal functions while O+C in 5.x
					if rawConfig := d.GetRawConfig().AsValueMap(); rawConfig["public_network_access"].IsNull() && rawConfig["public_network_access_enabled"].IsNull() {
						if err := d.SetNew("public_network_access", string(storageaccounts.PublicNetworkAccessEnabled)); err != nil {
							return err
						}
						if err := d.SetNew("public_network_access_enabled", true); err != nil {
							return err
						}
					}
				}

				return nil
			}),
		),
	}

	if !features.SixPointOh() {
		r.Schema["account_kind"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.All(
				validation.StringInSlice(storageaccounts.PossibleValuesForKind(), false),
				func(i any, k string) ([]string, []error) {
					v, ok := i.(string)
					if !ok {
						return nil, nil
					}

					if v == string(storageaccounts.KindStorage) || v == string(storageaccounts.KindBlobStorage) {
						return []string{fmt.Sprintf("type `%s` has been retired by Azure and will be removed in v6.0 of the AzureRM provider", v)}, nil
					}

					return nil, nil
				},
			),
			Default: string(storageaccounts.KindStorageVTwo),
		}

		r.Schema["public_network_access_enabled"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeBool,
			Optional: true,
			// Note: O+C because in 5.x this value may be affected by `public_network_access`
			Computed:      true,
			ConflictsWith: []string{"public_network_access"},
			Deprecated:    "`public_network_access_enabled` has been deprecated in favour of `public_network_access` and will be removed in v6.0 of the AzureRM provider",
		}

		r.Schema["public_network_access"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Optional: true,
			// Note: O+C because in 5.x this value may be affected by `public_network_access_enabled`
			Computed:      true,
			ConflictsWith: []string{"public_network_access_enabled"},
			ValidateFunc:  validation.StringInSlice(storageaccounts.PossibleValuesForPublicNetworkAccess(), false),
		}
	}

	return r
}
