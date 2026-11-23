// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_database

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/databases"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/resources/mssql_database/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name mssql_database -service-package-name mssql -properties "name" -compare-values "resource_group_name:server_id,server_name:server_id"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMsSqlDatabaseCreate,
		Read:   resourceMsSqlDatabaseRead,
		Update: resourceMsSqlDatabaseUpdate,
		Delete: resourceMsSqlDatabaseDelete,

		Importer: pluginsdk.ImporterValidatingIdentityThen(&commonids.SqlDatabaseId{}, resourceMsSqlDatabaseImporter),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&commonids.SqlDatabaseId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.DatabaseV0ToV1{},
		}),

		Schema: resourceMsSqlDatabaseSchema(),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIf("free_limit_enabled", func(ctx context.Context, d *schema.ResourceDiff, _ any) bool {
				// An existing database cannot be converted to the free tier.
				old, new := d.GetChange("free_limit_enabled")
				return d.Id() != "" && !old.(bool) && new.(bool)
			}),
			pluginsdk.ForceNewIf("free_limit_exhaustion_behavior", func(ctx context.Context, d *schema.ResourceDiff, _ any) bool {
				// Switching from BillOverUsage back to AutoPause is rejected by the API while the free offer remains enabled.
				old, new := d.GetChange("free_limit_exhaustion_behavior")
				return d.Id() != "" && d.Get("free_limit_enabled").(bool) &&
					old.(string) == string(databases.FreeLimitExhaustionBehaviorBillOverUsage) &&
					new.(string) == string(databases.FreeLimitExhaustionBehaviorAutoPause)
			}),
			pluginsdk.ForceNewIfChange("sku_name", func(ctx context.Context, old, new, _ any) bool {
				// hyperscale can not be changed to another sku
				return strings.HasPrefix(old.(string), "HS") && !strings.HasPrefix(new.(string), "HS")
			}),
			pluginsdk.ForceNewIfChange("enclave_type", func(ctx context.Context, old, new, _ any) bool {
				// enclave_type cannot be removed once it has been set
				// but can be changed between VBS and Default...
				// this Diff will not work until 4.0 when we remove
				// the computed property from the field schema.
				if old.(string) != "" && new.(string) == "" {
					return true
				}

				return false
			}),
			func(ctx context.Context, d *schema.ResourceDiff, meta any) error {
				transparentDataEncryption := d.Get("transparent_data_encryption_enabled").(bool)
				skuName := d.Get("sku_name").(string)

				if !strings.HasPrefix(strings.ToLower(skuName), "dw") && !transparentDataEncryption {
					return fmt.Errorf("transparent data encryption can only be disabled on Data Warehouse SKUs")
				}

				// NOTE: VBS enclaves are not supported by DW or DC skus...
				if d.Get("enclave_type").(string) == string(databases.AlwaysEncryptedEnclaveTypeVBS) {
					if strings.HasPrefix(strings.ToLower(skuName), "dw") || strings.Contains(strings.ToLower(skuName), "_dc_") {
						return fmt.Errorf("virtualization based security (VBS) enclaves are not supported for the %q sku", skuName)
					}
				}

				if strings.HasPrefix(strings.ToLower(skuName), "dw") {
					// NOTE: Got `PerDatabaseCMKDWNotSupported` error from API when `sku_name` is set to `DW100c` and `transparent_data_encryption_key_vault_key_id` is specified
					keyVaultKeyId := d.Get("transparent_data_encryption_key_vault_key_id").(string)
					if keyVaultKeyId != "" {
						return fmt.Errorf("database-level CMK is not supported for Data Warehouse SKUs")
					}
					// NOTE: Got `InternalServerError` error from API when `sku_name` is set to `DW100c` and `transparent_data_encryption_key_automatic_rotation_enabled` is specified
					if d.Get("transparent_data_encryption_key_automatic_rotation_enabled").(bool) {
						return fmt.Errorf("transparent_data_encryption_key_automatic_rotation_enabled should not be specified when using Data Warehouse SKUs")
					}
				}

				return nil
			},
			func(ctx context.Context, d *pluginsdk.ResourceDiff, i any) error {
				if !strings.HasPrefix(d.Get("sku_name").(string), "GP_S_") && !strings.HasPrefix(d.Get("sku_name").(string), "HS_S_") {
					rawConfig := d.GetRawConfig().AsValueMap()
					minCapacityVal := rawConfig["min_capacity"]
					if minCapacityVal.IsKnown() && !minCapacityVal.IsNull() {
						minCapacityValFloat, _ := minCapacityVal.AsBigFloat().Float64()
						if minCapacityValFloat > 0 {
							return fmt.Errorf("`min_capacity` should only be specified when using a serverless database")
						}
					}
					autoPauseDelayVal := rawConfig["auto_pause_delay_in_minutes"]
					if autoPauseDelayVal.IsKnown() && !autoPauseDelayVal.IsNull() {
						autoPauseDelayValInt, _ := autoPauseDelayVal.AsBigFloat().Int64()
						if autoPauseDelayValInt > 0 {
							return fmt.Errorf("`auto_pause_delay_in_minutes` should only be specified when using a serverless database")
						}
					}
				}
				return nil
			},
			func(ctx context.Context, d *pluginsdk.ResourceDiff, _ any) error {
				raw := d.GetRawConfig()
				if raw.IsNull() || !raw.IsKnown() {
					return nil
				}
				rawConfig := raw.AsValueMap()

				behaviorVal, hasBehavior := rawConfig["free_limit_exhaustion_behavior"]
				behaviorSet := hasBehavior && behaviorVal.IsKnown() && !behaviorVal.IsNull()

				enabledVal, hasEnabled := rawConfig["free_limit_enabled"]
				if hasEnabled && !enabledVal.IsKnown() {
					return nil
				}
				useFreeLimit := d.Get("free_limit_enabled").(bool)

				if behaviorSet && !useFreeLimit {
					return errors.New("`free_limit_exhaustion_behavior` can only be set when `free_limit_enabled` is `true`")
				}

				if useFreeLimit {
					if hasBehavior && !behaviorVal.IsKnown() {
						return nil
					}
					storageVal, hasStorage := rawConfig["storage_account_type"]
					if hasStorage && !storageVal.IsKnown() {
						return nil
					}

					behavior := d.Get("free_limit_exhaustion_behavior").(string)
					if behavior == "" {
						// Azure defaults to `AutoPause` when omitted.
						behavior = string(databases.FreeLimitExhaustionBehaviorAutoPause)
					}
					if behavior == string(databases.FreeLimitExhaustionBehaviorAutoPause) && d.Get("storage_account_type").(string) != string(databases.BackupStorageRedundancyLocal) {
						return errors.New("`storage_account_type` must be `Local` when `free_limit_enabled` is `true` and `free_limit_exhaustion_behavior` is `AutoPause`")
					}
				}

				return nil
			},
		),
	}
}
