// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_vault

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataProtectionBackupVaultRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataProtection.BackupVaultClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := backupvaultresources.ParseBackupVaultID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.BackupVaultsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] DataProtection BackupVault %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving DataProtection BackupVault (%q): %+v", id, err)
	}
	d.Set("name", id.BackupVaultName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(pointer.To(model.Location)))
		props := model.Properties

		immutability := backupvaultresources.ImmutabilityStateDisabled
		if securitySetting := model.Properties.SecuritySettings; securitySetting != nil {
			if immutabilitySettings := securitySetting.ImmutabilitySettings; immutabilitySettings != nil {
				if immutabilitySettings.State != nil {
					immutability = *immutabilitySettings.State
				}
			}
			if softDelete := securitySetting.SoftDeleteSettings; softDelete != nil {
				d.Set("soft_delete", pointer.FromEnum(softDelete.State))
				d.Set("retention_duration_in_days", pointer.From(softDelete.RetentionDurationInDays))
			}
		}
		d.Set("immutability", string(immutability))

		crossRegionStoreEnabled := false
		if featureSetting := model.Properties.FeatureSettings; featureSetting != nil {
			if crossRegionRestore := featureSetting.CrossRegionRestoreSettings; crossRegionRestore != nil {
				if pointer.From(crossRegionRestore.State) == backupvaultresources.CrossRegionRestoreStateEnabled {
					crossRegionStoreEnabled = true
				}
			}
		}
		d.Set("cross_region_restore_enabled", crossRegionStoreEnabled)

		identity, err := flattenBackupVaultDppIdentityDetails(model.Identity)
		if err != nil {
			return err
		}
		d.Set("identity", identity)

		if err = tags.FlattenAndSet(d, flattenTags(model.Tags)); err != nil {
			return err
		}

		if !features.SixPointOh() {
			datastoreType, redundancy := flattenBackupVaultStorageSettingsLegacy(props.StorageSettings, d.Get("datastore_type").(string))
			d.Set("datastore_type", datastoreType)
			d.Set("redundancy", redundancy)
			return pluginsdk.SetResourceIdentityData(d, id)
		}

		if err := d.Set("storage_setting", flattenBackupVaultStorageSettings(props.StorageSettings)); err != nil {
			return fmt.Errorf("setting `storage_setting`: %+v", err)
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenBackupVaultStorageSettings(input []backupvaultresources.StorageSetting) []any {
	result := make([]any, 0, len(input))
	for _, setting := range input {
		result = append(result, map[string]any{
			"datastore_type": pointer.FromEnum(setting.DatastoreType),
			"redundancy":     pointer.FromEnum(setting.Type),
		})
	}
	return result
}

func flattenBackupVaultStorageSettingsLegacy(input []backupvaultresources.StorageSetting, currentDatastoreType string) (string, string) {
	if len(input) == 0 {
		return "", ""
	}

	setting := input[0]
	// Preserve an existing store selection when it is still present in Azure.
	// Without prior state, prefer ArchiveStore over its required VaultStore companion.
	for _, candidate := range input {
		if currentDatastoreType != "" && pointer.FromEnum(candidate.DatastoreType) == currentDatastoreType {
			return pointer.FromEnum(candidate.DatastoreType), pointer.FromEnum(candidate.Type)
		}
		if pointer.From(candidate.DatastoreType) == backupvaultresources.StorageSettingStoreTypesArchiveStore {
			setting = candidate
		}
	}

	return pointer.FromEnum(setting.DatastoreType), pointer.FromEnum(setting.Type)
}

func flattenBackupVaultDppIdentityDetails(input *backupvaultresources.DppIdentityDetails) (*[]any, error) {
	var config *identity.SystemAndUserAssignedMap
	if input != nil {
		config = &identity.SystemAndUserAssignedMap{
			Type: identity.Type(*input.Type),
		}

		config.PrincipalId = pointer.From(input.PrincipalId)
		config.TenantId = pointer.From(input.TenantId)

		if len(pointer.From(input.UserAssignedIdentities)) > 0 {
			config.IdentityIds = make(map[string]identity.UserAssignedIdentityDetails, len(pointer.From(input.UserAssignedIdentities)))
			for k, v := range *input.UserAssignedIdentities {
				config.IdentityIds[k] = identity.UserAssignedIdentityDetails{
					ClientId:    v.ClientId,
					PrincipalId: v.PrincipalId,
				}
			}
		}
	}

	return identity.FlattenSystemAndUserAssignedMap(config)
}
