// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_vault

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataProtectionBackupVaultCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DataProtection.BackupVaultClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := backupvaultresources.NewBackupVaultID(subscriptionId, resourceGroup, name)

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.BackupVaultsGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing DataProtection BackupVault (%q): %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_data_protection_backup_vault", id.ID())
			}
		}
	}

	expandedIdentity, err := expandBackupVaultDppIdentityDetails(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	parameters := backupvaultresources.BackupVaultResource{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: backupvaultresources.BackupVault{
			SecuritySettings: &backupvaultresources.SecuritySettings{
				SoftDeleteSettings: &backupvaultresources.SoftDeleteSettings{
					State: pointer.ToEnum[backupvaultresources.SoftDeleteState](d.Get("soft_delete").(string)),
				},
				ImmutabilitySettings: &backupvaultresources.ImmutabilitySettings{
					State: pointer.ToEnum[backupvaultresources.ImmutabilityState](d.Get("immutability").(string)),
				},
			},
		},
		Identity: expandedIdentity,
		Tags:     expandTags(d.Get("tags").(map[string]any)),
	}

	if v, ok := d.GetOk("storage_setting"); ok {
		parameters.Properties.StorageSettings = expandBackupVaultStorageSettings(v.([]any))
	}
	if !features.SixPointOh() {
		parameters.Properties.StorageSettings = []backupvaultresources.StorageSetting{
			{
				DatastoreType: pointer.ToEnum[backupvaultresources.StorageSettingStoreTypes](d.Get("datastore_type").(string)),
				Type:          pointer.ToEnum[backupvaultresources.StorageSettingTypes](d.Get("redundancy").(string)),
			},
		}
		// Azure requires a VaultStore alongside ArchiveStore; use the configured redundancy for both.
		if strings.EqualFold(d.Get("datastore_type").(string), string(backupvaultresources.StorageSettingStoreTypesArchiveStore)) {
			parameters.Properties.StorageSettings = append(parameters.Properties.StorageSettings, backupvaultresources.StorageSetting{
				DatastoreType: pointer.To(backupvaultresources.StorageSettingStoreTypesVaultStore),
				Type:          pointer.ToEnum[backupvaultresources.StorageSettingTypes](d.Get("redundancy").(string)),
			})
		}
	}

	if !pluginsdk.IsExplicitlyNullInConfig(d, "cross_region_restore_enabled") {
		parameters.Properties.FeatureSettings = &backupvaultresources.FeatureSettings{
			CrossRegionRestoreSettings: &backupvaultresources.CrossRegionRestoreSettings{},
		}
		if d.Get("cross_region_restore_enabled").(bool) {
			parameters.Properties.FeatureSettings.CrossRegionRestoreSettings.State = pointer.To(backupvaultresources.CrossRegionRestoreStateEnabled)
		} else {
			parameters.Properties.FeatureSettings.CrossRegionRestoreSettings.State = pointer.To(backupvaultresources.CrossRegionRestoreStateDisabled)
		}
	}

	if v, ok := d.GetOk("retention_duration_in_days"); ok {
		parameters.Properties.SecuritySettings.SoftDeleteSettings.RetentionDurationInDays = pointer.To(v.(float64))
	}

	if d.IsNewResource() {
		if err := client.BackupVaultsCreateOrUpdateCallbackThenPoll(ctx, id, parameters, backupvaultresources.DefaultBackupVaultsCreateOrUpdateOperationOptions(), sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}

		d.SetId(id.ID())
		if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
			return err
		}
	} else {
		if err := client.BackupVaultsCreateOrUpdateThenPoll(ctx, id, parameters, backupvaultresources.DefaultBackupVaultsCreateOrUpdateOperationOptions()); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceDataProtectionBackupVaultRead(d, meta)
}

func expandBackupVaultStorageSettings(input []any) []backupvaultresources.StorageSetting {
	result := make([]backupvaultresources.StorageSetting, 0, len(input))
	for _, raw := range input {
		setting := raw.(map[string]any)
		result = append(result, backupvaultresources.StorageSetting{
			DatastoreType: pointer.ToEnum[backupvaultresources.StorageSettingStoreTypes](setting["datastore_type"].(string)),
			Type:          pointer.ToEnum[backupvaultresources.StorageSettingTypes](setting["redundancy"].(string)),
		})
	}
	return result
}

func expandBackupVaultDppIdentityDetails(input []any) (*backupvaultresources.DppIdentityDetails, error) {
	config, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	identity := backupvaultresources.DppIdentityDetails{
		Type: pointer.To(string(config.Type)),
	}

	if len(config.IdentityIds) > 0 {
		identityIds := make(map[string]backupvaultresources.UserAssignedIdentity, len(config.IdentityIds))
		for id := range config.IdentityIds {
			identityIds[id] = backupvaultresources.UserAssignedIdentity{}
		}
		identity.UserAssignedIdentities = pointer.To(identityIds)
	}

	return &identity, nil
}
